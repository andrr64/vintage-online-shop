package db_error

import (
	"net/http"

	"github.com/lib/pq"
	"vintage-server/pkg/apperror"
)

// constraintMessages berisi daftar constraint PostgreSQL dan pesan custom-nya.
// Tambahkan sesuai kebutuhan projekmu.
var constraintMessages = map[string]string{
	"shop_account_id_key": "Anda hanya boleh memiliki satu toko.",
	"accounts_email_key":  "Email sudah terdaftar.",
	"shops_name_key":      "Nama toko sudah digunakan.",
}

// HandlePgError memeriksa apakah error berasal dari PostgreSQL,
// lalu mengubahnya jadi apperror dengan pesan yang lebih jelas.
func HandlePgError(err error) error {
	if err == nil {
		return nil
	}

	pqErr, ok := err.(*pq.Error)
	if !ok {
		// Kalau bukan error dari PostgreSQL, kembalikan seperti biasa.
		return err
	}

	// Cek apakah constraint ada di mapping
	if msg, found := constraintMessages[pqErr.Constraint]; found {
		return apperror.New(http.StatusConflict, msg)
	}

	// Kalau tidak ditemukan di mapping, balikin pesan default-nya
	return apperror.New(http.StatusInternalServerError, pqErr.Message)
}
