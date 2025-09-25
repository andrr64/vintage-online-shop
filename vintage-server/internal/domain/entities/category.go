package entities

import (
    "errors"
    "time"
)

type Category struct {
    ID        uint32    `json:"id" db:"id"`
    Name      string    `json:"name" db:"name"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func (c *Category) Validate() error {
    if c.Name == "" {
        return errors.New("category name is required")
    }
    return nil
}
