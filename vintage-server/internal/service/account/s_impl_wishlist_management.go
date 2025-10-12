package service

import (
	"context"
	"log"
	"vintage-server/internal/domain/account"
	"vintage-server/internal/model"
	db_error "vintage-server/pkg"

	"github.com/google/uuid"
)

// AddToWishlist implements account.AccountService.
func (s *accountService) AddToWishlist(ctx context.Context, userID uuid.UUID, productID uuid.UUID) error {
	wishlist := model.Wishlist{
		AccountID: userID,
		ProductID: productID,
	}
	log.Printf("Logging %v", wishlist)
	return db_error.HandlePgError(s.store.SaveWishlistItem(ctx, wishlist))
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
