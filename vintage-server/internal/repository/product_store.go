package repository

import (
	"context"
	"fmt"
	"vintage-server/internal/domain/product"

	"github.com/jmoiron/sqlx"
)

type ProductStore interface {
	product.ProductRepository // pakai interface dari domain
	ExecTx(ctx context.Context, fn func(product.ProductRepository) error) error
}

// ProductSqlStore menyediakan implementasi dari Store untuk database SQL.
type ProductSqlStore struct {
	db *sqlx.DB
	product.ProductRepository
}

// NewProductStore adalah "pintu masuk" utama dari aplikasi kita ke lapisan data.
// Ia membuat sebuah Store baru yang siap digunakan.
func NewProductStore(db *sqlx.DB) ProductStore {
	return &ProductSqlStore{
		db:                db,
		ProductRepository: NewProductRepository(db), // Menggunakan konstruktor Repository yang sudah kita buat
	}
}

// ExecTx adalah metode andalan kita.
// Ia mengeksekusi sebuah fungsi di dalam sebuah transaksi database.
func (s *ProductSqlStore) ExecTx(ctx context.Context, fn func(product.ProductRepository) error) error {
	// 1. Memulai transaksi ("safety bubble")
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	repoTx := s.ProductRepository.WithTx(tx)
	err = fn(repoTx) // repoTx sekarang type ProductRepository
	if err != nil {
		// 4a. Jika resep gagal, batalkan semua perubahan (Rollback)
		if rbErr := tx.Rollback(); rbErr != nil {
			// Jika rollback juga gagal, laporkan kedua error
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err // Kembalikan error asli dari resep
	}

	// 4b. Jika resep berhasil, simpan semua perubahan (Commit)
	return tx.Commit()
}
