package routes

import (
	"BackEnd/internal/handlers"
	"BackEnd/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {

	recService := services.NewRecommendationService(db)
	stockHandler := handlers.NewStockHandler(db, recService)

	apiGroup := r.Group("/api")
	{
		apiGroup.GET("/stocks", stockHandler.GetStocks)
		apiGroup.GET("/recommend", stockHandler.RecommendStocks)
	}
}
