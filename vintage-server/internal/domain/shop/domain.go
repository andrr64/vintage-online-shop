package shop

import (
	"context"
	"vintage-server/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ShopRepository interface {
	WithTx(tx *sqlx.Tx) ShopRepository
	CreateShop(ctx context.Context, req model.Shop) (model.Shop, error)
	AddSellerRole(ctx context.Context, accountID uuid.UUID, sellerRoleID int) error

	GetRoleIDByName(ctx context.Context, roleName string) (int64, error)
	InsertAccountRole(ctx context.Context, accountID uuid.UUID, roleID int64) error
}

type ShopService interface {
	CreateShop(ctx context.Context, accountID uuid.UUID, req ShopRequest) (ShopDetail, error)
	UpdateShop(ctx context.Context, accountID uuid.UUID, shopID uuid.UUID, req ShopRequest) (ShopDetail, error)
}

type ShopHandler interface {
	CreateShop(c *gin.Context)
	UpdateShop(c *gin.Context)
}
