package main

import (
	"log"

	"github.com/adamhaiqal/go-auth/pkg/config"
	"github.com/adamhaiqal/go-auth/services/auth/internal/controllers"
	"github.com/adamhaiqal/go-auth/services/auth/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}
	config.ConnectToDatabase()
}

func main() {
	router := gin.Default()

	// Authentication routes (public)
	auth := router.Group("/api/v1/auth")
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
		auth.POST("/logout", middleware.AuthenticateCookie, controllers.Logout)
	}

	// Account routes (protected)
	account := router.Group("/api/v1/account")
	account.Use(middleware.AuthenticateCookie)
	{
		account.GET("/me", controllers.GetMe)
	}

	router.Run()
}
