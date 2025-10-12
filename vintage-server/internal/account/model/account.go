package model

import (
	"time"

	"github.com/google/uuid"
)

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
