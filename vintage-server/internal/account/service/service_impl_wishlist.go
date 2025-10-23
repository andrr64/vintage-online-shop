package service

import (
	"context"
	"fmt"
	"vintage-server/internal/account/domain"

	"vintage-server/internal/account/repository"
	db_error "vintage-server/pkg"
	serviceerror "vintage-server/pkg/service_error"

	"github.com/google/uuid"
)

type wishlistServiceImpl struct {
	store    repository.AccountStore
	maxItems int
}

// constructor
func NewWishlistService(store repository.AccountStore) WishlistService {
	return &wishlistServiceImpl{store: store, maxItems: 25}
}

// GetWishlist implements WishlistService.
func (s *wishlistServiceImpl) GetWishlist(ctx context.Context, accountID uuid.UUID, page int, size int, keyword string) ([]uuid.UUID, int, error) {
	w, total, err := s.store.GetWishlistRepo().FindWishlistProductIDsByAccountID(ctx, accountID, keyword, size, (page-1)*size)
	if err != nil {
		return []uuid.UUID{}, 0, err
	}
	return w, total, nil
}

func (s *wishlistServiceImpl) AddToWishlist(ctx context.Context, accountID uuid.UUID, productID uuid.UUID) error {
	// 1️⃣ Get user's wishlist from DB
	existingWishlist, _, err := s.store.GetWishlistRepo().FindWishlistByAccountIDAll(ctx, accountID)
	if err != nil {
		// Only return error if it's a real DB issue, not because wishlist is empty
		return serviceerror.New(serviceerror.ErrInternal, "failed to fetch wishlist", err)
	}

	// 2️⃣ Create domain Wishlist
	w := domain.NewWishlist(accountID, existingWishlist.ProductIDs) // empty slice is fine

	// 3️⃣ Add product
	if !w.CanAddMore(s.maxItems) {
		return serviceerror.New(serviceerror.ErrConflict, fmt.Sprintf("wishlist is full, maximum %d items allowed", s.maxItems), nil)
	}

	if err := w.AddProduct(productID); err != nil {
		return serviceerror.New(serviceerror.ErrBadRequest, "product already exists in wishlist", nil)
	}

	// 4️⃣ Save to DB
	return db_error.HandlePgError(s.store.GetWishlistRepo().SaveWishlistItem(ctx, accountID, productID))
}

// helper convert []model.Wishlist ke []uuid.UUID
func (s *wishlistServiceImpl) RemoveFromWishlist(ctx context.Context, accountID uuid.UUID, productID uuid.UUID) error {
	return db_error.HandlePgError(s.store.GetWishlistRepo().DeleteWishlistItem(ctx, accountID, productID))
}
