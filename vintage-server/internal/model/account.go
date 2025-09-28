package model

import (
	"time"

	"github.com/google/uuid"
)

const (
	RoleCustomer = iota + 1
	RoleSeller
	RoleAdmin
)

// Account merepresentasikan tabel 'accounts'
type Account struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	Firstname string    `json:"firstname" db:"firstname"`
	Lastname  *string   `json:"lastname" db:"lastname"`
	Password  string    `json:"-" db:"password"`
	Email     string    `json:"email" db:"email"`
	AvatarURL *string   `json:"avatar_url" db:"avatar_url"`
	Active    bool      `json:"active" db:"active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Roles struct {
	ID   int64  `json:"role_id" db:"id"`
	Name string `json:"role" db:"name"`
}

// Address merepresentasikan tabel 'addresses'
type Address struct {
	ID             int64     `json:"id" db:"id"`
	AccountID      uuid.UUID `json:"account_id" db:"account_id"`
	DistrictID     string    `json:"district_id" db:"district_id"`
	RegencyID      string    `json:"regency_id" db:"regency_id"`
	ProvinceID     string    `json:"province_id" db:"province_id"`
	Label          string    `json:"label" db:"label"`
	RecipientName  string    `json:"recipient_name" db:"recipient_name"`
	RecipientPhone string    `json:"recipient_phone" db:"recipient_phone"`
	Street         string    `json:"street" db:"street"`
	PostalCode     string    `json:"postal_code" db:"postal_code"`
	IsPrimary      bool      `json:"is_primary" db:"is_primary"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// Wishlist merepresentasikan tabel 'wishlist'
type Wishlist struct {
	ID        int64     `json:"id" db:"id"`
	AccountID int64     `json:"account_id" db:"account_id"`
	ProductID int64     `json:"product_id" db:"product_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// AdminLog merepresentasikan tabel 'admin_logs'
type AdminLog struct {
	ID          int64     `json:"id" db:"id"`
	AdminID     int64     `json:"admin_id" db:"admin_id"`
	Action      string    `json:"action" db:"action"`
	Description *string   `json:"description" db:"description"`
	IPAddress   string    `json:"ip_address" db:"ip_address"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}
