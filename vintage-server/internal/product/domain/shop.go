package domain

import (
	"time"

	"github.com/google/uuid"
)

type Shop struct {
	ID          uuid.UUID
	AccountID   uuid.UUID
	Name        string
	Summary     *string
	Description *string
	Active      bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
