package routes

import (
	"BackEnd/internal/handlers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {

	stockHandler := handlers.NewStockHandler(db)
	apiGroup := r.Group("/api")
	{
		apiGroup.GET("/stocks", stockHandler.GetStocks)
		apiGroup.GET("/recommend", stockHandler.RecommendStocks)
	}
}
