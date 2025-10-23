package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"vintage-server/internal/account/domain"
	"vintage-server/internal/shared/db"

	"github.com/google/uuid"
)

func (a *addressRepoPostgres) AddAddress(ctx context.Context, address domain.Address) error {
	// 🔹 Cek apakah user sudah punya alamat
	var count int
	countQuery := `SELECT COUNT(*) FROM addresses WHERE account_id = $1`
	if err := a.db.QueryRowContext(ctx, countQuery, address.AccountID).Scan(&count); err != nil {
		return fmt.Errorf("failed to count existing addresses: %w", err)
	}

	// 🔹 Jika belum ada, jadikan alamat ini sebagai primary
	if count == 0 {
		address.IsPrimary = true
	}

	// 🔹 Insert alamat baru
	query := `
		INSERT INTO addresses (
			account_id,
			district_id,
			regency_id,
			province_id,
			village_id,
			label,
			recipient_name,
			recipient_phone,
			street,
			postal_code,
			is_primary
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at, updated_at
	`

	err := a.db.QueryRowContext(
		ctx,
		query,
		address.AccountID,
		address.DistrictID,
		address.RegencyID,
		address.ProvinceID,
		address.VillageID,
		address.Label,
		address.RecipientName,
		address.RecipientPhone,
		address.Street,
		address.PostalCode,
		address.IsPrimary,
	).Scan(&address.ID, &address.CreatedAt, &address.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to insert address: %w", err)
	}

	return nil
}

// RemoveAddress implements AddressRepository.
func (a *addressRepoPostgres) RemoveAddress(ctx context.Context, accountID uuid.UUID, addressID int64) error {
	// 🔹 Cek apakah alamat milik akun tersebut
	var isPrimary bool
	checkQuery := `
		SELECT is_primary 
		FROM addresses 
		WHERE id = $1 AND account_id = $2
	`
	err := a.db.QueryRowContext(ctx, checkQuery, addressID, accountID).Scan(&isPrimary)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("address not found or does not belong to this account")
		}
		return fmt.Errorf("failed to check address ownership: %w", err)
	}

	// 🔹 Hapus alamat
	deleteQuery := `DELETE FROM addresses WHERE id = $1 AND account_id = $2`
	result, err := a.db.ExecContext(ctx, deleteQuery, addressID, accountID)
	if err != nil {
		return fmt.Errorf("failed to delete address: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to verify deletion: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no address deleted")
	}

	// 🔹 Jika yang dihapus adalah alamat utama, set salah satu alamat lain jadi primary
	if isPrimary {
		updatePrimaryQuery := `
			UPDATE addresses
			SET is_primary = true, updated_at = CURRENT_TIMESTAMP
			WHERE id = (
				SELECT id FROM addresses 
				WHERE account_id = $1 
				ORDER BY created_at ASC 
				LIMIT 1
			)
		`
		if _, err := a.db.ExecContext(ctx, updatePrimaryQuery, accountID); err != nil {
			return fmt.Errorf("failed to set new primary address: %w", err)
		}
	}

	return nil
}

// GetAddresses implements AddressRepository.
func (a *addressRepoPostgres) GetAddresses(ctx context.Context, accountID uuid.UUID, addressID *int64) ([]domain.Address, error) {
	query := `
		SELECT 
			id,
			account_id,
			district_id,
			regency_id,
			province_id,
			village_id,
			label,
			recipient_name,
			recipient_phone,
			street,
			postal_code,
			is_primary,
			created_at,
			updated_at
		FROM addresses
		WHERE account_id = $1
	`
	args := []interface{}{accountID}

	// tambahkan filter jika addressID tidak nil
	if addressID != nil {
		query += " AND id = $2"
		args = append(args, *addressID)
	}

	query += " ORDER BY is_primary DESC, created_at ASC"

	rows, err := a.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query addresses: %w", err)
	}
	defer rows.Close()

	var addresses []domain.Address
	for rows.Next() {
		var addr domain.Address
		if err := rows.Scan(
			&addr.ID,
			&addr.AccountID,
			&addr.DistrictID,
			&addr.RegencyID,
			&addr.ProvinceID,
			&addr.VillageID,
			&addr.Label,
			&addr.RecipientName,
			&addr.RecipientPhone,
			&addr.Street,
			&addr.PostalCode,
			&addr.IsPrimary,
			&addr.CreatedAt,
			&addr.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan address: %w", err)
		}
		addresses = append(addresses, addr)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating address rows: %w", err)
	}

	return addresses, nil
}

