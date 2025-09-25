package entities

import (
    "errors"
    "time"
)

type Shop struct {
    ID          uint64    `json:"id" db:"id"`
    AccountID   uint64    `json:"account_id" db:"account_id"`
    Name        string    `json:"name" db:"name"`
    Summary     *string   `json:"summary" db:"summary"`
    Description *string   `json:"description" db:"description"`
    Active      bool      `json:"active" db:"active"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

func (s *Shop) Validate() error {
    if len(s.Name) < 2 {
        return errors.New("shop name must be at least 2 characters")
    }
    if s.AccountID == 0 {
        return errors.New("account_id is required")
    }
    return nil
}

func (s *Shop) Activate() {
    s.Active = true
}

func (s *Shop) Deactivate() {
    s.Active = false
}
