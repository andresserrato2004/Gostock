package handlers_test

import (
	"BackEnd/internal/handlers"
	"BackEnd/internal/services"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGetStocks_BasicStructure(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := handlers.NewStockHandler(nil, nil)
	if handler == nil {
		t.Fatal("NewStockHandler devolvió nil")
	}
}

func TestRecommendStocks_IntegrationAttempt(t *testing.T) {

	svc := services.NewRecommendationService(nil)
	handler := handlers.NewStockHandler(nil, svc)

	if handler == nil {
		t.Error("Handler initialization failed")
	}
}
