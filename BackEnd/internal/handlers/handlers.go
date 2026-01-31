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

// GetStocks obtiene una lista paginada de acciones con opciones de búsqueda y ordenamiento.
//
// Recibe query params: `page`, `limit`, `search`, `sort_by`, `order`.
//
// Devuelve un objeto JSON con la data de las acciones y metadatos de paginación.
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

// RecommendStocks maneja la solicitud de recomendaciones de inversión avanzadas.
//
// Recibe un int `limit` (via query param) que especifica el número máximo de recomendaciones a devolver (default 10).
//
// Devuelve una lista de `ScoredStock` calculada en tiempo real combinando datos de analistas y mercado.
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
	})

}
