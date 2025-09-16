package models

import (
	"gorm.io/gorm"
)

type Wishlist struct {
	gorm.Model
	UserID   uint
	User     User
	ProdukID uint
	Produk   Produk
}
