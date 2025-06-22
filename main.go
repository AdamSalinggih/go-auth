package main

import (
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
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "LOL, it works!",
		})
	})
	router.Run() // listen and serve on 0.0.0.0:8080
}
