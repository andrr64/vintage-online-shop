package entities

import (
    "errors"
    "time"
)

type CartItem struct {
    ID        uint64    `json:"id" db:"id"`
    CartID    uint64    `json:"cart_id" db:"cart_id"`
    ProductID uint64    `json:"product_id" db:"product_id"`
    Quantity  uint32    `json:"quantity" db:"quantity"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func (ci *CartItem) Validate() error {
    if ci.CartID == 0 {
        return errors.New("cart_id is required")
    }
    if ci.ProductID == 0 {
        return errors.New("product_id is required")
    }
    if ci.Quantity == 0 {
        return errors.New("quantity must be greater than 0")
    }
    return nil
}

func (ci *CartItem) IncreaseQuantity(amount uint32) {
    ci.Quantity += amount
}

func (ci *CartItem) DecreaseQuantity(amount uint32) error {
    if amount > ci.Quantity {
        return errors.New("cannot decrease below zero")
    }
    ci.Quantity -= amount
    return nil
}
