package models

import "gorm.io/gorm"

type AccountMessage struct {
	gorm.Model
	AccountID uint   `json:"account_id" gorm:"not null"`
	Message   string `json:"message" gorm:"not null;size:500"`
	IsRead    bool   `json:"is_read" gorm:"default:false"`
	CreatedAt string `json:"created_at" gorm:"not null"`
	UpdatedAt string `json:"updated_at" gorm:"not null"`
}
