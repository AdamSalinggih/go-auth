package controllers

import (
	"github.com/adamhaiqal/go-auth/initializers"
	"github.com/adamhaiqal/go-auth/models"
	"github.com/gin-gonic/gin"
)

func AccountCreate(c *gin.Context) {

	var account models.Account

	err := c.BindJSON(&account)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := initializers.DB.Create(&account).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Account created successfully", "account": account})
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
