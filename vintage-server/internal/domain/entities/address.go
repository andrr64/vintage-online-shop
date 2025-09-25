package entities

import (
    "errors"
    "time"
)

type Address struct {
    ID             uint64    `json:"id" db:"id"`
    AccountID      uint64    `json:"account_id" db:"account_id"`
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

func (a *Address) Validate() error {
    if a.AccountID == 0 {
        return errors.New("account_id is required")
    }
    if a.RecipientName == "" {
        return errors.New("recipient name is required")
    }
    if a.RecipientPhone == "" {
        return errors.New("recipient phone is required")
    }
    if a.Street == "" {
        return errors.New("street address is required")
    }
    return nil
}

func (a *Address) SetAsPrimary() {
    a.IsPrimary = true
}

func (a *Address) SetAsNonPrimary() {
    a.IsPrimary = false
}

func (a *Address) GetFullAddress() string {
    return a.Street + ", " + a.PostalCode
}
