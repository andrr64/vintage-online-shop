package db_error

import (
	"database/sql"
	"net/http"

	"github.com/lib/pq"
	"vintage-server/pkg/apperror"
)

// constraintMessages berisi daftar constraint PostgreSQL dan pesan custom-nya.
var constraintMessages = map[string]string{
	"shop_account_id_key":               "Anda hanya boleh memiliki satu toko.",
	"accounts_email_key":                "Email sudah terdaftar.",
	"shops_name_key":                    "Nama toko sudah digunakan.",
	"product_condition_name_lower_idx":  "Nama kondisi sudah digunakan.",
	"product_categories_name_lower_idx": "Nama kategori sudah digunakan.",
	"products_condition_id_fkey":        "Condition tidak ditemukan.",
	"products_category_id_fkey":         "Kategori tidak ditemukan.",
	"products_brand_id_fkey":            "Brand tidak ditemukan.",
	"products_shop_id_fkey":             "Toko tidak ditemukan.",
	"products_size_id_fkey":             "Size tidak ditemukan.",
	"wishlist_account_id_fkey" :		"Akun tidak ditemukan.",
	"wishlist_product_id_fkey" :		"Produk tidak ditemukan.",
}

// HandlePgError memeriksa apakah error berasal dari PostgreSQL,
// lalu mengubahnya jadi apperror dengan pesan yang lebih jelas.
func HandlePgError(err error) error {
	if err == nil {
		return nil
	}

	// ✅ Tangani error ketika data tidak ditemukan
	if err == sql.ErrNoRows {
		return apperror.New(http.StatusNotFound, "Data tidak ditemukan.")
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
