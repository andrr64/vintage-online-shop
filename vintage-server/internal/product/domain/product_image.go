package domain

import (
	"time"

	"github.com/google/uuid"
)

type ProductImage struct {
	ID         int64
	ProductID  uuid.UUID
	ImageIndex int16
	URL        string
	CreatedAt  time.Time
}
