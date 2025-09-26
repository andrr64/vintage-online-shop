// domain/entities/order.go
package entities

import "time"

type OrderStatus uint8

const (
	OrderPending OrderStatus = iota + 1
	OrderPaid
	OrderShipped
	OrderCompleted
	OrderCancelled
)

type Order struct {
	ID         uint64
	AccountID  uint64
	TotalPrice uint64
	Status     OrderStatus
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type OrderItem struct {
	ID              uint64
	OrderID         uint64
	ProductID       uint64
	Quantity        uint32
	PriceAtPurchase uint64
	CreatedAt       time.Time
}
