package shop

import (
	"context"
	"vintage-server/internal/domain/shop"
	"vintage-server/internal/repository"
	"vintage-server/pkg/auth"
	"vintage-server/pkg/uploader"
)

type shopService struct {
	store    repository.ShopStore
	jwt      auth.JWTService
	uploader uploader.Uploader
}

// CreateShop implements shop.ShopService.
func (s *shopService) CreateShop(ctx context.Context, req shop.CreateShop) {
	panic("unimplemented")
}

func NewShopService(
	store repository.ShopStore,
	jwt auth.JWTService,
	uploader uploader.Uploader,
) shop.ShopService {
	return &shopService{
		store:    store,
		jwt:      jwt,
		uploader: uploader,
	}
}
