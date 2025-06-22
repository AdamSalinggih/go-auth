package models

import gorm "gorm.io/gorm"

type AccountAuth struct {
	gorm.Model
	AccountID uint   `json:"account_id" gorm:"not null"`
	Username  string `json:"username" gorm:"not null;unique"`
	Password  string `json:"password" gorm:"not null"`
	Email     string `json:"email" gorm:"not null;unique"`
}
