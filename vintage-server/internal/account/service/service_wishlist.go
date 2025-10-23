package service

import (
	"context"
	"github.com/google/uuid"
)

type WishlistService interface {
	GetWishlist(ctx context.Context, accountID uuid.UUID, page int, size int, keyword string) ([]uuid.UUID, int, error)
	AddToWishlist(ctx context.Context, accountID uuid.UUID, productID uuid.UUID) error
	RemoveFromWishlist(ctx context.Context, accountID uuid.UUID, productID uuid.UUID) error
}
