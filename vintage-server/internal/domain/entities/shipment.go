package entities

import (
    "errors"
    "time"
)

type Shipment struct {
    ID              uint64    `json:"id" db:"id"`
    OrderID         uint64    `json:"order_id" db:"order_id"`
    AddressID       uint64    `json:"address_id" db:"address_id"`
    Courier         string    `json:"courier" db:"courier"`
    Service         string    `json:"service" db:"service"`
    ShippingCost    uint64    `json:"shipping_cost" db:"shipping_cost"`
    TrackingNumber  *string   `json:"tracking_number" db:"tracking_number"`
    Status          uint8     `json:"status" db:"status"`
    CreatedAt       time.Time `json:"created_at" db:"created_at"`
    UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

func (s *Shipment) Validate() error {
    if s.OrderID == 0 {
        return errors.New("order_id is required")
    }
    if s.AddressID == 0 {
        return errors.New("address_id is required")
    }
    if s.Courier == "" {
        return errors.New("courier is required")
    }
    if s.Service == "" {
        return errors.New("service is required")
    }
    return nil
}

func (s *Shipment) HasTracking() bool {
    return s.TrackingNumber != nil && *s.TrackingNumber != ""
}

func (s *Shipment) UpdateTracking(trackingNumber string) {
    s.TrackingNumber = &trackingNumber
}
