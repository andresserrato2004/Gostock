package services

import (
	"BackEnd/internal/models"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"gorm.io/gorm"
)

type RecommendationService struct {
	db     *gorm.DB
	apiKey string
	apiUrl string
}

func NewRecommendationService(db *gorm.DB) *RecommendationService {
	return &RecommendationService{
		db:     db,
		apiKey: os.Getenv("FINNHUB_API_KEY"),
		apiUrl: os.Getenv("FINNHUB_API_URL"),
	}
}

func (s *RecommendationService) GetTopRecommendations(limit int) ([]models.ScoredStock, error) {

	// 1. Obtener candidatos de la BD (filtrando por compras y actualizaciones positivas)
	var candidates []models.Stock
	// RESTRICTION: Limitamos a 30 candidatos debido al LÍMITE ESTRICTO de la API externa.
	// Si intentamos procesar más, recibiremos errores 429 (Too Many Requests).
	hardLimit := 30

	err := s.db.Where("target_to_num > ?", 0).
		Where(s.db.Where("rating_to ILIKE ?", "%Buy%").Or("action ILIKE ?", "%raised%")).
		Order("updated_at DESC").
		Limit(hardLimit).
		Find(&candidates).Error

	if err != nil {
		return nil, err
	}

	var results []models.ScoredStock
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Usamos un ticker conservador para espaciar las llamadas y evitar bursts
	rateLimiter := time.Tick(100 * time.Millisecond)

	for _, stock := range candidates {
		// Esperar al ticker para asegurar rate limit
		<-rateLimiter

		wg.Add(1)
		go func(st models.Stock) {
			defer wg.Done()

			// 2. Obtener datos en tiempo real (Precio y Momentum)
			quote, err := s.getRealTimeQuote(st.Ticker)
			if err != nil {
				log.Println("Error fetching data for", st.Ticker, ":", err)
				return // Skip this stock
			}

			currentPrice := quote.C
			momentumPct := quote.Dp // e.g. 2.5 for 2.5%

			if currentPrice <= 0 {
				return
			}

			// --- CÁLCULO DE VARIABLES DE LA FÓRMULA ---

			// P: Upside Potential (Potencial de Subida)
			// (Target - Current) / Current
			upside := (st.TargetToNum - currentPrice) / currentPrice

			// C: Consensus Rating Score (Consenso)
			// Normalizado 0.0 - 1.0
			ratingScore := s.normalizeRating(st.RatingTo)

			// E: Analyst Conviction (Convicción/Esperanza)
			// Cuánto subieron el target: (TargetTo - TargetFrom) / TargetFrom
			conviction := 0.0
			if st.TargetFromNum > 0 {
				conviction = (st.TargetToNum - st.TargetFromNum) / st.TargetFromNum
			}

			// M: Market Momentum (Momentum)
			// Convertir el porcentaje de Finnhub (enteros) a decimal para equiparar escala.
			// e.g. 5% -> 0.05
			momentumDecimal := momentumPct / 100.0

			// --- FÓRMULA 2.0 ---
			// Score = (0.4 * P) + (0.2 * C) + (0.2 * E) + (0.2 * M)
			wP, wC, wE, wM := 0.4, 0.2, 0.2, 0.2

			// Normalizamos P y E para que no se disparen infinitamente si hay un error de datos
			// Por ahora los dejamos crudos ("raw") pero asumimos que outliers serán raros en acciones serias.
			// Una corrección simple: Si P > 1 (100%), lo capeamos a 1 para el score, pero guardamos el real.
			pClamped := upside
			if pClamped > 1.0 {
				pClamped = 1.0
			} else if pClamped < -1.0 {
				pClamped = -1.0
			}

			// Lo mismo para Conviction
			eClamped := conviction
			if eClamped > 1.0 {
				eClamped = 1.0
			} else if eClamped < -1.0 {
				eClamped = -1.0
			}

			finalScore := (wP * pClamped) + (wC * ratingScore) + (wE * eClamped) + (wM * momentumDecimal)

			mu.Lock()
			results = append(results, models.ScoredStock{
				Ticker:       st.Ticker,
				Company:      st.Company,
				CurrentPrice: currentPrice,
				TargetPrice:  st.TargetToNum,
				Upside:       math.Round(upside*10000) / 100,     // Display as %
				RatingScore:  ratingScore,                        // 0-1
				Conviction:   math.Round(conviction*10000) / 100, // Display as %
				Momentum:     math.Round(momentumPct*100) / 100,  // Display as % (Finnhub format)
				FinalScore:   math.Round(finalScore*10000) / 100, // Scale to 0-100 roughly
				Reason:       fmt.Sprintf("Rating: %s, Action: %s", st.RatingTo, st.Action),
			})
			mu.Unlock()

		}(stock)
	}

	wg.Wait()

	// 3. Ordenar por FinalScore DESC
	sort.Slice(results, func(i, j int) bool {
		return results[i].FinalScore > results[j].FinalScore
	})

	if len(results) > limit {
		results = results[:limit]
	}

	return results, nil
}

func (s *RecommendationService) getRealTimeQuote(ticker string) (*models.FinnhubQuote, error) {
	if s.apiKey == "" {
		return nil, fmt.Errorf("no API key")
	}

	url := fmt.Sprintf("%s?symbol=%s&token=%s", s.apiUrl, ticker, s.apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("api error: %s", resp.Status)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var q models.FinnhubQuote
	if err := json.Unmarshal(bodyBytes, &q); err != nil {
		return nil, err
	}

	// Si Finnhub devuelve 0 (ticker no encontrado o error), manejamos
	// nota: precio 0 es sospechoso.
	if q.C == 0 {
		return nil, fmt.Errorf("price is 0 (ticker invalid?)")
	}

	return &q, nil
}

// Helper: Normalizar Rating
func (s *RecommendationService) normalizeRating(rating string) float64 {
	r := strings.ToLower(rating)
	if strings.Contains(r, "strong buy") {
		return 1.0
	}
	if strings.Contains(r, "buy") || strings.Contains(r, "outperform") || strings.Contains(r, "overweight") || strings.Contains(r, "top pick") {
		return 0.8
	}
	if strings.Contains(r, "hold") || strings.Contains(r, "neutral") || strings.Contains(r, "perform") || strings.Contains(r, "equal-weight") {
		return 0.5
	}
	return 0.2 // Sell, Underperform o desconocido
}
