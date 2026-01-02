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

	api := router.Group("/api/v1/account")
	{
		api.POST("/signup", controllers.AccountSignup)
		api.POST("/signin/:id", controllers.AccountSignin)
		api.POST("/signout", middleware.AuthenticateCookie, controllers.AccountSignout)
		api.GET("/home", middleware.AuthenticateCookie, controllers.AccountHome)
	}

	router.Run()
}
