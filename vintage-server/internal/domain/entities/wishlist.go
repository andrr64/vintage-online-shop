package entities

import (
    "errors"
    "time"
)

type Wishlist struct {
    ID        uint64    `json:"id" db:"id"`
    AccountID uint64    `json:"account_id" db:"account_id"`
    ProductID uint64    `json:"product_id" db:"product_id"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func (w *Wishlist) Validate() error {
    if w.AccountID == 0 {
        return errors.New("account_id is required")
    }
    if w.ProductID == 0 {
        return errors.New("product_id is required")
    }
    return nil
}
