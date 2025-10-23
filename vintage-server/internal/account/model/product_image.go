package model

import (
	"time"

	"github.com/google/uuid"
)

type ProductImage struct {
	ID         int64     `json:"id" db:"id"`
	ProductID  uuid.UUID `json:"product_id" db:"product_id"`
	ImageIndex int16     `json:"image_index" db:"image_index"`
	URL        string    `json:"url" db:"url"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}
