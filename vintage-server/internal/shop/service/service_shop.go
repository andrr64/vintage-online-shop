package service

import (
	"vintage-server/internal/shop/repository"
	"vintage-server/pkg/uploader"
)

type ShopService interface {
	
}

type shopSvc struct {
	store repository.ShopStore
	uploader uploader.Uploader
}


func NewShopSvc(
	store repository.ShopStore,
	uploader uploader.Uploader,
) ShopService{
	return &shopSvc {
		store: store,
		uploader: uploader,
	}
}