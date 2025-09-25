package entities

import (
    "errors"
    "time"
)

type Account struct {
    ID        uint64    `json:"id" db:"id"`
    Username  string    `json:"username" db:"username"`
    Password  string    `json:"-" db:"password"`
    Email     *string   `json:"email" db:"email"`
    AvatarURL *string   `json:"avatar_url" db:"avatar_url"`
    Active    bool      `json:"active" db:"active"`
    Role      uint8     `json:"role" db:"role"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func (a *Account) Validate() error {
    if len(a.Username) < 3 {
        return errors.New("username must be at least 3 characters")
    }
    if len(a.Password) < 6 {
        return errors.New("password must be at least 6 characters")
    }
    return nil
}

func (a *Account) SetPassword(plainPassword string) error {
    // TODO: Implement password hashing
    a.Password = plainPassword // Temporary - replace with hashing
    return nil
}

func (a *Account) IsAdmin() bool {
    return a.Role == 1 // Assuming 1 is admin role
}

func (a *Account) IsSeller() bool {
    return a.Role == 2 // Assuming 2 is seller role
}

func (a *Account) IsCustomer() bool {
    return a.Role == 0 // Assuming 0 is customer role
}
