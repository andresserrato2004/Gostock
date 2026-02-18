package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Stock representa la información de una acción de una empresa almacenada en la base de datos.
// Se utiliza GORM para el mapeo objeto-relacional.

type Stock struct {
	ID            uint    `gorm:"primaryKey" json:"id"`
	Ticker        string  `gorm:"uniqueIndex;not null" json:"ticker"` // Símbolo único de la empresa (e.g., AAPL)
	Company       string  `gorm:"not null" json:"company"`            // Nombre completo de la compañía
	Brokerage     string  `json:"brokerage"`                          // Firma de corretaje que emite la calificación
	Action        string  `json:"action"`                             // Acción recomendada
	RatingFrom    string  `json:"rating_from"`                        // Calificación anterior
	RatingTo      string  `json:"rating_to"`                          // Nueva calificación
	TargetFrom    string  `gorm:"-" json:"target_from"`               // Precio objetivo anterior
	TargetTo      string  `gorm:"-" json:"target_to"`                 // Precio objetivo nuevo
	TargetFromNum float64 `gorm:"type:decimal(10,2)" json:"-"`        // Precio objetivo anterior
	TargetToNum   float64 `gorm:"type:decimal(10,2)" json:"-"`        // Precio objetivo nuevo

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Empresas struct {
	Ticker      string `json:"ticker"`
	Target_from string `json:"target_from"`
	Target_to   string `json:"target_to"`
	Company     string `json:"company"`
	Action      string `json:"action"`
	Brokerage   string `json:"brokerage"`
	Rating_from string `json:"rating_from"`
	Rating_to   string `json:"rating_to"`
	Time        string `json:"time"`
}

// AfterFind es un hook de GORM que se ejecuta después de consultar una Stock.
// Se utiliza para formatear los valores numéricos a string para la respuesta JSON.

func (s *Stock) AfterFind(tx *gorm.DB) (err error) {
	s.TargetFrom = fmt.Sprintf("%.2f", s.TargetFromNum)
	s.TargetTo = fmt.Sprintf("%.2f", s.TargetToNum)
	return
}

// ApiResponse estructura la respuesta de la API externa (KarenAI) que se ingesta.
type ApiResponse struct {
	Data     []Stock `json:"items"`     // Lista de stocks devueltos en la página
	NextPage string  `json:"next_page"` // Token o URL para la siguiente página de resultados
}

type Creacion struct {
	Items    []Empresas `json:"items"`     // Lista de stocks devueltos en la página
	NextPage string     `json:"next_page"` // Token o URL para la siguiente página de resultados
}

// ScoredStock representa una acción evaluada y puntuada por el algoritmo de recomendación.
type ScoredStock struct {
	Ticker       string  `json:"ticker"`
	Company      string  `json:"company"`
	CurrentPrice float64 `json:"current_price"` // Precio actual de mercado
	TargetPrice  float64 `json:"target_price"`  // Precio objetivo promedio de analistas
	Upside       float64 `json:"upside"`        // P: Potencial de subida en porcentaje (Target - Current) / Current
	RatingScore  float64 `json:"rating_score"`  // C: Puntuación normalizada del consenso de analistas (0.0 - 1.0)
	Conviction   float64 `json:"conviction"`    // E: Convicción (Cambio porcentual en el precio objetivo)
	Momentum     float64 `json:"momentum"`      // M: Momentum de mercado (Cambio porcentual diario)
	FinalScore   float64 `json:"final_score"`   // Puntuación final calculada (0 - 100)
	Reason       string  `json:"reason"`        // Razón descriptiva de la recomendación
}

// FinnhubQuote mapea la respuesta de cotización en tiempo real de la API de Finnhub.
type FinnhubQuote struct {
	C  float64 `json:"c"`  // Current price (Precio actual)
	Dp float64 `json:"dp"` // Percent change (Porcentaje de cambio diario)
}
