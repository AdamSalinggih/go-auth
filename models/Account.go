package models

import "gorm.io/gorm"

type Account struct {
	gorm.Model
	Username    string `json:"username" gorm:"not null;unique"`
	Password    string `json:"password" gorm:"not null"`
	Email       string `json:"email" gorm:"not null;unique"`
	FirstName   string `json:"first_name" gorm:"not null"`
	LastName    string `json:"last_name" gorm:"not null"`
	Address     string `json:"address" gorm:"not null"`
	HomePhone   string `json:"home_phone" gorm:"not null"`
	MobilePhone string `json:"mobile_phone" gorm:"not null"`
	WorkPhone   string `json:"work_phone" gorm:"not null"`
	StateCode   string `json:"state_code" gorm:"not null"`
	ZipCode     string `json:"zip_code" gorm:"not null"`
	Country     string `json:"country" gorm:"not null"`
	IsVerified  bool   `json:"is_verified" gorm:"default:false"`
}
