package entities

import (
    "errors"
    "time"
)

type Payment struct {
    ID                     uint64    `json:"id" db:"id"`
    OrderID                uint64    `json:"order_id" db:"order_id"`
    PaymentStatus          string    `json:"payment_status" db:"payment_status"`
    MidtransOrderID        string    `json:"midtrans_order_id" db:"midtrans_order_id"`
    MidtransTransactionID  *string   `json:"midtrans_transaction_id" db:"midtrans_transaction_id"`
    PaymentMethod          *string   `json:"payment_method" db:"payment_method"`
    CreatedAt              time.Time `json:"created_at" db:"created_at"`
    UpdatedAt              time.Time `json:"updated_at" db:"updated_at"`
}

func (p *Payment) Validate() error {
    if p.OrderID == 0 {
        return errors.New("order_id is required")
    }
    if p.PaymentStatus == "" {
        return errors.New("payment_status is required")
    }
    if p.MidtransOrderID == "" {
        return errors.New("midtrans_order_id is required")
    }
    return nil
}

func (p *Payment) IsPaid() bool {
    return p.PaymentStatus == "paid" || p.PaymentStatus == "captured"
}

func (p *Payment) IsPending() bool {
    return p.PaymentStatus == "pending"
}

func (p *Payment) IsFailed() bool {
    return p.PaymentStatus == "failed" || p.PaymentStatus == "expired"
}
