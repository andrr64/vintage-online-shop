package entities

import (
    "errors"
    "time"
)

type Order struct {
    ID         uint64    `json:"id" db:"id"`
    AccountID  uint64    `json:"account_id" db:"account_id"`
    TotalPrice uint64    `json:"total_price" db:"total_price"`
    Status     uint8     `json:"status" db:"status"`
    CreatedAt  time.Time `json:"created_at" db:"created_at"`
    UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

func (o *Order) Validate() error {
    if o.AccountID == 0 {
        return errors.New("account_id is required")
    }
    if o.TotalPrice == 0 {
        return errors.New("total_price must be greater than 0")
    }
    return nil
}

func (o *Order) CanBeCancelled() bool {
    // Define cancellable statuses
    return o.Status == 1 || o.Status == 2 // Pending or Processing
}

func (o *Order) CalculateTotal(items []OrderItem) {
    total := uint64(0)
    for _, item := range items {
        total += item.PriceAtPurchase * uint64(item.Quantity)
    }
    o.TotalPrice = total
}
