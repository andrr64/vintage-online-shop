package models

import (
	"gorm.io/gorm"
)

type Kategori struct {
	gorm.Model
	NamaKategori string
	Produk       []Produk
}