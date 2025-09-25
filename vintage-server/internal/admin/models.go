package admin

import (
	"time"
)

type Admin struct {
	ID           uint `gorm:"primaryKey"`
	Username     string
	PasswordHash string
	Logs         []AdminLog
}

type AdminLog struct {
	ID        uint `gorm:"primaryKey"`
	AdminID   uint
	Action    string
	Timestamp time.Time
}
