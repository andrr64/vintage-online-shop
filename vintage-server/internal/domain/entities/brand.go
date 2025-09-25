package entities

import (
    "errors"
    "time"
)

type Brand struct {
    ID        uint32    `json:"id" db:"id"`
    Name      string    `json:"name" db:"name"`
    LogoURL   *string   `json:"logo_url" db:"logo_url"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func (b *Brand) Validate() error {
    if b.Name == "" {
        return errors.New("brand name is required")
    }
    return nil
}
