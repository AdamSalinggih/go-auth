package controllers

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/adamhaiqal/go-auth/config"
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
	if err := config.DB.Where("username = ?", account.Username).First(&existingAccount).Error; err == nil {
		c.JSON(400, gin.H{"error": "Username already exists"})
		return
	}

	if err := config.DB.Where("email = ?", account.Email).First(&existingAccount).Error; err == nil {
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

	if err := config.DB.Create(&account).Error; err != nil {
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

	if err := config.DB.Where("username = ?", signinRequest.Username).First(&account).Error; err != nil {
		c.JSON(404, gin.H{"error": "Account not found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(signinRequest.Password)); err != nil {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":      account.ID,
		"username": account.Username,
		"email":    account.Email,
		"exp":      time.Now().Add(time.Minute * 15).Unix(),
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

func AccountHome(c *gin.Context) {
	accountID, exists := c.Get("accountID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var account models.Account
	if err := config.DB.First(&account, accountID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Account not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Welcome to your account",
		"username":  account.Username,
		"email":     account.Email,
		"verified":  account.IsVerified,
		"address":   account.Address,
		"state":     account.StateCode,
		"zip":       account.ZipCode,
		"country":   account.Country,
		"firstname": account.FirstName,
		"lastname":  account.LastName,
		"status":    http.StatusOK,
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

func AccountSignout(c *gin.Context) {
	log.Print("Inside AccountSignout controller")
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("token", "", -1, "", "", false, true) // Clear the cookie
	c.JSON(http.StatusOK, gin.H{"message": "Signout successful"})
}
