package config

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDatabase() *gorm.DB {
	var err error

	postgresURL := os.Getenv("DATABASE_URL")
	if postgresURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	log.Printf("Attempting to connect to database...")

	DB, err = gorm.Open(postgres.Open(postgresURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connection established successfully")
	return DB
}
