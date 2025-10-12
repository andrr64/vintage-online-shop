// wishlist_repo_impl_postgres.go

package repository

import (
	"context"
	"log"
	"vintage-server/internal/account/domain"
	"vintage-server/internal/account/model"
	"vintage-server/internal/shared/db"

	"github.com/google/uuid"
)

// CheckWishlistItemExists implements WishlistRepository.
func (r *wishlistRepoPostgres) CheckWishlistItemExists(ctx context.Context, accountID uuid.UUID, productID uuid.UUID) (bool, error) {

	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM wishlist WHERE account_id = $1 AND product_id = $2)`

	if err := r.db.GetContext(ctx, &exists, query, accountID, productID); err != nil {
		return false, err
	}

	return exists, nil
}

// DeleteWishlistItem implements WishlistRepository.
func (r *wishlistRepoPostgres) DeleteWishlistItem(ctx context.Context, accountID uuid.UUID, productID uuid.UUID) error {

	query := `DELETE FROM wishlist WHERE account_id = $1 AND product_id = $2`
	_, err := r.db.ExecContext(ctx, query, accountID, productID)
	return err
}

func (r *wishlistRepoPostgres) FindWishlistProductIDsByAccountID(ctx context.Context, accountID uuid.UUID, keyword string, limit, offset int) ([]uuid.UUID, int, error) {
	query := `
    SELECT w.product_id
    FROM wishlist w
    JOIN products p ON w.product_id = p.id
    WHERE w.account_id=$1 AND ($2='' OR p.name ILIKE '%' || $2 || '%')
    ORDER BY w.created_at DESC
    LIMIT $3 OFFSET $4
`
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	log.Printf("FindWishlistProductIDsByAccountID: accountID=%s, keyword=%s, limit=%d, offset=%d", accountID, keyword, limit, offset)
	var ids []uuid.UUID
	if err := r.db.SelectContext(ctx, &ids, query, accountID, keyword, limit, offset); err != nil {
		return nil, 0, err
	}
	if ids == nil {
		ids = []uuid.UUID{}
	}

	countQuery := `
        SELECT COUNT(*)
        FROM wishlist w
        JOIN products p ON w.product_id = p.id
        WHERE w.account_id=$1 AND ($2='' OR p.name ILIKE '%'||$2||'%')
    `
	var total int
	if err := r.db.GetContext(ctx, &total, countQuery, accountID, keyword); err != nil {
		return nil, 0, err
	}

	return ids, total, nil
}

func (r *wishlistRepoPostgres) FindWishlistByAccountIDAll(ctx context.Context, accountID uuid.UUID) (domain.Wishlist, int, error) {
	// 1) Ambil semua wishlist user
	query := `
		SELECT 
			w.id,
			w.account_id,
			w.product_id,
			w.created_at,
			w.updated_at
		FROM wishlist w
		WHERE w.account_id = $1
		ORDER BY w.created_at DESC
	`

	var models []model.Wishlist
	if err := r.db.SelectContext(ctx, &models, query, accountID); err != nil {
		return domain.Wishlist{}, 0, err
	}

	// konversi model ke domain entity
	productIDs := make([]uuid.UUID, 0, len(models)) // slice kosong jika tidak ada
	for _, m := range models {
		productIDs = append(productIDs, m.ProductID)
	}

	w := domain.NewWishlist(accountID, productIDs)

	return *w, len(productIDs), nil
}

// IsUsernameUsed implements WishlistRepository.
func (r *wishlistRepoPostgres) IsUsernameUsed(ctx context.Context, username string) (bool, error) {
	panic("unimplemented")
}

// SaveWishlistItem implements WishlistRepository.
func (r *wishlistRepoPostgres) SaveWishlistItem(ctx context.Context, accountID uuid.UUID, productID uuid.UUID) error {
	query := `
		INSERT INTO wishlist (account_id, product_id)
		VALUES ($1, $2)
	`
	_, err := r.db.ExecContext(ctx, query, accountID, productID)
	return err
}

func (r *wishlistRepoPostgres) WithTx(tx db.DBTX) WishlistRepository {
	return &wishlistRepoPostgres{
		db: tx,
	}
}
