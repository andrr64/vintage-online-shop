package domain

import (
	"time"

	"github.com/google/uuid"
)

type Review struct {
	ID        uuid.UUID
	ProductID int64
	AccountID int64
	OrderID   int64
	Rating    int16
	Comment   *string
	CreatedAt time.Time
	UpdatedAt time.Time
}
