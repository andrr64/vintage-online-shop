package service

import "vintage-server/internal/shop/repository"

type ShopServices struct {
	Shop ShopService
}


func NewShopService(store repository.ShopStore) ShopServices {
	return ShopServices{
		Shop: NewShopService(store),
	}
}