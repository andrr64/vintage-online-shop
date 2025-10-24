package domain

import "time"

type Brand struct {
	ID        int
	Name      string
	LogoURL   *string
	CreatedAt time.Time
	UpdatedAt time.Time
}
