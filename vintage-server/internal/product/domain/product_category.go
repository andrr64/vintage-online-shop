package domain

import "time"

type ProductCategory struct {
	ID           int
	Name         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	ProductCount *int64
}