// SetPrimaryAddress implements AddressRepository.
func (a *addressRepoPostgres) SetPrimaryAddress(ctx context.Context, accountID uuid.UUID, addressID int64) error {
	// 🔹 Verify address exists and belongs to the account
	var exists bool
	checkQuery := `
		SELECT EXISTS (
			SELECT 1 
			FROM addresses 
			WHERE id = $1 AND account_id = $2
		)`
	if err := a.db.QueryRowContext(ctx, checkQuery, addressID, accountID).Scan(&exists); err != nil {
		return fmt.Errorf("failed to verify address: %w", err)
	}

	if !exists {
		return fmt.Errorf("address not found or does not belong to this account")
	}

	// 🔹 Unset any existing primary address for this account
	unsetQuery := `
		UPDATE addresses 
		SET is_primary = false, updated_at = CURRENT_TIMESTAMP
		WHERE account_id = $1 AND is_primary = true
	`
	if _, err := a.db.ExecContext(ctx, unsetQuery, accountID); err != nil {
		return fmt.Errorf("failed to unset existing primary address: %w", err)
	}

	// 🔹 Set the specified address as primary
	setQuery := `
		UPDATE addresses 
		SET is_primary = true, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND account_id = $2
	`
	result, err := a.db.ExecContext(ctx, setQuery, addressID, accountID)
	if err != nil {
		return fmt.Errorf("failed to set primary address: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to verify update: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no address updated")
	}

	return nil
}

// UpdateAddress implements AddressRepository.
// UpdateAddress implements AddressRepository.
func (a *addressRepoPostgres) UpdateAddress(ctx context.Context, address domain.Address) (domain.Address, error) {
	// 🔹 Verify address exists and belongs to the account
	var exists bool
	checkQuery := `
		SELECT EXISTS (
			SELECT 1 
			FROM addresses 
			WHERE id = $1 AND account_id = $2
		)`
	if err := a.db.QueryRowContext(ctx, checkQuery, address.ID, address.AccountID).Scan(&exists); err != nil {
		return domain.Address{}, fmt.Errorf("failed to verify address: %w", err)
	}

	if !exists {
		return domain.Address{}, fmt.Errorf("address not found or does not belong to this account")
	}

	// 🔹 If setting as primary, unset any existing primary address
	if address.IsPrimary {
		unsetQuery := `
			UPDATE addresses 
			SET is_primary = false, updated_at = CURRENT_TIMESTAMP
			WHERE account_id = $1 AND is_primary = true AND id != $2
		`
		if _, err := a.db.ExecContext(ctx, unsetQuery, address.AccountID, address.ID); err != nil {
			return domain.Address{}, fmt.Errorf("failed to unset existing primary address: %w", err)
		}
	}

	// 🔹 Update the address
	updateQuery := `
		UPDATE addresses
		SET 
			district_id = $1,
			regency_id = $2,
			province_id = $3,
			village_id = $4,
			label = $5,
			recipient_name = $6,
			recipient_phone = $7,
			street = $8,
			postal_code = $9,
			is_primary = $10,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $11 AND account_id = $12
		RETURNING id, account_id, district_id, regency_id, province_id, village_id, 
		          label, recipient_name, recipient_phone, street, postal_code, 
		          is_primary, created_at, updated_at
	`
	var updatedAddress domain.Address
	err := a.db.QueryRowContext(
		ctx,
		updateQuery,
		address.DistrictID,
		address.RegencyID,
		address.ProvinceID,
		address.VillageID,
		address.Label,
		address.RecipientName,
		address.RecipientPhone,
		address.Street,
		address.PostalCode,
		address.IsPrimary,
		address.ID,
		address.AccountID,
	).Scan(
		&updatedAddress.ID,
		&updatedAddress.AccountID,
		&updatedAddress.DistrictID,
		&updatedAddress.RegencyID,
		&updatedAddress.ProvinceID,
		&updatedAddress.VillageID,
		&updatedAddress.Label,
		&updatedAddress.RecipientName,
		&updatedAddress.RecipientPhone,
		&updatedAddress.Street,
		&updatedAddress.PostalCode,
		&updatedAddress.IsPrimary,
		&updatedAddress.CreatedAt,
		&updatedAddress.UpdatedAt,
	)

	if err != nil {
		return domain.Address{}, fmt.Errorf("failed to update address: %w", err)
	}

	return updatedAddress, nil
}

func (a *addressRepoPostgres) WithTx(tx db.DBTX) AddressRepository {
	return &addressRepoPostgres{db: tx}
}
