package helper

import (
	"errors"
	"fmt"
	"net/http"
	"vintage-server/pkg/apperror"
	"vintage-server/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// checkAuthAndRole adalah helper untuk mengekstrak token, memvalidasi, dan memeriksa role.
// Fungsi ini mengembalikan accountID jika otentikasi dan otorisasi berhasil.
// Jika allowedRoles tidak disediakan, fungsi ini hanya akan memeriksa otentikasi (token valid).
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


// handleError adalah helper internal untuk menangani error dari service secara konsisten
func HandleError(c *gin.Context, err error) {
	var appErr *apperror.AppError
	if errors.As(err, &appErr) {
		response.Error(c, appErr.Code, appErr.Message)
	} else {
		// Sembunyikan detail error internal dari client
		response.Error(c, http.StatusInternalServerError, "an unexpected internal error occurred")
	}
}