package entities

import (
    "errors"
    "time"
)

type ProductCondition struct {
    ID        uint8     `json:"id" db:"id"`
    Name      string    `json:"name" db:"name"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func (pc *ProductCondition) Validate() error {
    if pc.Name == "" {
        return errors.New("condition name is required")
    }
    return nil
}
