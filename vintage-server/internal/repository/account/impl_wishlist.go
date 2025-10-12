package repository

import (
	"context"
	"vintage-server/internal/model"

	"github.com/google/uuid"
)

// CheckWishlistItemExists implements account.AccountRepository.
// CheckWishlistItemExists implements account.AccountRepository.
func (r *sqlAccountRepository) CheckWishlistItemExists(ctx context.Context, accountID int64, productID int64) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, DefaultQueryTimeout)
	defer cancel()

	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM wishlist WHERE account_id = $1 AND product_id = $2)`

	if err := r.db.GetContext(ctx, &exists, query, accountID, productID); err != nil {
		return false, err
	}

	return exists, nil
}

// DeleteAddress implements account.AccountRepository.
func (r *sqlAccountRepository) DeleteAddress(ctx context.Context, addressID int64, accountID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, DefaultQueryTimeout)
	defer cancel()

	query := `DELETE FROM addresses WHERE id = $1 AND account_id = $2`
	_, err := r.db.ExecContext(ctx, query, addressID, accountID)
	return err
}

// DeleteWishlistItem implements account.AccountRepository.
func (r *sqlAccountRepository) DeleteWishlistItem(ctx context.Context, accountID int64, productID int64) error {
	ctx, cancel := context.WithTimeout(ctx, DefaultQueryTimeout)
	defer cancel()

	query := `DELETE FROM wishlist WHERE account_id = $1 AND product_id = $2`
	_, err := r.db.ExecContext(ctx, query, accountID, productID)
	return err
}

// SaveWishlistItem implements account.AccountRepository.
func (r *sqlAccountRepository) SaveWishlistItem(ctx context.Context, item model.Wishlist) error {
	ctx, cancel := context.WithTimeout(ctx, DefaultQueryTimeout)
	defer cancel()
	query := `
		INSERT INTO wishlist (account_id, product_id)
		VALUES (:account_id, :product_id)
	`
	_, err := r.db.NamedExecContext(ctx, query, item)
	return err
}

// FindWishlistByAccountID implements account.AccountRepository.
// Mengembalikan slice model.Wishlist + total count (untuk pagination)
func (r *sqlAccountRepository) FindWishlistByAccountID(
	ctx context.Context,
	accountID uuid.UUID,
	keyword string,
	limit, offset int,
) ([]model.Wishlist, int, error) {
	ctx, cancel := context.WithTimeout(ctx, DefaultQueryTimeout)
	defer cancel()

	// 1) Ambil rows wishlist (filter optional by product name via JOIN)
	query := `
		SELECT 
			w.id,
			w.account_id,
			w.product_id,
			w.created_at,
			w.updated_at
		FROM wishlist w
		JOIN products p ON w.product_id = p.id
		WHERE w.account_id = $1
		  AND ($2 = '' OR p.name ILIKE '%' || $2 || '%')
		ORDER BY w.created_at DESC
		LIMIT $3 OFFSET $4
	`

	var wishlists []model.Wishlist
	if err := r.db.SelectContext(ctx, &wishlists, query, accountID, keyword, limit, offset); err != nil {
		return nil, 0, err
	}

	// 2) Hitung total item sesuai filter (tanpa limit/offset)
	countQuery := `
		SELECT COUNT(*)
		FROM wishlist w
		JOIN products p ON w.product_id = p.id
		WHERE w.account_id = $1
		  AND ($2 = '' OR p.name ILIKE '%' || $2 || '%')
	`

	var total int
	if err := r.db.GetContext(ctx, &total, countQuery, accountID, keyword); err != nil {
		return nil, 0, err
	}

	return wishlists, total, nil
}
