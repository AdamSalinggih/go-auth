package utils

import (
	"log"
	"os"
)

func GetJWTKey() []byte {
	key := os.Getenv("SIGNIN_KEY")
	if key == "" {
		log.Fatal("SIGNIN_KEY environment variable is not set. Please set it to a secure secret.")
	}
	return []byte(key)
}
