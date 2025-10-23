package domain

import (
	"fmt"
	"time"
	"vintage-server/pkg/hash"

	"github.com/google/uuid"
)

type AccountRole int

const (
	RoleCustomer AccountRole = iota + 1
	RoleSeller
	RoleAdmin
)

type Account struct {
	ID        uuid.UUID
	Username  string
	Firstname string
	Lastname  *string
	Password  string
	Email     string
	AvatarURL *string
	Active    bool
	Roles     []AccountRole
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Logika bisnis contoh: aktifkan atau nonaktifkan akun
func (a *Account) Activate() {
	a.Active = true
	a.UpdatedAt = time.Now()
}

func (a *Account) Deactivate() {
	a.Active = false
	a.UpdatedAt = time.Now()
}

func (a *Account) UpdateEmail(email string) error {
	if email == "" || !containsAt(email) {
		return fmt.Errorf("email tidak valid")
	}
	a.Email = email
	a.UpdatedAt = time.Now()
	return nil
}

func containsAt(email string) bool {
	for _, c := range email {
		if c == '@' {
			return true
		}
	}
	return false
}

func (a *Account) ComparePassword(password string) bool {
	return hash.Verify(a.Password, password) == nil
}
