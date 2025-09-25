package product

import "time"

type Brand struct {
	ID        int64  `gorm:"primaryKey;autoIncrement"`
	Name      string `gorm:"size:255;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`
}

type Category struct {
	ID        int64  `gorm:"primaryKey;autoIncrement"`
	Name      string `gorm:"size:255;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`
}

type Product struct {
	ID        int64  `gorm:"primaryKey;autoIncrement"`
	BrandID   int64  `gorm:"not null;index"`
	BaseName  string `gorm:"size:255;not null"`
	BaseDesc  string `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`
}

type ProductVariant struct {
	ID                   int64  `gorm:"primaryKey;autoIncrement"`
	ProductID            int64  `gorm:"not null;index"`
	VariantNameExtension string `gorm:"size:255"`
	Price                float64
	StockQuantity        int
}

type ProductImage struct {
	ID        int64  `gorm:"primaryKey;autoIncrement"`
	VariantID int64  `gorm:"not null;index"`
	ImageURL  string `gorm:"size:255;not null"`
	AltText   string `gorm:"size:255"`
	Order     int
}

type Option struct {
	ID   int64  `gorm:"primaryKey;autoIncrement"`
	Name string `gorm:"size:255;not null"`
}

type OptionValue struct {
	ID       int64  `gorm:"primaryKey;autoIncrement"`
	OptionID int64  `gorm:"not null;index"`
	Name     string `gorm:"size:255;not null"`
}

type VariantOption struct {
	VariantID int64 `gorm:"primaryKey;autoIncrement:false"`
	ValueID   int64 `gorm:"primaryKey;autoIncrement:false"`
}

type ProductCategory struct {
	ProductID  int64 `gorm:"primaryKey;autoIncrement:false"`
	CategoryID int64 `gorm:"primaryKey;autoIncrement:false"`
}
