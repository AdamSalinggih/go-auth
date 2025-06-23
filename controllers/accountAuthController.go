package controllers

import (
	"github.com/adamhaiqal/go-auth/initializers"
	"github.com/adamhaiqal/go-auth/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func AccountAuthSignup(c *gin.Context) {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	if c.Bind(&body) != nil {
		c.Status(400)
		c.JSON(400, gin.H{"error": "Failed to read body"})
		return
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.Status(400)
		c.JSON(400, gin.H{"error": "Failed to hash password"})
		return
	}

	account := models.AccountAuth{
		Username: body.Username,
		Password: string(hashPassword),
		Email:    body.Email,
	}

	result := initializers.DB.Create(&account)
	if result.Error != nil {
		c.Status(400)
		c.JSON(400, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(200, gin.H{
		"data":    account,
		"status":  "Success",
		"message": "Account successfully created",
	})
}
