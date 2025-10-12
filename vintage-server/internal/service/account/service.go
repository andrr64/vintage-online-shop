package service

import (
	"vintage-server/internal/domain/account"
	repository "vintage-server/internal/repository/account"
	"vintage-server/pkg/auth"
	"vintage-server/pkg/uploader"
)

type accountService struct {
	store    repository.AccountStore
	jwt      *auth.JWTService
	uploader uploader.Uploader
}

func NewService(store repository.AccountStore, jwtSecret string, uploader uploader.Uploader) account.AccountService {
	return &accountService{
		store:    store,
		jwt:      auth.NewJWTService(jwtSecret),
		uploader: uploader,
	}
}
