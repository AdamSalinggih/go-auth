package main

import (
	"os"
	"path/filepath"

	"github.com/adamhaiqal/go-auth/initializers"
	"github.com/adamhaiqal/go-auth/models"
	"github.com/joho/godotenv"
)

func init() {
	// Load environment variables and connect to the database
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	//use ../.env because main.go inside /cmd
	godotenv.Load(filepath.Join(pwd, "../.env"))
	initializers.LoadEnvVariables()
	initializers.ConnectToDatabase()
}

func main() {
	// Migrate the schema
	initializers.DB.AutoMigrate(&models.AccountAuth{})
	initializers.DB.AutoMigrate(&models.Account{})

	// You can add more models to migrate here if needed
	// initializers.DB.AutoMigrate(&models.AnotherModel{})

	// Print a message indicating migration is complete
	println("Database migration completed successfully.")
}
