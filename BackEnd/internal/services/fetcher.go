package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"BackEnd/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type StockService struct {
	db *gorm.DB
}

func NewStockService(db *gorm.DB) *StockService {
	return &StockService{db: db}
}

func (s *StockService) IngestStocks() error {
	// Verificar si ya existen datos para evitar recargas innecesarias al reiniciar
	var count int64
	if err := s.db.Model(&models.Stock{}).Count(&count).Error; err != nil {
		log.Println("Advertencia: No se pudo verificar el conteo de registros:", err)
	}
	if count > 0 {
		log.Printf("La base de datos ya contiene %d registros. Saltando ingestión inicial.\n", count)
		return nil
	}

	baseURL := os.Getenv("API_URL")
	if baseURL == "" {
		return fmt.Errorf("API_URL no esta en el archivo .env")
	}

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		return fmt.Errorf("API_KEY no esta en el archivo .env")
	}

	nextPage := ""

	for {

		FullUrl := baseURL
		if nextPage != "" {
			FullUrl = fmt.Sprintf("%s?next_page=%s", baseURL, nextPage)
		}

		log.Println("Fetching URL:", FullUrl)

		var apiResp models.ApiResponse

		err := s.fetchWithRetries(FullUrl, apiKey, &apiResp, 3)
		if err != nil {
			return err
		}

		if len(apiResp.Data) > 0 {

			s.cleanData(apiResp.Data)

			err := s.db.Clauses(clause.OnConflict{
				Columns: []clause.Column{{Name: "ticker"}},
				DoUpdates: clause.AssignmentColumns([]string{
					"company", "brokerage", "action", "rating_from", "rating_to", "target_from_num", "target_to_num",
				}),
			}).Create(&apiResp.Data).Error

			if err != nil {
				log.Println("Error inserting stocks:", err)
			}
		}

		nextPage = apiResp.NextPage
		if nextPage == "" {
			break
		}

		time.Sleep(500 * time.Millisecond)
	}

	log.Println("datos melos")

	return nil
}

func (s *StockService) cleanData(stocks []models.Stock) {
	for i := range stocks {
		stocks[i].TargetFromNum = parseCurrency(stocks[i].TargetFrom)
		stocks[i].TargetToNum = parseCurrency(stocks[i].TargetTo)
	}
}

func parseCurrency(value string) float64 {
	cleaned := strings.TrimSpace(strings.ReplaceAll(value, "$", ""))
	val, err := strconv.ParseFloat(cleaned, 64)
	if err != nil {
		return 0.0
	}
	return val
}

func (s *StockService) fetchWithRetries(url, token string, target interface{}, maxRetries int) error {
	client := &http.Client{Timeout: 10 * time.Second}

	for i := 0; i < maxRetries; i++ {
		req, _ := http.NewRequest("GET", url, nil)

		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := client.Do(req)
		if err == nil && resp.StatusCode == http.StatusOK {
			defer resp.Body.Close()
			return json.NewDecoder(resp.Body).Decode(target)
		}

		if err != nil {
			log.Printf("Error en la solicitud: %v", err)
		} else {
			log.Printf("Respuesta no exitosa: %s", resp.Status)
			resp.Body.Close()
		}

		time.Sleep(time.Duration(i+1) * time.Second)
	}
	return fmt.Errorf("no se pudo obtener datos tras %d intentos", maxRetries)
}
