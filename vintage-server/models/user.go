package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username       string `gorm:"size:100;uniqueIndex"`
	Fullname       string `gorm:"size:100"`
	Password       string
	Email          string `gorm:"size:150;uniqueIndex"`
	ProfilePicture string
	Alamat      []Alamat `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
