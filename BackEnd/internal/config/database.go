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
	// El DSN usualmente tiene este formato:
	// "postgresql://<user>:<pass>@<host>:<port>/<db>?sslmode=verify-full"
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL no está configurada")
	}
	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error al conectar a la base de datos:", err)
	}

	fmt.Println("Conexión exitosa a CockroachDB")

	// AutoMigrate crea las tablas automáticamente basadas en el struct
	// Forzar drop table si el schema está corrupto (SOLO PARA DEBUG)
	// db.Migrator().DropTable(&models.Stock{})

	err = db.AutoMigrate(&models.Stock{})
	if err != nil {
		log.Println("Error en la migración (ignorando para continuar):", err)
		// No detener ejecución con log.Fatal
	}

	DB = db
}
