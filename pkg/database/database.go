package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() (*gorm.DB, error) {
	postgresURL := os.Getenv("DATABASE_URL")
	if postgresURL == "" {
		return nil, fmt.Errorf("DATABASE_URL not set")
	}

	log.Printf("Attempting to connect to database...")

	DB, err := gorm.Open(postgres.Open(postgresURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	log.Println("Database connected")
	return DB, nil
}
