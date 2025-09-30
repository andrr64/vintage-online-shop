// File: pkg/errors/app_error.go
package apperror

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
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
	// Jika tidak ada error, kembalikan nil.
	if err == nil {
		return nil
	}

	// Jika error adalah sql.ErrNoRows, kita kembalikan error "NotFound".
	// Nanti di service layer, ini bisa diubah lagi jika perlu (misal: menjadi Unauthorized).
	if errors.Is(err, sql.ErrNoRows) {
		return New(ErrCodeNotFound, "resource not found")
	}

	// Untuk semua error database lainnya, kita catat error teknisnya...
	log.Printf("%s: %v", logContext, err)
	// ...dan kembalikan error internal yang generik ke user.
	return New(ErrCodeInternal, "an internal error occurred")
}