package order

import (
	"time"
)

type Order struct {
	ID                uint `gorm:"primaryKey"`
	UserID            uint
	ShippingAddressID uint
	OrderDate         time.Time
	CurrentStatus     string
	LineItems         []OrderLineItem
	StatusHistory     []OrderStatusHistory
}

type OrderLineItem struct {
	ID                 uint `gorm:"primaryKey"`
	OrderID            uint
	ProductVariantID   uint
	Quantity           int
	PriceAtTransaction float64
}

type OrderStatusHistory struct {
	ID        uint `gorm:"primaryKey"`
	OrderID   uint
	Status    string
	Notes     string
	CreatedAt time.Time
}
