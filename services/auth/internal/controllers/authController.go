package controllers

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/adamhaiqal/go-auth/pkg/models"
	"github.com/adamhaiqal/go-auth/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthController struct {
	DB *gorm.DB
}

func NewAuthController(db *gorm.DB) *AuthController {
	return &AuthController{DB: db}
}

func (a *AuthController) Register(c *gin.Context) {
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
	usernameErr := a.DB.Where("username = ?", account.Username).First(&existingAccount).Error
	emailErr := a.DB.Where("email = ?", account.Email).First(&existingAccount).Error

	if usernameErr == nil || emailErr == nil {
		c.JSON(400, gin.H{"error": "Unable to create account"})
		return
	}

	if !errors.Is(usernameErr, gorm.ErrRecordNotFound) {
		log.Printf("Database error checking username: %v", usernameErr)
	}
	if !errors.Is(emailErr, gorm.ErrRecordNotFound) {
		log.Printf("Database error checking email: %v", emailErr)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to process password"})
		return
	}
	account.Password = string(hashedPassword)

	account.IsVerified = false

	if err := a.DB.Create(&account).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to create account"})
		return
	}

	c.JSON(201, gin.H{
		"message":  "Account created successfully",
		"username": account.Username,
		"email":    account.Email,
	})
}

func (a *AuthController) Login(c *gin.Context) {
	var signinRequest struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.BindJSON(&signinRequest); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	var account models.Account

	err := a.DB.Where("username = ?", signinRequest.Username).First(&account).Error
	if err != nil {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
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

	tokenString, err := token.SignedString(utils.GetJWTKey())

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate token"})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("token", tokenString, 3600, "", "", false, true)
	c.JSON(200, gin.H{"message": "Login successful"})
}

func (a *AuthController) GetMe(c *gin.Context) {
	accountID, exists := c.Get("accountID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var account models.Account
	if err := a.DB.First(&account, accountID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Resource not found"})
		} else {
			log.Printf("Database error in GetMe: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
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

func Logout(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("token", "", -1, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}
