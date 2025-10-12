package repository

import (
	"context"
	"fmt"
	"vintage-server/internal/domain/shop"

	"github.com/jmoiron/sqlx"
)

type ShopStore interface {
	shop.ShopRepository
	ExecTx(ctx context.Context, fn func(shop.ShopRepository) error) error
}

type ShopSqlStore struct {
	db *sqlx.DB
	shop.ShopRepository
}

func NewShopStore(db *sqlx.DB) ShopStore {
	return &ShopSqlStore{
		db:             db,
		ShopRepository: NewShopRepository(db),
	}
}

func (s *ShopSqlStore) ExecTx(ctx context.Context, fn func(shop.ShopRepository) error) error {
	// 1. Memulai transaksi ("safety bubble")
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	repoTx := s.ShopRepository.WithTx(tx)
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
