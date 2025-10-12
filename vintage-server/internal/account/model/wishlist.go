package model

import (
	"time"

	"github.com/google/uuid"
)

// Wishlist merepresentasikan tabel 'wishlist'
type Wishlist struct {
	ID        int64     `json:"id" db:"id"`
	AccountID uuid.UUID `json:"account_id" db:"account_id"`
	ProductID uuid.UUID `json:"product_id" db:"product_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
