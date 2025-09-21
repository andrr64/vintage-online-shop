package models

import "gorm.io/gorm"

type Admin struct {
	gorm.Model
	Username string `gorm:"size:100;uniqueIndex;not null"`
	Password string `gorm:"not null"`
	Logs     []AdminLog
}

type AdminLog struct {
	gorm.Model
	AdminID uint   `gorm:"not null"`
	Action  string `gorm:"size:255"`
}
