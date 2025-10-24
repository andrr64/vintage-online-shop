package repository

import (
	"context"
	"fmt"
	"time"
	"vintage-server/internal/shop/domain"

	"github.com/google/uuid"
)

// === PRIVATE  ====
func (s *shopRepo) GetRoleIDByName(ctx context.Context, roleName string) (int, error) {
	var roleID int
	err := s.db.GetContext(ctx, &roleID, `SELECT id FROM roles WHERE name = $1 LIMIT 1`, roleName)
	if err != nil {
		return 0, fmt.Errorf("repo_error: failed to get role ID for '%s': %w", roleName, err)
	}
	return roleID, nil
}

// insertAccountRole inserts a record to account_roles table.
func (s *shopRepo) InsertAccountRole(ctx context.Context, accountID uuid.UUID, roleID int) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO account_roles (account_id, role_id) VALUES ($1, $2)`,
		accountID, roleID)
	if err != nil {
		return fmt.Errorf("repo_error: failed to assign role %d to account %s: %w", roleID, accountID, err)
	}
	return nil
}

// IsShopExists implements ShopRepository.
func (s *shopRepo) IsShopExists(ctx context.Context, accountId uuid.UUID) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM shop WHERE account_id = $1
		)
	`
	var exists bool = false

	err := s.db.QueryRowxContext(ctx, query, &accountId).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// addSellerRole assigns the "seller" role to the account.
func (s *shopRepo) AddSellerRole(ctx context.Context, accountID uuid.UUID, sellerRoleID int) error {
	query := `
		INSERT INTO account_roles (account_id, role_id)
		VALUES ($1, $2)
	`
	_, err := s.db.ExecContext(ctx, query, accountID, sellerRoleID)
	if err != nil {
		return fmt.Errorf("repo_error: failed to add seller role to account %s: %w", accountID, err)
	}
	return nil
}

// === PUBLIC ====
// CreateShop implements ShopRepository.
func (s *shopRepo) CreateShop(ctx context.Context, data domain.Shop) (domain.Shop, error) {
	query := `
		INSERT INTO shop
			(account_id, name, summary, description, active, created_at, updated_at)
		VALUES
			($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, account_id, name, summary, description, active, created_at, updated_at
	`

	var created domain.Shop
	err := s.db.QueryRowxContext(
		ctx,
		query,
		data.AccountID,
		data.Name,
		data.Summary,
		data.Description,
		true,
		time.Now(),
		time.Now(),
	).Scan(
		&created.ID,
		&created.AccountID,
		&created.Name,
		&created.Summary,
		&created.Description,
		&created.Active,
		&created.CreatedAt,
		&created.UpdatedAt,
	)

	if err != nil {
		return domain.Shop{}, fmt.Errorf("repo_error: failed to insert shop: %w", err)
	}
	return created, nil
}