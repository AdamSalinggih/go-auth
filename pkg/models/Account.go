package models

import "gorm.io/gorm"

type Account struct {
	gorm.Model
	Username    string `json:"username" gorm:"not null;unique" validate:"required,min=3,max=50,alphanum"`
	Password    string `json:"password" gorm:"not null" validate:"required,min=8"`
	Email       string `json:"email" gorm:"not null;unique" validate:"required,email"`
	FirstName   string `json:"first_name" gorm:"not null" validate:"required,min=2,max=50"`
	LastName    string `json:"last_name" gorm:"not null" validate:"required,min=2,max=50"`
	Address     string `json:"address" gorm:"not null" validate:"required,min=5,max=200"`
	HomePhone   string `json:"home_phone" gorm:"not null" validate:"required,min=10,max=15"`
	MobilePhone string `json:"mobile_phone" gorm:"not null" validate:"required,min=10,max=15"`
	WorkPhone   string `json:"work_phone" gorm:"not null" validate:"required,min=10,max=15"`
	StateCode   string `json:"state_code" gorm:"not null" validate:"required,min=2,max=10"`
	ZipCode     string `json:"zip_code" gorm:"not null" validate:"required,min=3,max=10"`
	Country     string `json:"country" gorm:"not null" validate:"required,min=2,max=50"`
	IsVerified  bool   `json:"is_verified" gorm:"default:false"`
}
