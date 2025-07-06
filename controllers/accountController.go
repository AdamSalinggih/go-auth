package controllers

import (
	"strings"

	"github.com/adamhaiqal/go-auth/initializers"
	"github.com/adamhaiqal/go-auth/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Welcome(c *gin.Context) {
	c.HTML(200, "welcome.html", gin.H{
		"title": "Welcome to Go Auth",
	})
}

func AccountCreate(c *gin.Context) {

	var account models.Account

	err := c.BindJSON(&account)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(400, gin.H{"error": "Failed to hash password"})
		return
	}
	account.Password = string(hashedPassword)

	if err := initializers.DB.Create(&account).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			c.JSON(400, gin.H{"error": "Username already exists"})
			return
		}
		c.JSON(400, gin.H{"error": "Failed to create account", "details": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Account created successfully", "accountId": account.ID, "username": account.Username})
}

func AccountGet(c *gin.Context) {
	var account models.Account

	id := c.Param("id")

	if err := initializers.DB.First(&account, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Account not found"})
		return
	}

	c.JSON(200, gin.H{"account": account})
}

func AccountUpdate(c *gin.Context) {
	var account models.Account

	id := c.Param("id")

	if err := initializers.DB.First(&account, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Account not found"})
		return
	}

	if err := c.BindJSON(&account); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := initializers.DB.Save(&account).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Account updated successfully", "account": account})
}

func AccountDelete(c *gin.Context) {
	var account models.Account

	id := c.Param("id")

	if err := initializers.DB.First(&account, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Account not found"})
		return
	}

	if err := initializers.DB.Unscoped().Delete(&account).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"accountId": account.ID,
		"message":   "Account deleted successfully"})
}
