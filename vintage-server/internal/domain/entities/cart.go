package entities

import (
    "errors"
    "time"
)

type Cart struct {
    ID        uint64    `json:"id" db:"id"`
    AccountID uint64    `json:"account_id" db:"account_id"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func (c *Cart) Validate() error {
    if c.AccountID == 0 {
        return errors.New("account_id is required")
    }
    return nil
}
