package repository

import (
	"context"
	"vintage-server/internal/domain/shop"
	"vintage-server/internal/model"

	"github.com/jmoiron/sqlx"
)

type shopRepository struct {
	db *sqlx.DB
	tx *sqlx.Tx
}

// CreateShop implements shop.ShopRepository.
func (s *shopRepository) CreateShop(ctx context.Context, req shop.CreateShop) (model.Shop, error) {
	panic("unimplemented")
}

// WithTx implements shop.ShopRepository.
func (s *shopRepository) WithTx(tx *sqlx.Tx) shop.ShopRepository {
	return &shopRepository{db: s.db, tx: tx}
}

func NewShopRepository(db *sqlx.DB) shop.ShopRepository {
	return &shopRepository{
		db: db,
	}
}
