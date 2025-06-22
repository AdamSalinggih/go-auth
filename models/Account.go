package models

import "gorm.io/gorm"

type Account struct {
	gorm.Model
	FirstName string `json:"first_name" gorm:"not null"`
	LastName  string `json:"last_name" gorm:"not null"`
	Email     string `json:"email" gorm:"not null;unique"`
	Address   string `json:"address" gorm:"not null"`
	Phone     string `json:"phone" gorm:"not null"`
	StateCode string `json:"state_code" gorm:"not null"`
	ZipCode   string `json:"zip_code" gorm:"not null"`
	Country   string `json:"country" gorm:"not null"`
}
