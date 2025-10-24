package service

import (
	"vintage-server/internal/shop/repository"
	"vintage-server/pkg/uploader"
)

type ShopServices struct {
	Shop ShopService
}

func NewShopServices(store repository.ShopStore, up uploader.Uploader) ShopServices {
	return ShopServices{
		Shop: NewShopSvc(store, up),
	}
}
