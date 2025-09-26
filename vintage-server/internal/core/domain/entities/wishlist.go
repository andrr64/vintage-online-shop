// domain/entities/wishlist.go
package entities

import "time"

type Wishlist struct {
	ID        uint64
	AccountID uint64
	ProductID uint64
	CreatedAt time.Time
	UpdatedAt time.Time
}
