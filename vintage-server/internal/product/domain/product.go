package domain

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID          uuid.UUID
	ShopID      uuid.UUID
	ConditionID int16
	CategoryID  int
	BrandID     *int
	SizeID      *int
	Name        string
	Summary     *string
	Description *string
	Price       int64
	Stock       int
	IsLatest    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time

	// Relasi
	Images    []ProductImage
	Brand     *Brand
	Category  *ProductCategory
	Condition *ProductCondition
	Size      *ProductSize
	Shop      *Shop
}
