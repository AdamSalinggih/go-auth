package main

import (
	"github.com/adamhaiqal/go-auth/controllers"
	"github.com/adamhaiqal/go-auth/initializers"
	"github.com/gin-gonic/gin"
)

func init() {
	// Load environment variables and connect to the database
	initializers.LoadEnvVariables()
	initializers.ConnectToDatabase()
}
func main() {
	router := gin.Default()
	router.POST("/account/create", controllers.AccountCreate)
	router.GET("/account/get/:id", controllers.AccountGet)
	router.PUT("/account/update/:id", controllers.AccountUpdate)
	router.Run() // listen and serve on 0.0.0.0:8080
}
