package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Stock struct {
	ID            uint    `gorm:"primaryKey" json:"id"`
	Ticker        string  `gorm:"uniqueIndex;not null" json:"ticker"`
	Company       string  `gorm:"not null" json:"company"`
	Brokerage     string  `json:"brokerage"`
	Action        string  `json:"action"`
	RatingFrom    string  `json:"rating_from"`
	RatingTo      string  `json:"rating_to"`
	TargetFrom    string  `gorm:"-" json:"target_from"`
	TargetTo      string  `gorm:"-" json:"target_to"`
	TargetFromNum float64 `gorm:"type:decimal(10,2)" json:"-"`
	TargetToNum   float64 `gorm:"type:decimal(10,2)" json:"-"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (s *Stock) AfterFind(tx *gorm.DB) (err error) {
	s.TargetFrom = fmt.Sprintf("%.2f", s.TargetFromNum)
	s.TargetTo = fmt.Sprintf("%.2f", s.TargetToNum)
	return
}

type ApiResponse struct {
	Data     []Stock `json:"items"`
	NextPage string  `json:"next_page"`
}

type ScoredStock struct {
	Ticker       string  `json:"ticker"`
	Company      string  `json:"company"`
	CurrentPrice float64 `json:"current_price"`
	TargetPrice  float64 `json:"target_price"`
	Upside       float64 `json:"upside"`       // P: Potencial de subida (decimal)
	RatingScore  float64 `json:"rating_score"` // C: Consenso Analistas (0-1)
	Conviction   float64 `json:"conviction"`   // E: Esperanza/Convicción (cambio en target)
	Momentum     float64 `json:"momentum"`     // M: Momentum de mercado (diario)
	FinalScore   float64 `json:"final_score"`
	Reason       string  `json:"reason"`
}

type FinnhubQuote struct {
	C  float64 `json:"c"`  // Precio actual
	Dp float64 `json:"dp"` // Porcentaje de cambio diario
}
