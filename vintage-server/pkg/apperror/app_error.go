// File: pkg/errors/app_error.go
package apperror

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/lib/pq"
)

const (
	// ErrCodeValidation - 400 Bad Request
	// Digunakan untuk error validasi input dari user.
	ErrCodeValidation = http.StatusBadRequest

	ErrCodeBadRequest = http.StatusBadRequest

	// ErrCodeUnauthorized - 401 Unauthorized
	// Digunakan untuk kegagalan autentikasi (login gagal, token tidak valid).
	ErrCodeUnauthorized = http.StatusUnauthorized

	// ErrCodeForbidden - 403 Forbidden
	// Digunakan saat user terautentikasi tapi tidak punya hak akses ke resource.
	ErrCodeForbidden = http.StatusForbidden

	// ErrCodeNotFound - 404 Not Found
	// Digunakan saat resource yang dicari tidak ditemukan.
	ErrCodeNotFound = http.StatusNotFound

	// ErrCodeConflict - 409 Conflict
	// Digunakan saat ada konflik data, misal: mencoba mendaftar dengan email yang sudah ada.
	ErrCodeConflict = http.StatusConflict

	// ErrCodeInternal - 500 Internal Server Error
	// Digunakan untuk error-error tak terduga di sisi server.
	ErrCodeInternal = http.StatusInternalServerError
)

type AppError struct {
	Code    int
	Message string
	// Kamu bisa tambah field lain di sini, misal: TraceID, etc.
}

// INI ADALAH KUNCINYA
// Karena ada method ini, maka *AppError sekarang adalah sebuah 'error'
func (e *AppError) Error() string {
	return e.Message
}

func New(code int, message string) error {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

func HandleDBError(err error, logContext string) error {
	if err == nil {
		return nil
	}

	// SQL no rows
	if errors.Is(err, sql.ErrNoRows) {
		return New(ErrCodeNotFound, "resource not found")
	}

	// bisa juga handle constraint violation
	if pqErr, ok := err.(*pq.Error); ok {
		switch pqErr.Code {
		case "23505": // unique violation
			return New(ErrCodeConflict, "duplicate entry")
		case "23503": // foreign key violation
			return New(ErrCodeBadRequest, "foreign key constraint violation")
		}
	}

	log.Printf("%s: %v", logContext, err)
	return New(ErrCodeInternal, "an internal error occurred")
}
