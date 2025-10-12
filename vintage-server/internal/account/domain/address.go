package domain

import (
	"time"

	"github.com/google/uuid"
)

type Address struct {
	ID             int64
	AccountID      uuid.UUID
	DistrictID     string
	RegencyID      string
	ProvinceID     string
	VillageID      string
	Label          string
	RecipientName  string
	RecipientPhone string
	Street         string
	PostalCode     string
	IsPrimary      bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// Contoh logika: jadikan alamat utama
func (a *Address) SetPrimary() {
	a.IsPrimary = true
	a.UpdatedAt = time.Now()
}
