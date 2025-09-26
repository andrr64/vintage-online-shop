// domain/entities/shop.go
package entities

import "time"

type Shop struct {
	ID          uint64
	AccountID   uint64
	Name        string
	Summary     string
	Description string
	Active      bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
