// domain/entities/review.go
package entities

import "time"

type Review struct {
	ID        uint64
	ProductID uint64
	AccountID uint64
	OrderID   uint64
	Rating    uint8
	Comment   *string
	CreatedAt time.Time
	UpdatedAt time.Time
}
