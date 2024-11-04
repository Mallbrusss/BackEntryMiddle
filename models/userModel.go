package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Token    string `json:"token" gorm:"uniqueIndex;size:255"`
	Login    string `json:"login" gorm:"uniqueIndex;size:255;not null"`
	Password string `json:"pswd" gorm:"size:255;not null"`
}
