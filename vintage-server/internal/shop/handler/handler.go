package handler

import (
	"vintage-server/internal/shop/repository"
	"vintage-server/internal/shop/service"

	"github.com/gin-gonic/gin"
)

type ShopHandler interface {
	// Shop
	CreateShop(c *gin.Context)
	UpdateShop(c *gin.Context)
}

type shopHandler struct {
	svc service.ShopServices
}

func NewShopHandler(store repository.ShopStore) ShopHandler {
	return &shopHandler{
		svc: service.NewShopService(store),
	}
}
