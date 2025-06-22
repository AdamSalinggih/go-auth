package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/adamhaiqal/go-auth/initializers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	// Load environment variables and connect to the database
	initializers.LoadEnvVariables()
	initializers.ConnectToDatabase()
}
func main() {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	//use ../.env because main.go inside /cmd
	err = godotenv.Load(filepath.Join(pwd, "../.env"))
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	fmt.Println(fmt.Sprintf("MYVAR=%s", os.Getenv("MYVAR")))

	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "LOL, it works!",
		})
	})
	router.Run() // listen and serve on 0.0.0.0:8080
}
