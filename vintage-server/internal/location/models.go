package location

import (
	"time"
	"gorm.io/gorm"
)

type Provinsi struct {
	ID             uint           `gorm:"primaryKey"`
	Name           string
	KabupatenKota  []KabupatenKota
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

type KabupatenKota struct {
	ID         uint      `gorm:"primaryKey"`
	ProvinsiID uint
	Name       string
	Kecamatan  []Kecamatan
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

type Kecamatan struct {
	ID             uint   `gorm:"primaryKey"`
	KabupatenKotaID uint
	Name           string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

type Address struct {
	ID            uint `gorm:"primaryKey"`
	UserID        uint
	KecamatanID   uint
	Label         string
	RecipientName string
	PhoneNumber   string
	FullAddress   string
	PostalCode    string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}
