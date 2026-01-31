package services_test

import (
	"BackEnd/internal/models"
	"BackEnd/internal/services"
	"testing"
)

func TestCalculateScore(t *testing.T) {
	
	service := services.NewRecommendationService(nil) 

	tests := []struct {
		name     string
		stock    models.Stock
		quote    models.FinnhubQuote
		expected float64 
	}{
		{
			name: "Perfect Scenario: High Upside, Strong Buy, High Conviction, Good Momentum",
			stock: models.Stock{
				Ticker:        "WIN",
				TargetToNum:   200.0,
				TargetFromNum: 150.0, 
				RatingTo:      "Strong Buy",
			},
			quote: models.FinnhubQuote{
				C:  100.0, 
				Dp: 5.0,   
			},
			expected: 67.6,
		},
		{
			name: "Bad Scenario: Overvalued, Sell, Target Lowered, Bad Momentum",
			stock: models.Stock{
				Ticker:        "LOSE",
				TargetToNum:   50.0,
				TargetFromNum: 80.0, 
				RatingTo:      "Sell",
			},
			quote: models.FinnhubQuote{
				C:  100.0, 
				Dp: -10.0, 
			},
			expected: -25.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.CalculateScore(tt.stock, tt.quote)

			epsilon := 0.1
			if result.FinalScore < (tt.expected-epsilon) || result.FinalScore > (tt.expected+epsilon) {
				t.Errorf("CalculateScore() = %v, want %v", result.FinalScore, tt.expected)
			}
		})
	}
}
