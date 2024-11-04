package models

import (
	"time"

	"gorm.io/gorm"
)

type Document struct {
	gorm.Model
	ID        string    `json:"id" gorm:"size:255;not null"`
	Name      string    `json:"name" gorm:"size:255;not null"`
	Mime      string    `json:"mime" gorm:"size:100;not null"`
	FilePath  string    `json:"-" gorm:"size:500;not null"`
	File      bool      `json:"file" gorm:"not null;default:false"`
	Public    bool      `json:"public" gorm:"not null;default:false"`
	Token     string    `json:"token" gorm:"-"`
	CreatedAt time.Time `json:"created"`
	UpdatedAt time.Time
}


type DocumentAccess struct {
	gorm.Model
	PkID uint `gorm:"primaryKey"`
	ID string `gorm:"not null;index"`
	Login      string `gorm:"size:255;not null;index"`
}

type Meta struct {
	Name   string   `json:"name"`
	File   bool     `json:"file"`
	Public bool     `json:"public"`
	Token  string   `json:"token"`
	Mime   string   `json:"mime"`
	Grant  []string `json:"grant"`
}
