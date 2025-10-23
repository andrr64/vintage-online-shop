package repository

import (
	"context"
	"vintage-server/internal/account/domain"
	"vintage-server/internal/shared/db"

	"github.com/google/uuid"
)

// CheckWishlistItemExists implements AccountStore.
func (s *accountSqlStore) CheckWishlistItemExists(ctx context.Context, accountID uuid.UUID, productID uuid.UUID) (bool, error) {
	return s.WishlistRepo.CheckWishlistItemExists(ctx, accountID, productID)
}

// DeleteWishlistItem implements AccountStore.
func (s *accountSqlStore) DeleteWishlistItem(ctx context.Context, accountID uuid.UUID, productID uuid.UUID) error {
	return s.WishlistRepo.DeleteWishlistItem(ctx, accountID, productID)
}

// FindWishlistByAccountID implements AccountStore.
func (s *accountSqlStore) FindWishlistProductIDsByAccountID(ctx context.Context, accountID uuid.UUID, keyword string, limit int, offset int) ([]uuid.UUID, int, error) {
	return s.WishlistRepo.FindWishlistProductIDsByAccountID(ctx, accountID, keyword, limit, offset)
}

// SaveWishlistItem implements AccountStore.
func (s *accountSqlStore) SaveWishlistItem(ctx context.Context, accountID uuid.UUID, productID uuid.UUID) error {
	return s.WishlistRepo.SaveWishlistItem(ctx, accountID, productID)
}

// FindWishlistByAccountIDAll implements AccountStore.
func (s *accountSqlStore) FindWishlistByAccountIDAll(ctx context.Context, accountID uuid.UUID) (domain.Wishlist, int, error) {
	return s.WishlistRepo.FindWishlistByAccountIDAll(ctx, accountID)
}

// WithTx implements AccountStore.
func (s *accountSqlStore) WithTx(tx db.DBTX) WishlistRepository {
	return s.WishlistRepo.WithTx(tx)
}
