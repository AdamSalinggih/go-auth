package main

import (
	"log"

	"github.com/adamhaiqal/go-auth/pkg/config"
	"github.com/adamhaiqal/go-auth/pkg/database"
	"github.com/adamhaiqal/go-auth/services/auth/internal/controllers"
	"github.com/adamhaiqal/go-auth/services/auth/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}
}

func main() {
	router := gin.Default()
	db, err := database.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	authController := controllers.NewAuthController(db)
	config.GoogleOauthConfig()

	auth := router.Group("/api/v1/auth")
	{
		auth.POST("/register", authController.Register)
		auth.POST("/login", authController.Login)
		auth.POST("/logout", middleware.AuthenticateCookie, controllers.Logout)
	}

	router.Run()
}
