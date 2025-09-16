package models

import (
	"gorm.io/gorm"
)

type Brand struct {
	gorm.Model
	NamaBrand string
	Produk    []Produk
}
