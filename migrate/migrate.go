package main

import (
	"os"
	"path/filepath"

	"github.com/adamhaiqal/go-auth/config"
	"github.com/adamhaiqal/go-auth/models"
	"github.com/joho/godotenv"
)

func init() {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	godotenv.Load(filepath.Join(pwd, "../.env"))
	config.LoadEnvVariables()
	config.ConnectToDatabase()
}

func main() {

	config.DB.AutoMigrate(&models.Account{})

	println("Database migration completed successfully.")
}
