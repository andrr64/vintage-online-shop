package admin

import "time"

type Admin struct {
	ID           int64  `gorm:"primaryKey;autoIncrement"`
	Username     string `gorm:"size:255;not null"`
	PasswordHash string `gorm:"size:255;not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type AdminLog struct {
	ID      int64  `gorm:"primaryKey;autoIncrement"`
	AdminID int64  `gorm:"not null;index"`
	Action  string `gorm:"type:text;not null"`
	Time    time.Time
}
