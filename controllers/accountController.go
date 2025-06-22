package controllers

import (
	"github.com/adamhaiqal/go-auth/initializers"
	"github.com/adamhaiqal/go-auth/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func AccountCreate(c *gin.Context) {
	var account models.Account

	// Bind JSON and validate
	if err := c.BindJSON(&account); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}
	if err := validate.Struct(account); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	result := initializers.DB.Create(&account)
	if result.Error != nil {
		c.JSON(400, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(200, gin.H{
		"accountID": account.ID,
		"status":    "Success",
		"message":   "Account successfully created",
	})
}

func AccountGet(c *gin.Context) {
	id := c.Param("id")
	var account models.Account
	result := initializers.DB.First(&account, id)
	if result.Error != nil {
		c.Status(400)
		c.JSON(400, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(200, gin.H{
		"status": "Success",
		"data":   account,
	})
}

func AccountUpdate(c *gin.Context) {
	id := c.Param("id")
	var account struct {
		FirstName string
		LastName  string
		Email     string
		Address   string
		Phone     string
		StateCode string
		ZipCode   string
		Country   string
	}
	c.Bind(&account)
	var updateAccount models.Account
	result := initializers.DB.First(&updateAccount, id)
	if result.Error != nil {
		c.Status(400)
		c.JSON(400, gin.H{"error": result.Error.Error()})
		return
	}
	updateAccount.FirstName = account.FirstName
	updateAccount.LastName = account.LastName
	updateAccount.Email = account.Email
	updateAccount.Address = account.Address
	updateAccount.Phone = account.Phone
	updateAccount.StateCode = account.StateCode
	updateAccount.ZipCode = account.ZipCode
	updateAccount.Country = account.Country

	saveResult := initializers.DB.Save(&updateAccount)
	if saveResult.Error != nil {
		c.Status(400)
		c.JSON(400, gin.H{"error": saveResult.Error.Error()})
		return
	}
	c.JSON(200, gin.H{
		"status":  "Success",
		"message": "Account successfully updated",
	})
}
