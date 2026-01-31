package config

import (
	"fmt"
	"log"
	"os"

	"BackEnd/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL no está configurada")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error al conectar a la base de datos:", err)
	}

	fmt.Println("Conexión exitosa a CockroachDB")

	err = db.AutoMigrate(&models.Stock{})
	if err != nil {
		log.Println("Error en la migración", err)
	}

	DB = db
}
