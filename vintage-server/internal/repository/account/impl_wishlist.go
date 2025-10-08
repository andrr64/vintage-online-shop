package repository

import (
	"context"
	"vintage-server/internal/domain/account"
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
		INSERT INTO wishlist (account_id, product_id, created_at, updated_at)
		VALUES (:account_id, :product_id, :created_at, :updated_at)
	`

	_, err := r.db.NamedExecContext(ctx, query, item)
	return err
}

// FindWishlistByAccountID implements account.AccountRepository.
func (r *sqlAccountRepository) FindWishlistByAccountID(ctx context.Context, accountID int64) ([]account.WishlistItemDetail, error) {
	ctx, cancel := context.WithTimeout(ctx, DefaultQueryTimeout)
	defer cancel()

	var wishlistItems []account.WishlistItemDetail
	query := `
		SELECT 
			w.product_id,
			p.name AS product_name,
			p.price,
			pi.url AS product_image_url,
			w.created_at
		FROM wishlist w
		JOIN products p ON w.product_id = p.id
		LEFT JOIN product_images pi ON p.id = pi.product_id AND pi.image_index = 0
		WHERE w.account_id = $1
		ORDER BY w.created_at DESC
	`

	if err := r.db.SelectContext(ctx, &wishlistItems, query, accountID); err != nil {
		return nil, err
	}

	return wishlistItems, nil
}
