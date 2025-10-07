package shop

import (
	"vintage-server/internal/domain/shop"
)

type Handler struct {
	svc shop.ShopService
}

func NewHandler(svc shop.ShopService) shop.ShopHandler {
	return &Handler{svc: svc}
}