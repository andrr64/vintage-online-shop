package helper

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
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

func HandleErrorBadRequest(c *gin.Context) {
	response.Success(c, http.StatusBadRequest, "Invalid request.")
}

// CheckBodyJSON melakukan binding JSON dari body request ke struct tujuan (t).
// Jika gagal, akan mengembalikan error 400 (Bad Request).
func CheckBodyJSON[T any](c *gin.Context, t *T) bool {
	if err := c.ShouldBindJSON(t); err != nil {
		response.ErrorBadRequest(c)
		return false
	}
	return true
}

func CheckBody[T any](c *gin.Context, t *T) bool {
	contentType := c.GetHeader("Content-Type")

	var err error
	if strings.HasPrefix(contentType, "application/json") {
		err = c.ShouldBindJSON(t)
	} else {
		// Untuk form dan multipart form
		err = c.ShouldBind(t)
	}

	if err != nil {
		response.ErrorBadRequest(c)
		return false
	}

	return true
}
