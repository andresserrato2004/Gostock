package handlers

import (
	"net/http"
	"strconv"

	"BackEnd/internal/models"
	"BackEnd/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type StockHandler struct {
	db *gorm.DB
}

func NewStockHandler(db *gorm.DB) *StockHandler {
	return &StockHandler{db: db}
}

func (h *StockHandler) GetStocks(c *gin.Context) {
	var stocks []models.Stock
	var total int64

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	search := c.Query("search")
	sortBy := c.DefaultQuery("sort_by", "target_to_num")
	order := c.DefaultQuery("order", "desc")

	query := h.db.Model(&models.Stock{})

	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("ticker ILIKE ? OR company ILIKE ?", searchPattern, searchPattern)
	}

	query.Count(&total)

	// Validate sort column to prevent SQL injection
	validSorts := map[string]bool{
		"ticker": true, "company": true, "target_to_num": true, "rating_to": true, "action": true,
	}
	if !validSorts[sortBy] {
		sortBy = "target_to_num"
	}
	if order != "asc" && order != "desc" {
		order = "desc"
	}

	result := query.Order(sortBy + " " + order).Offset(offset).Limit(limit).Find(&stocks)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": stocks,
		"meta": gin.H{
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}

func (h *StockHandler) RecommendStocks(c *gin.Context) {

	recServices := services.NewRecommendationService(h.db)

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	scoredStocks, err := recServices.GetTopRecommendations(limit)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return

	}

	if scoredStocks == nil {
		scoredStocks = []models.ScoredStock{}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": scoredStocks,
		"meta": gin.H{
			"count": len(scoredStocks),
			"info":  "algoritmo: upside (50%), Rating (30%), Sentiment (20%)",
		},
	})

}
