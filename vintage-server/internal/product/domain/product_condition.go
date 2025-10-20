package domain

import "time"

type ProductCondition struct {
	ID        int16
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
