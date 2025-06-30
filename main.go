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
	router.POST("api/v1/account/create", controllers.AccountCreate)
	router.GET("api/v1/account/get/:id", controllers.AccountGet)
	//router.PUT("api/v1/account/update/:id", controllers.AccountUpdate)
	//router.DELETE("api/v1/account/delete/:id", controllers.AccountDelete)
	router.Run()
}
