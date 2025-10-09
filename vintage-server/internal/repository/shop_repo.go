package repository

import (
	"context"
	"vintage-server/internal/domain/shop"
	"vintage-server/internal/model"
	"vintage-server/pkg/controller"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type sqlShopRepo struct {
	db *sqlx.DB
	tx *sqlx.Tx
}

// GetRoleIDByName implements account.AccountRepository.
func (r *sqlShopRepo) GetRoleIDByName(ctx context.Context, roleName string) (int64, error) {
	ctx, cancel := controller.WithQueryTimeout(ctx)
	defer cancel()

	var roleID int64
	err := r.db.GetContext(ctx, &roleID, `SELECT id FROM roles WHERE name = $1 LIMIT 1`, roleName)
	if err != nil {
		return 0, err
	}

	return roleID, nil
}

// InsertAccountRole implements account.AccountRepository.
func (r *sqlShopRepo) InsertAccountRole(ctx context.Context, accountID uuid.UUID, roleID int64) error {
	ctx, cancel := controller.WithQueryTimeout(ctx)
	defer cancel()

	_, err := r.db.ExecContext(ctx,
		`INSERT INTO account_roles (account_id, role_id) VALUES ($1, $2)`,
		accountID, roleID)
	return err
}

// AddSellerRole implements shop.ShopRepository.
func (r *sqlShopRepo) AddSellerRole(ctx context.Context, accountID uuid.UUID, sellerRoleID int) error {
	ctx, cancel := controller.WithQueryTimeout(ctx)
	defer cancel()
	query := `
		INSERT INTO account_roles (account_id, role_id)
		VALUES ($1, $2)
	`
	_, err := r.db.ExecContext(ctx, query, accountID, sellerRoleID)
	if err != nil {
		return err
	}
	return nil
}

// CreateShop implements shop.ShopRepository.
func (r *sqlShopRepo) CreateShop(ctx context.Context, req model.Shop) (model.Shop, error) {
	ctx, cancel := controller.WithQueryTimeout(ctx)
	defer cancel()

	query := `
		INSERT INTO shop (account_id, name, summary, description, active) 
		VALUES
			(:account_id, :name, :summary, :description, :active)
		RETURNING
			*`
	var savedShop model.Shop
	nstmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return model.Shop{}, err
	}
	defer nstmt.Close()
	if err := nstmt.GetContext(ctx, &savedShop, req); err != nil {
		return model.Shop{}, err
	}
	return savedShop, nil
}

// WithTx implements shop.ShopRepository.
func (r *sqlShopRepo) WithTx(tx *sqlx.Tx) shop.ShopRepository {
	return &sqlShopRepo{db: r.db, tx: tx}
}

func NewShopRepository(db *sqlx.DB) shop.ShopRepository {
	return &sqlShopRepo{
		db: db,
	}
}
