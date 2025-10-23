package model

import (
	"time"

	"github.com/google/uuid"
)

// ProductCondition merepresentasikan tabel 'product_conditions'
type ProductCondition struct {
	ID        int16     `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// ProductCategory merepresentasikan tabel 'product_categories'
type ProductCategory struct {
	ID           int       `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	ProductCount *int64    `json:"product_count" db:"product_count"`
}

// Brand merepresentasikan tabel 'brands'
type Brand struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	LogoURL   *string   `json:"logo_url" db:"logo_url"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// ProductSize merepresentasikan tabel 'product_size'
type ProductSize struct {
	ID   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

// Shop merepresentasikan tabel 'shop'
type Shop struct {
	ID          uuid.UUID `json:"id" db:"id"`
	AccountID   uuid.UUID `json:"account_id" db:"account_id"`
	Name        string    `json:"name" db:"name"`
	Summary     *string   `json:"summary" db:"summary"`
	Description *string   `json:"description" db:"description"`
	Active      bool      `json:"active" db:"active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Product merepresentasikan tabel 'products'
type Product struct {
	ID          uuid.UUID `json:"id" db:"id"`
	ShopID      uuid.UUID `json:"shop_id" db:"shop_id"`
	ConditionID int16     `json:"condition_id" db:"condition_id"`
	CategoryID  int       `json:"category_id" db:"category_id"`
	BrandID     *int      `json:"brand_id" db:"brand_id"`
	SizeID      *int      `json:"size_id" db:"size_id"`
	Name        string    `json:"name" db:"name"`
	Summary     *string   `json:"summary" db:"summary"`
	Description *string   `json:"description" db:"description"`
	Price       int64     `json:"price" db:"price"`
	Stock       int       `json:"stock" db:"stock"`
	IsLatest    bool      `json:"is_latest" db:"is_latest"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`

	// 🔽 Relasi (bukan dari DB langsung)
	Images    []ProductImage    `json:"images,omitempty" db:"-"`
	Brand     *Brand            `json:"brand,omitempty" db:"-"`
	Category  *ProductCategory  `json:"category,omitempty" db:"-"`
	Condition *ProductCondition `json:"condition,omitempty" db:"-"`
	Size      *ProductSize      `json:"size,omitempty" db:"-"`
	Shop      *Shop             `json:"shop,omitempty" db:"-"`
}

// ProductImage merepresentasikan tabel 'product_images'
type ProductImage struct {
	ID         int64     `json:"id" db:"id"`
	ProductID  uuid.UUID `json:"product_id" db:"product_id"`
	ImageIndex int16     `json:"image_index" db:"image_index"`
	URL        string    `json:"url" db:"url"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

// Review merepresentasikan tabel 'reviews'
type Review struct {
	ID        uuid.UUID `json:"id" db:"id"`
	ProductID int64     `json:"product_id" db:"product_id"`
	AccountID int64     `json:"account_id" db:"account_id"`
	OrderID   int64     `json:"order_id" db:"order_id"`
	Rating    int16     `json:"rating" db:"rating"`
	Comment   *string   `json:"comment" db:"comment"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
