package controllers

import (
	"github.com/adamhaiqal/go-auth/initializers"
	"github.com/adamhaiqal/go-auth/models"
	"github.com/gin-gonic/gin"
)

func AccountCreate(c *gin.Context) {
	var account models.Account

	c.BindJSON(&account)
	if err := initializers.DB.Create(&account).Error; err != nil {
		c.JSON(400, gin.H{"error": "Account creation failed"})
		return
	}

	c.JSON(200, gin.H{"message": "Account created successfully", "account": account})
}
