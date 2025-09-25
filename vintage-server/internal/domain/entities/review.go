package entities

import (
    "errors"
    "time"
)

type Review struct {
    ID        uint64    `json:"id" db:"id"`
    ProductID uint64    `json:"product_id" db:"product_id"`
    AccountID uint64    `json:"account_id" db:"account_id"`
    OrderID   uint64    `json:"order_id" db:"order_id"`
    Rating    uint8     `json:"rating" db:"rating"`
    Comment   *string   `json:"comment" db:"comment"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func (r *Review) Validate() error {
    if r.ProductID == 0 {
        return errors.New("product_id is required")
    }
    if r.AccountID == 0 {
        return errors.New("account_id is required")
    }
    if r.OrderID == 0 {
        return errors.New("order_id is required")
    }
    if r.Rating < 1 || r.Rating > 5 {
        return errors.New("rating must be between 1 and 5")
    }
    return nil
}

func (r *Review) HasComment() bool {
    return r.Comment != nil && *r.Comment != ""
}
