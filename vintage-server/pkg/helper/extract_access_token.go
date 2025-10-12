package helper

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var (
	ErrAccountIDNotFound = errors.New("accountID not found in context")
	ErrRolesNotFound     = errors.New("roles not found in context")
	ErrInvalidAccountID  = errors.New("invalid accountID format")
	ErrInvalidRoles      = errors.New("invalid roles format")
)

// ExtractAccountInfoFromToken ngambil accountID dan roles dari gin.Context
func ExtractAccountInfoFromToken(c *gin.Context) (uuid.UUID, *string, error) {
	// Ambil accountID
	accountIDv, exists := c.Get("accountID")
	if !exists {
		return uuid.Nil, nil, ErrAccountIDNotFound
	}
	accountID, ok := accountIDv.(uuid.UUID)
	if !ok {
		return uuid.Nil, nil, ErrInvalidAccountID
	}

	// Ambil roles
	rolesV, exists := c.Get("role")
	if !exists {
		return uuid.Nil, nil, ErrRolesNotFound
	}
	role, ok := rolesV.(string)
	if !ok {
		return uuid.Nil, nil, ErrInvalidRoles
	}

	return accountID, &role, nil
}

func ExtractAccountID(c *gin.Context) (uuid.UUID, error) {
	accountIDv, exists := c.Get("accountID")
	if !exists {
		return uuid.Nil, ErrAccountIDNotFound
	}

	accountID, ok := accountIDv.(uuid.UUID)
	if !ok {
		return uuid.Nil, ErrInvalidAccountID
	}

	return accountID, nil
}