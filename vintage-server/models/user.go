package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username       string `gorm:"size:100;uniqueIndex;not null"`
	Password       string `gorm:"not null"`
	Fullname       string `gorm:"not null"`
	Email          string `gorm:"size:100;uniqueIndex;not null"`
	ProfilePicture string `gorm:"size:255"`
	Addresses      []Address
	Carts          []Cart
	Transactions   []Transaction
	Wishlists      []Wishlist
}

type Address struct {
	gorm.Model
	UserID      uint
	User        User
	FullAddress string `gorm:"type:text"`
	Kel         int
	Kec         int
	Kab         int
	Prov        int
	KodePos     string `gorm:"size:10"`
}
