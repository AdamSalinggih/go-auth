package main

import (
	"fmt"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading godotevent:", err)
		return
	}
	fmt.Println("godotevent loaded successfully")
}
func main() {
	fmt.Println("Hello, World!")
}
