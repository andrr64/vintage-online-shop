// domain/entities/account.go
package entities

import "time"

type Role uint8

const (
	RoleCustomer Role = iota + 1
	RoleSeller
	RoleAdmin
)

type Account struct {
	ID        uint64
	Username  string
	Password  string
	Email     *string
	AvatarURL *string
	Active    bool
	Role      Role
	CreatedAt time.Time
	UpdatedAt time.Time
}
