package model

import (
	"time"

	"github.com/google/uuid"
)

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
}
