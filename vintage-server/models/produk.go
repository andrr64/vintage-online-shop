package models

import (
	"gorm.io/gorm"
	"time"
)

type Produk struct {
	gorm.Model
	NamaProduk  string
	BrandID     uint
	Brand       Brand
	Size        string
	Warna       string
	UploadedAt  time.Time
	LastUpdate  time.Time
	KategoriID  uint
	Kategori    Kategori
	Description string
	Harga       float64
	Stok        int
	CartDetails []CartDetail
	Wishlists   []Wishlist
}
