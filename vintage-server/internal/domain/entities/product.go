package entities

import (
	"errors"
	"time"
)

type Product struct {
	ID          uint64    `json:"id" db:"id"`
	ShopID      uint64    `json:"shop_id" db:"shop_id"`
	ConditionID uint8     `json:"condition_id" db:"condition_id"`
	CategoryID  uint32    `json:"category_id" db:"category_id"`
	BrandID     *uint32   `json:"brand_id" db:"brand_id"`
	Name        string    `json:"name" db:"name"`
	Summary     *string   `json:"summary" db:"summary"`
	Description *string   `json:"description" db:"description"`
	Price       uint64    `json:"price" db:"price"`
	Size        *string   `json:"size" db:"size"`
	Stock       uint32    `json:"stock" db:"stock"`
	IsLatest    bool      `json:"is_latest" db:"is_latest"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

func (p *Product) Validate() error {
	if len(p.Name) < 2 {
		return errors.New("product name must be at least 2 characters")
	}
	if p.Price == 0 {
		return errors.New("price must be greater than 0")
	}
	if p.Stock <= 0 {
		return errors.New("stock cannot be negative")
	}
	return nil
}

func (p *Product) ReduceStock(quantity uint32) error {
	if quantity > p.Stock {
		return errors.New("insufficient stock")
	}
	p.Stock -= quantity
	return nil
}

func (p *Product) IncreaseStock(quantity uint32) {
	p.Stock += quantity
}

func (p *Product) IsAvailable() bool {
	return p.Stock > 0 && p.IsLatest
}

func (p *Product) MarkAsLatest() {
	p.IsLatest = true
}

func (p *Product) MarkAsArchived() {
	p.IsLatest = false
}
