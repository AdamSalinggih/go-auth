package initializers

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDatabase() *gorm.DB {
	var err error

	postgresURL := os.Getenv("DATABASE_URL")
	fmt.Print("Connecting to database at: ", postgresURL)

	DB, err = gorm.Open(postgres.Open(postgresURL), &gorm.Config{})

	if err != nil {
		log.Fatal("failed to connect to database")
	}
	return DB
}
