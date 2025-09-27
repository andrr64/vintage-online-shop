package model

import "time"

// ProductCondition merepresentasikan tabel 'product_conditions'
type ProductCondition struct {
	ID        int16     `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// ProductCategory merepresentasikan tabel 'product_categories'
type ProductCategory struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
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
	ID       int    `json:"id" db:"id"`
	SizeName string `json:"size_name" db:"size_name"`
}

// Shop merepresentasikan tabel 'shop'
type Shop struct {
	ID          int64     `json:"id" db:"id"`
	AccountID   int64     `json:"account_id" db:"account_id"`
	Name        string    `json:"name" db:"name"`
	Summary     *string   `json:"summary" db:"summary"`
	Description *string   `json:"description" db:"description"`
	Active      bool      `json:"active" db:"active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Product merepresentasikan tabel 'products'
type Product struct {
	ID          int64     `json:"id" db:"id"`
	ShopID      int64     `json:"shop_id" db:"shop_id"`
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
}

// ProductImage merepresentasikan tabel 'product_images'
type ProductImage struct {
	ID         int64     `json:"id" db:"id"`
	ProductID  int64     `json:"product_id" db:"product_id"`
	ImageIndex int16     `json:"image_index" db:"image_index"`
	URL        string    `json:"url" db:"url"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

// Review merepresentasikan tabel 'reviews'
type Review struct {
	ID        int64     `json:"id" db:"id"`
	ProductID int64     `json:"product_id" db:"product_id"`
	AccountID int64     `json:"account_id" db:"account_id"`
	OrderID   int64     `json:"order_id" db:"order_id"`
	Rating    int16     `json:"rating" db:"rating"`
	Comment   *string   `json:"comment" db:"comment"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
