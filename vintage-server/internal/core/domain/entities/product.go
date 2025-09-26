// domain/entities/product.go
package entities

import "time"

type Product struct {
	ID          uint64
	ShopID      uint64
	ConditionID uint8
	CategoryID  uint32
	BrandID     *uint32
	Name        string
	Summary     string
	Description string
	Price       uint64
	Size        string
	Stock       uint32
	IsLatest    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ProductImage struct {
	ID        uint64
	ProductID uint64
	Index     uint8
	URL       string
	CreatedAt time.Time
}
