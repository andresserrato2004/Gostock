package services

import (
	"BackEnd/internal/models"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"gorm.io/gorm"
)

var ErrRateLimited = errors.New("rate limited by finnhub")

type RecommendationService struct {
	db         *gorm.DB
	apiKey     string
	apiUrl     string
	httpClient *http.Client
}

func NewRecommendationService(db *gorm.DB) *RecommendationService {
	return &RecommendationService{
		db:         db,
		apiKey:     os.Getenv("FINNHUB_API_KEY"),
		apiUrl:     os.Getenv("FINNHUB_API_URL"),
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// GetTopRecommendations obtiene las mejores recomendaciones de acciones para realizar inversiones.
//
// Recibe un int `limit` que especifica el número máximo de recomendaciones a devolver.
// esto por que Finnhub tiene un limite de peticiones por minuto en su plan gratuito que es 30.
//
// Devuelve una lista de `ScoredStock` ordenada por el puntaje calculado, junto con un error si ocurre alguno.
func (s *RecommendationService) GetTopRecommendations(limit int) ([]models.ScoredStock, error) {

	var candidates []models.Stock
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
	var rateLimited atomic.Bool

	rateLimiter := time.Tick(100 * time.Millisecond)

	for _, stock := range candidates {
		<-rateLimiter

		wg.Add(1)
		go func(st models.Stock) {
			defer wg.Done()

			quote, err := s.getRealTimeQuote(st.Ticker)
			if err != nil {
				log.Println("Error fetching data for", st.Ticker, ":", err)
				if errors.Is(err, ErrRateLimited) {
					rateLimited.Store(true)
				}
				return
			}

			// Extract logic to method for better testing
			scoredStock := s.CalculateScore(st, *quote)

			mu.Lock()
			results = append(results, scoredStock)
			mu.Unlock()

		}(stock)
	}

	wg.Wait()

	if rateLimited.Load() {
		return nil, ErrRateLimited
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].FinalScore > results[j].FinalScore
	})

	if len(results) > limit {
		results = results[:limit]
	}

	return results, nil
}

// CalculateScore realiza el cálculo puro de la puntuación de una acción.
// Separado para facilitar pruebas unitarias.
func (s *RecommendationService) CalculateScore(st models.Stock, quote models.FinnhubQuote) models.ScoredStock {
	currentPrice := quote.C
	momentumPct := quote.Dp

	if currentPrice <= 0 {
		return models.ScoredStock{}
	}

	// P: Upside Potential (Potencial de Subida)
	upside := (st.TargetToNum - currentPrice) / currentPrice

	// C: Consensus Rating Score (Consenso)
	ratingScore := s.normalizeRating(st.RatingTo)

	// E: Analyst Conviction (Convicción/Esperanza)
	conviction := 0.0
	if st.TargetFromNum > 0 {
		conviction = (st.TargetToNum - st.TargetFromNum) / st.TargetFromNum
	}

	// M: Market Momentum (Momentum)
	momentumDecimal := momentumPct / 100.0

	// Score = (0.4 * P) + (0.2 * C) + (0.2 * E) + (0.2 * M)
	wP, wC, wE, wM := 0.4, 0.2, 0.2, 0.2

	pClamped := upside
	if pClamped > 1.0 {
		pClamped = 1.0
	} else if pClamped < -1.0 {
		pClamped = -1.0
	}

	eClamped := conviction
	if eClamped > 1.0 {
		eClamped = 1.0
	} else if eClamped < -1.0 {
		eClamped = -1.0
	}

	finalScore := (wP * pClamped) + (wC * ratingScore) + (wE * eClamped) + (wM * momentumDecimal)

	return models.ScoredStock{
		Ticker:       st.Ticker,
		Company:      st.Company,
		CurrentPrice: currentPrice,
		TargetPrice:  st.TargetToNum,
		Upside:       math.Round(upside*10000) / 100,
		RatingScore:  ratingScore,
		Conviction:   math.Round(conviction*10000) / 100,
		Momentum:     math.Round(momentumPct*100) / 100,
		FinalScore:   math.Round(finalScore*10000) / 100,
		Reason:       fmt.Sprintf("Rating: %s, Action: %s", st.RatingTo, st.Action),
	}
}

// getRealTimeQuote tiene como parametro el ticker de una acción/compañia
// obtiene la cotización en tiempo real de una acción utilizando la API de Finnhub.
//
// retorna un puntero a FinnhubQuote y un error si ocurre alguno.
func (s *RecommendationService) getRealTimeQuote(ticker string) (*models.FinnhubQuote, error) {
	if s.apiKey == "" {
		return nil, fmt.Errorf("no API key")
	}

	url := fmt.Sprintf("%s?symbol=%s&token=%s", s.apiUrl, ticker, s.apiKey)

	resp, err := s.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		if resp.StatusCode == http.StatusTooManyRequests {
			return nil, fmt.Errorf("%w: %s", ErrRateLimited, resp.Status)
		}
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

	if q.C == 0 {
		return nil, fmt.Errorf("price is 0 (ticker invalid?)")
	}

	return &q, nil
}

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
	return 0.2
}
