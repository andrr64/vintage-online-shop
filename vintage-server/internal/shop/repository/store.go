package repository

import (
	"context"
	"fmt"
	"vintage-server/internal/shared/db"

	"github.com/jmoiron/sqlx"
)

type ShopStore interface {
	ExecTx(ctx context.Context, fn func(ShopStore) error) error
}

type shopSqlStore struct {
	db db.DBTX
}

func NewShopStore(db db.DBTX) ShopStore {
	return &shopSqlStore {
		db: db,
	}
}

func (s *shopSqlStore) ExecTx(ctx context.Context, fn func(ShopStore) error) error {
	sqlDB, ok := s.db.(*sqlx.DB)
	if !ok {
		return fmt.Errorf("error on ShopStore.go: this db (db.DBTX) can't convert into *sqlx.DB")
	}
	dbWithTx, err := sqlDB.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	txStore := &shopSqlStore{
		db: dbWithTx,
	}
	if err := fn(txStore); err != nil {
		_ = dbWithTx.Rollback()
		return err
	}
	return dbWithTx.Commit()
}