package main

import (
	"log"

	"BackEnd/internal/config" // Asegúrate de que coincida con go.mod
	// Usamos la carpeta
	"BackEnd/internal/routes"
	"BackEnd/internal/services"

	// Usamos la carpeta 'services'

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No se encontró archivo .env, usando variables de entorno del sistema")
	}
	config.ConnectDB()

	go func() {
		log.Println("Starting data ingestion...")
		svc := services.NewStockService(config.DB)
		if err := svc.IngestStocks(); err != nil {
			log.Println("Error during data ingestion:", err)
		}
	}()

	r := gin.Default()

	// Configuración CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Rutas lo tengo que pasar a la carpeta de routes
	routes.SetupRoutes(r, config.DB)

	log.Println("Server running on port 8080")
	r.Run(":8080")
}
