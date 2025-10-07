package security

import (
	"errors"
	"fmt"
	"vintage-server/pkg/apperror"

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

// checkAuthAndRole adalah helper untuk mengekstrak token, memvalidasi, dan memeriksa role.
func CheckAuthAndRole(c *gin.Context, allowedRoles ...string) (uuid.UUID, error) {
	accountID, role, err := ExtractAccountInfoFromToken(c)
	if err != nil {
		// Jika token tidak valid atau tidak ada
		return uuid.Nil, apperror.New(apperror.ErrCodeUnauthorized, "unauthorized access")
	}

	// Jika ada role yang diizinkan, lakukan pengecekan
	if len(allowedRoles) > 0 {
		if role == nil {
			return uuid.Nil, apperror.New(apperror.ErrCodeForbidden, "access forbidden: role not found in token")
		}

		isAllowed := false
		for _, allowedRole := range allowedRoles {
			if *role == allowedRole {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			errMsg := fmt.Sprintf("access forbidden: role '%s' is not allowed", *role)
			return uuid.Nil, apperror.New(apperror.ErrCodeForbidden, errMsg)
		}
	}

	return accountID, nil
}
