package domain

import (
	"errors"
	"strings"
)

type Email string
type Username string
type AccountPassword string

// =====================
// Email
// =====================
func NewEmail(email string) (Email, error) {
	email = strings.TrimSpace(email)
	if email == "" {
		return "", errors.New("email cannot be empty")
	}
	if !strings.Contains(email, "@") {
		return "", errors.New("email is invalid")
	}
	return Email(email), nil
}

func (e Email) String() string {
	return string(e)
}

// =====================
// Username
// =====================
func NewUsername(username string) (Username, error) {
	username = strings.TrimSpace(username)
	if len(username) < 8 {
		return "", errors.New("username must be at least 8 characters")
	}
	return Username(username), nil
}

func (u Username) String() string {
	return string(u)
}

// =====================
// AccountPassword
// =====================
func NewAccountPassword(password string) (AccountPassword, error) {
	if len(password) < 10 {
		return "", errors.New("password must be at least 10 characters")
	}
	return AccountPassword(password), nil
}

func (p AccountPassword) String() string {
	return string(p)
}
