package models

import (
	"time"

	"gorm.io/gorm"
)

type Document struct {
	gorm.Model
	Name      string `gorm:"size:255;not null"`
	MimeType  string `gorm:"size:100;not null"`
	FilePath  string `gorm:"size:500;not null"`
	Public    bool   `gorm:"not null;default:false"`
	OwnerID   uint   `gorm:"not null"` // ID владельца документа
	CreatedAt time.Time
	UpdatedAt time.Time
	Grant     []string `gorm:"-"` // Пользователи с доступом
}
