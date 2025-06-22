package initializers

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
		log.Fatal("DATABASE_URL environment variable not set")
	}
	DB, err = gorm.Open(postgres.Open(postgresURL), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database")
	}
	return DB
}
