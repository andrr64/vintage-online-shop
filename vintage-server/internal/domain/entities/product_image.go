package entities

import (
    "errors"
    "time"
)

type ProductImage struct {
    ID        uint64    `json:"id" db:"id"`
    ProductID uint64    `json:"product_id" db:"product_id"`
    Index     uint8     `json:"index" db:"index"`
    URL       string    `json:"url" db:"url"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
}

func (pi *ProductImage) Validate() error {
    if pi.ProductID == 0 {
        return errors.New("product_id is required")
    }
    if pi.URL == "" {
        return errors.New("image URL is required")
    }
    return nil
}
