package service

import (
	"context"
	"vintage-server/internal/domain/account"
	"vintage-server/internal/model"
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

// AddToWishlist implements account.AccountService.
func (s *accountService) AddToWishlist(ctx context.Context, userID int64, productID int64) error {
	panic("unimplemented")
}

// DeactivateUser implements account.AccountService.
func (s *accountService) DeactivateUser(ctx context.Context, userID int64, reason string) error {
	panic("unimplemented")
}

// GetUserProfile implements account.AccountService.
func (s *accountService) GetUserProfile(ctx context.Context, userID int64) (model.Account, error) {
	panic("unimplemented")
}

// GetWishlistByUserID implements account.AccountService.
func (s *accountService) GetWishlistByUserID(ctx context.Context, userID int64) ([]account.WishlistItemDetail, error) {
	panic("unimplemented")
}

// RemoveFromWishlist implements account.AccountService.
func (s *accountService) RemoveFromWishlist(ctx context.Context, userID int64, productID int64) error {
	panic("unimplemented")
}
