package main

import (
	"os"
	"path/filepath"

	"github.com/adamhaiqal/go-auth/pkg/database"
	"github.com/adamhaiqal/go-auth/pkg/models"
	"github.com/joho/godotenv"
)

func init() {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	godotenv.Load(filepath.Join(pwd, "../.env"))
	database.Connect()
}

func main() {
	database.DB.AutoMigrate(&models.Account{})
	database.DB.AutoMigrate(&models.AccountMessage{})

	println("Database migration completed successfully.")
}
