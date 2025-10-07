package shop

import (
	"context"
	"vintage-server/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type ShopRepository interface {
	WithTx(tx *sqlx.Tx) ShopRepository
	CreateShop(ctx context.Context, req CreateShop) (model.Shop, error)
}

type ShopService interface {
	CreateShop(ctx context.Context, req CreateShop)
}

type ShopHandler interface {
	CreateShop(c *gin.Context)
	UpdateShop(c *gin.Context)
}
