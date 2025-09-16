package models

import (
	"time"

	"gorm.io/gorm"
)

type AdminLog struct {
	gorm.Model
	AdminID   uint
	Action    string
	Timestamp time.Time
}
