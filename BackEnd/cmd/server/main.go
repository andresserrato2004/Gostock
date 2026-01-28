package main

import (
	"log"

	"BackEnd/internal/config"

	"BackEnd/internal/middleware"
	"BackEnd/internal/routes"
	"BackEnd/internal/services"

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

	r.Use(middleware.CORSMiddleware())

	routes.SetupRoutes(r, config.DB)

	log.Println("Server running on port 8080")
	r.Run(":8080")
}
