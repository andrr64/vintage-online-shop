// domain/entities/cart.go
package entities

import "time"

type Cart struct {
	ID        uint64
	AccountID uint64
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CartItem struct {
	ID        uint64
	CartID    uint64
	ProductID uint64
	Quantity  uint32
	CreatedAt time.Time
	UpdatedAt time.Time
}
