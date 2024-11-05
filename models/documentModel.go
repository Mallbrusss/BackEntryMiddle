package models

import (
	"time"
)

type Document struct {
	// gorm.Model
	ID        string `gorm:"primaryKey;size:255"`
	Name      string `json:"name" gorm:"size:255;not null"`
	Mime      string `json:"mime" gorm:"size:100;not null"`
	FilePath  string `json:"-" gorm:"size:500;not null"`
	File      bool   `json:"file" gorm:"not null"`
	Public    bool   `json:"public" gorm:"not null"`
	Token     string `json:"-" gorm:"-"`
	CreatedAt time.Time
	Grant     []string   `json:"grant" gorm:"-"`
}

type DocumentAccess struct {
	// gorm.Model
	DocID string `json:"-" gorm:"size:255;not null;index"`
	Login string `json:"login" gorm:"size:255;not null;index"`
}

type Meta struct {
	Name   string   `json:"name"`
	File   bool     `json:"file"`
	Public bool     `json:"public"`
	Token  string   `json:"token"`
	Mime   string   `json:"mime"`
	Grant  []string `json:"grant"`
}
