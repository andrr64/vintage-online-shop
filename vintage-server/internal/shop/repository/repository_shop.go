package repository

import (
	"context"
	"vintage-server/internal/shared/db"
	"vintage-server/internal/shop/domain"

	"github.com/google/uuid"
)

type ShopRepository interface {
	CreateShop(ctx context.Context, data domain.Shop) (domain.Shop, error)
	IsShopExists(ctx context.Context, accountId uuid.UUID) (bool, error)

	// private helper
	AddSellerRole(ctx context.Context, accountID uuid.UUID, sellerRoleID int) error
	GetRoleIDByName(ctx context.Context, roleName string) (int, error)
	InsertAccountRole(ctx context.Context, accountID uuid.UUID, roleID int) error
}

type shopRepo struct {
	db db.DBTX
}



func NewShopPostgres(db db.DBTX) ShopRepository {
	return &shopRepo{
		db: db,
	}
}
