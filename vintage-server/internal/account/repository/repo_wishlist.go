package repository

import (
	"context"
	"vintage-server/internal/account/domain"
	"vintage-server/internal/shared/db"

	"github.com/google/uuid"
)

type WishlistRepository interface {
	WithTx(tx db.DBTX) WishlistRepository

	SaveWishlistItem(ctx context.Context, accountID uuid.UUID, productID uuid.UUID) error
	FindWishlistProductIDsByAccountID(ctx context.Context, accountID uuid.UUID, keyword string, limit, offset int) ([]uuid.UUID, int, error)
	FindWishlistByAccountIDAll(ctx context.Context, accountID uuid.UUID) (domain.Wishlist, int, error)
	DeleteWishlistItem(ctx context.Context, accountID, productID uuid.UUID) error
	CheckWishlistItemExists(ctx context.Context, accountID, productID uuid.UUID) (bool, error)
}

// implementasi konkret
type wishlistRepoPostgres struct {
	db db.DBTX
}

// constructor
func NewWishlistRepositoryPostgres(db db.DBTX) WishlistRepository {
	return &wishlistRepoPostgres{db: db}
}
