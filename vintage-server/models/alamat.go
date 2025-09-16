package models

import (
	"gorm.io/gorm"
)

type Alamat struct {
	gorm.Model
	UserID        uint `gorm:"index" json:"user_id"`
	AlamatLengkap string
	Kel           int
	Kec           int
	Kab           int
	Prov          int
	KodePos       string
	IsPrimary     bool `gorm:"default:false" json:"is_primary"` // menandai alamat utama
}