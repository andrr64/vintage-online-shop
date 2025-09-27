package model

import "time"

// Order merepresentasikan tabel 'orders'
type Order struct {
	ID         int64     `json:"id" db:"id"`
	AccountID  int64     `json:"account_id" db:"account_id"`
	TotalPrice int64     `json:"total_price" db:"total_price"`
	Status     int16     `json:"status" db:"status"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// OrderItem merepresentasikan tabel 'order_items'
type OrderItem struct {
	ID              int64     `json:"id" db:"id"`
	OrderID         int64     `json:"order_id" db:"order_id"`
	ProductID       int64     `json:"product_id" db:"product_id"`
	Quantity        int       `json:"quantity" db:"quantity"`
	PriceAtPurchase int64     `json:"price_at_purchase" db:"price_at_purchase"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}

// OrderStatusLog merepresentasikan tabel 'order_status_logs'
type OrderStatusLog struct {
	ID        int64     `json:"id" db:"id"`
	OrderID   int64     `json:"order_id" db:"order_id"`
	OldStatus *int16    `json:"old_status" db:"old_status"`
	NewStatus int16     `json:"new_status" db:"new_status"`
	Note      *string   `json:"note" db:"note"`
	CreatedBy *int64    `json:"created_by" db:"created_by"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Payment merepresentasikan tabel 'payments'
type Payment struct {
	ID                    int64     `json:"id" db:"id"`
	OrderID               int64     `json:"order_id" db:"order_id"`
	PaymentStatus         string    `json:"payment_status" db:"payment_status"`
	MidtransOrderID       string    `json:"midtrans_order_id" db:"midtrans_order_id"`
	MidtransTransactionID *string   `json:"midtrans_transaction_id" db:"midtrans_transaction_id"`
	PaymentMethod         *string   `json:"payment_method" db:"payment_method"`
	CreatedAt             time.Time `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time `json:"updated_at" db:"updated_at"`
}

// Shipment merepresentasikan tabel 'shipments'
type Shipment struct {
	ID             int64     `json:"id" db:"id"`
	OrderID        int64     `json:"order_id" db:"order_id"`
	AddressID      int64     `json:"address_id" db:"address_id"`
	Courier        string    `json:"courier" db:"courier"`
	Service        string    `json:"service" db:"service"`
	ShippingCost   int64     `json:"shipping_cost" db:"shipping_cost"`
	TrackingNumber *string   `json:"tracking_number" db:"tracking_number"`
	Status         int16     `json:"status" db:"status"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// Cart merepresentasikan tabel 'cart'
type Cart struct {
	ID        int64     `json:"id" db:"id"`
	AccountID int64     `json:"account_id" db:"account_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// CartItem merepresentasikan tabel 'cart_items'
type CartItem struct {
	ID        int64     `json:"id" db:"id"`
	CartID    int64     `json:"cart_id" db:"cart_id"`
	ProductID int64     `json:"product_id" db:"product_id"`
	Quantity  int       `json:"quantity" db:"quantity"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
