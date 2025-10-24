package handler

import (
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

func NewShopHandler(svc service.ShopServices) ShopHandler {
	return &shopHandler{
		svc: svc,
	}
}
