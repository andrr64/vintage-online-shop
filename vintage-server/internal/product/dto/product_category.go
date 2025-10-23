package dto

import "time"

type ProductCategoryBase struct {
	ID           int        `json:"id"`
	Name         string     `json:"name" binding:"required"` 
	CreatedAt    *time.Time `json:"created_at,omitempty"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`
	ProductCount *int64     `json:"product_count,omitempty"`
}
