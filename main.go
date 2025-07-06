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
	router.LoadHTMLGlob("templates/*")
	// Set up BasicAuth credentials using a map
	accounts := gin.Accounts{
		"admin": "admin123",
	}

	router.GET("/", controllers.Welcome)

	api := router.Group("/api/v1/account", gin.BasicAuth(accounts))
	{
		api.POST("/create", controllers.AccountCreate)
		api.GET("/get/:id", controllers.AccountGet)
		api.PUT("/update/:id", controllers.AccountUpdate)
		api.DELETE("/delete/:id", controllers.AccountDelete)
	}

	router.Run()
}
