package product

import (
	"gorm.io/gorm"
	"time"
)

type Brand struct {
	ID        uint `gorm:"primaryKey"`
	Name      string
	Products  []Product
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Category struct {
	ID        uint `gorm:"primaryKey"`
	Name      string
	Products  []Product `gorm:"many2many:product_categories;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Product struct {
	ID              uint `gorm:"primaryKey"`
	BrandID         uint
	BaseName        string
	BaseDescription string
	Variants        []ProductVariant
	Categories      []Category `gorm:"many2many:product_categories;"`
	Reviews         []Review
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}

type ProductVariant struct {
	ID                   uint `gorm:"primaryKey"`
	ProductID            uint
	VariantNameExtension string
	Price                float64
	StockQuantity        int
	Images               []ProductImage
	VariantOptions       []VariantOption
}

type ProductImage struct {
	ID               uint `gorm:"primaryKey"`
	ProductVariantID uint
	ImageURL         string
	AltText          string
	DisplayOrder     int
}

type Option struct {
	ID     uint `gorm:"primaryKey"`
	Name   string
	Values []OptionValue
}

type OptionValue struct {
	ID       uint `gorm:"primaryKey"`
	OptionID uint
	Name     string
}

type VariantOption struct {
	ProductVariantID uint `gorm:"primaryKey;autoIncrement:false"`
	OptionValueID    uint `gorm:"primaryKey;autoIncrement:false"`
}

type Review struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint
	ProductID uint
	Rating    int
	Comment   string
	CreatedAt time.Time
}
