package models

import (
	"time"
)

// Document - модель документа 
type Document struct {
	// gorm.Model
	ID        string `gorm:"primaryKey;size:255;index"`
	Name      string `json:"name" gorm:"size:255;not null"`
	Mime      string `json:"mime" gorm:"size:100;not null"`
	FilePath  string `json:"-" gorm:"size:500;not null"`
	File      bool   `json:"file" gorm:"not null"`
	Public    bool   `json:"public" gorm:"not null"`
	Token     string `json:"-" gorm:"-"`
	CreatedAt time.Time
	Grant     []string `json:"grant" gorm:"-"`
}

// CacheDocument - модель документа для кеширования
type CacheDocument struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Mime      string    `json:"mime"`
	FilePath  string    `json:"file_path"`
	File      bool      `json:"file"`
	Public    bool      `json:"public"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
	Grant     []string  `json:"grant"`
}

// DocumentAccess - модель юзеров имеющих доступ
type DocumentAccess struct {
	DocID string `json:"-" gorm:"size:255;not null;index"`
	Login string `json:"login" gorm:"size:255;not null;index"`
}

// Meta - модель мета данных
type Meta struct {
	Name   string   `json:"name"`
	File   bool     `json:"file"`
	Public bool     `json:"public"`
	Token  string   `json:"token"`
	Mime   string   `json:"mime"`
	Grant  []string `json:"grant"`
}
