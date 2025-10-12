package service

import (
	"vintage-server/internal/account/repository"
	"vintage-server/pkg/uploader"
)

// AccountServices nge-embed semua service yang ada
type AccountServices struct {
	Wishlist WishlistService
	Login    LoginService
}

// constructor
func NewAccountServices(store repository.AccountStore, jwtSecret string, upSvc uploader.Uploader) *AccountServices {
	return &AccountServices{
		Wishlist: NewWishlistService(store.GetWishlistRepo()),
		Login:    NewLoginService(store.GetAccountRepo(), jwtSecret, upSvc),
	}
}
