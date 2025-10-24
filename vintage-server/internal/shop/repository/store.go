package repository

import (
	"context"
	"fmt"
	"vintage-server/internal/shared/db"

	"github.com/jmoiron/sqlx"
)

type ShopStore interface {
	ExecTx(ctx context.Context, fn func(ShopStore) error) error
	GetShopRepo() ShopRepository
}

type shopSqlStore struct {
	ShopRepo ShopRepository
	db       db.DBTX
}

// GetShopRepo implements ShopStore.
func (s *shopSqlStore) GetShopRepo() ShopRepository {
	return s.ShopRepo
}

func NewShopStore(database db.DBTX) ShopStore {
	return &shopSqlStore{
		db:       database,
		ShopRepo: NewShopPostgres(database),
	}
}

func (s *shopSqlStore) ExecTx(ctx context.Context, fn func(ShopStore) error) error {
	sqlDB, ok := s.db.(*sqlx.DB)
	if !ok {
		return fmt.Errorf("error on ShopStore.go: db (db.DBTX) is not *sqlx.DB, cannot begin transaction")
	}

	tx, err := sqlDB.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	// buat instance baru shopSqlStore tapi dengan transaksi aktif
	txStore := &shopSqlStore{
		db:       tx,
		ShopRepo: NewShopPostgres(tx), // penting: gunakan tx
	}

	// jalankan fungsi di dalam transaksi
	if err := fn(txStore); err != nil {
		_ = tx.Rollback()
		return err
	}

	// commit transaksi
	if err := tx.Commit(); err != nil {
		_ = tx.Rollback()
		return err
	}

	return nil
}
