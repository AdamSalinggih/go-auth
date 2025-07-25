package main

import (
	"github.com/adamhaiqal/go-auth/controllers"
	"github.com/adamhaiqal/go-auth/initializers"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDatabase()
}

func main() {
	router := gin.Default()

	api := router.Group("/api/v1/account")
	{
		api.POST("/signup", controllers.AccountSignup)
		api.POST("/signin/:id", controllers.AccountSignin)
	}

	router.Run()
}
