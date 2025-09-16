package models

import (
	"gorm.io/gorm"
)

type Admin struct {
	gorm.Model
	Username string `gorm:"size:100;uniqueIndex"`
	Password string
	Logs     []AdminLog
}
