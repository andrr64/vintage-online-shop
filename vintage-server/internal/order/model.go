package order

import "time"

type Order struct {
	ID                int64 `gorm:"primaryKey;autoIncrement"`
	UserID            int64 `gorm:"not null;index"`
	ShippingAddressID int64 `gorm:"not null;index"`
	OrderDate         time.Time
	CurrentStatus     string `gorm:"size:50"`
}

type OrderLineItem struct {
	ID        int64 `gorm:"primaryKey;autoIncrement"`
	OrderID   int64 `gorm:"not null;index"`
	VariantID int64 `gorm:"not null;index"`
	Quantity  int
	Price     float64
}

type OrderStatusHistory struct {
	ID        int64  `gorm:"primaryKey;autoIncrement"`
	OrderID   int64  `gorm:"not null;index"`
	Status    string `gorm:"size:50"`
	Notes     string `gorm:"type:text"`
	CreatedAt time.Time
}
