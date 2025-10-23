package repository

import (
	"context"
	"fmt"
	"vintage-server/internal/shared/db"

	"github.com/jmoiron/sqlx"
)

type ProductStore interface {
	ExecTx(ctx context.Context, fn func(ProductStore) error) error

	GetBrandRepository() BrandRepository
	GetProductConditionRepository() ProductConditionRepository
}

type productSqlStore struct {
	db                   db.DBTX
	ProductConditionRepo ProductConditionRepository
	BrandRepo            BrandRepository
}

func NewProductStore(db db.DBTX) ProductStore {
	return &productSqlStore{
		db:                   db,
		ProductConditionRepo: NewProductConditionRepoPostgres(db),
		BrandRepo:            NewBrandRepositoryPostgres(db),
	}
}

func (p *productSqlStore) ExecTx(ctx context.Context, fn func(ProductStore) error) error {
	sqlDB, ok := p.db.(*sqlx.DB)
	if !ok {
		return fmt.Errorf("error on ProductStore.go: this db (db.DBTX) can't convert into *sqlx.DB")
	}

	dbWithTx, err := sqlDB.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	// ✅ inject repo baru pakai koneksi transaksi
	txStore := &productSqlStore{
		db:                   dbWithTx,
		BrandRepo:            NewBrandRepositoryPostgres(dbWithTx),
		ProductConditionRepo: NewProductConditionRepoPostgres(dbWithTx),
	}

	if err := fn(txStore); err != nil {
		_ = dbWithTx.Rollback()
		return err
	}
	return dbWithTx.Commit()
}

// GetBrandRepository implements ProductStore.
func (p *productSqlStore) GetBrandRepository() BrandRepository {
	return p.BrandRepo
}

// GetBrandRepository implements ProductStore.
func (p *productSqlStore) GetProductConditionRepository() ProductConditionRepository {
	return p.ProductConditionRepo
}
