package controllers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/adamhaiqal/go-auth/initializers"
	"github.com/adamhaiqal/go-auth/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func AccountSignup(c *gin.Context) {
	var account models.Account

	err := c.BindJSON(&account)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON input"})
		return
	}

	validate := validator.New()
	if err := validate.Struct(account); err != nil {
		c.JSON(400, gin.H{"error": "Validation failed", "details": err.Error()})
		return
	}

	var existingAccount models.Account
	if err := initializers.DB.Where("username = ?", account.Username).First(&existingAccount).Error; err == nil {
		c.JSON(400, gin.H{"error": "Username already exists"})
		return
	}

	if err := initializers.DB.Where("email = ?", account.Email).First(&existingAccount).Error; err == nil {
		c.JSON(400, gin.H{"error": "Email already exists"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to process password"})
		return
	}
	account.Password = string(hashedPassword)

	account.IsVerified = false

	if err := initializers.DB.Create(&account).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to create account"})
		return
	}

	c.JSON(201, gin.H{
		"message":  "Account created successfully",
		"username": account.Username,
		"email":    account.Email,
	})
}

func AccountSignin(c *gin.Context) {
	var signinRequest struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.BindJSON(&signinRequest); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	var account models.Account
	if err := initializers.DB.Where("username = ?", signinRequest.Username).First(&account).Error; err != nil {
		c.JSON(404, gin.H{"error": "Account not found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(signinRequest.Password)); err != nil {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	fmt.Println("After password verification, setting cookie")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": account.Username,
		"email":    account.Email,
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SIGNIN_KEY"))) // Replace with your secret key

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate token"})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("token", tokenString, 3600, "", "", false, true)
	c.JSON(200, gin.H{"message": "Signin successful"})
}
