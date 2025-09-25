package entities

import (
    "errors"
    "time"
)

type OrderItem struct {
    ID               uint64    `json:"id" db:"id"`
    OrderID          uint64    `json:"order_id" db:"order_id"`
    ProductID        uint64    `json:"product_id" db:"product_id"`
    Quantity         uint32    `json:"quantity" db:"quantity"`
    PriceAtPurchase  uint64    `json:"price_at_purchase" db:"price_at_purchase"`
    CreatedAt        time.Time `json:"created_at" db:"created_at"`
}

func (oi *OrderItem) Validate() error {
    if oi.OrderID == 0 {
        return errors.New("order_id is required")
    }
    if oi.ProductID == 0 {
        return errors.New("product_id is required")
    }
    if oi.Quantity == 0 {
        return errors.New("quantity must be greater than 0")
    }
    if oi.PriceAtPurchase == 0 {
        return errors.New("price_at_purchase must be greater than 0")
    }
    return nil
}

func (oi *OrderItem) CalculateSubtotal() uint64 {
    return oi.PriceAtPurchase * uint64(oi.Quantity)
}
