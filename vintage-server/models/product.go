package models

import (
	"time"

	"gorm.io/gorm"
)

//
// ===== Master Data =====
//

type Brand struct {
	gorm.Model
	Name     string `gorm:"size:255;not null"`
	Products []Product
}

type Category struct {
	gorm.Model
	Name     string `gorm:"size:255;not null"`
	Products []Product
}

//
// ===== Product & Versioning =====
//

type Product struct {
	gorm.Model
	Name       string `gorm:"size:255;not null"`
	BrandID    uint   `gorm:"not null"`
	Brand      Brand
	CategoryID uint `gorm:"not null"`
	Category   Category
	Versions   []ProductVersion
}

type ProductVersion struct {
	gorm.Model
	ProductID         uint `gorm:"index;not null"`
	Product           Product
	PreviousVersionID *uint
	PreviousVersion   *ProductVersion `gorm:"foreignKey:PreviousVersionID"`

	Size        string  `gorm:"size:50"`
	Color       string  `gorm:"size:50"`
	Description string  `gorm:"type:text"`
	Price       float64 `gorm:"not null"`
	Stock       int     `gorm:"not null"`

	ActiveFrom *time.Time
	ActiveTo   *time.Time

	CartDetails      []CartDetail        `gorm:"foreignKey:VersionID"`
	WishlistItems    []Wishlist          `gorm:"foreignKey:VersionID"`
	TransactionItems []TransactionDetail `gorm:"foreignKey:VersionID"`
}

//
// ===== User Related =====
//

type Wishlist struct {
	gorm.Model
	UserID    uint
	User      User
	VersionID uint
	Version   ProductVersion
}

type Cart struct {
	gorm.Model
	UserID uint
	User   User
	Items  []CartDetail
}

type CartDetail struct {
	gorm.Model
	CartID    uint
	Cart      Cart
	VersionID uint
	Version   ProductVersion
	Quantity  int
}

//
// ===== Transactions =====
//

type Transaction struct {
	gorm.Model
	UserID        uint
	User          User
	TotalAmount   float64
	TotalQuantity int
	Status        string `gorm:"size:50"`
	Details       []TransactionDetail
}

type TransactionDetail struct {
	gorm.Model
	TransactionID uint
	Transaction   Transaction
	VersionID     uint
	Version       ProductVersion
	Quantity      int
	Total         float64
}
