package models

import (
	"time"
)

type AccountMessage struct {
	ID          uint64     `gorm:"primaryKey;autoIncrement"`
	SenderID    uint64     `gorm:"not null;index"`
	RecipientID uint64     `gorm:"not null;index"`
	Body        string     `gorm:"type:text;not null"`
	SentAt      time.Time  `gorm:"autoCreateTime"`
	ReadAt      *time.Time `gorm:"index"`
	EditedAt    *time.Time
	IsDeleted   bool    `gorm:"default:false"`
	ThreadID    *string `gorm:"type:uuid;index"`
}
