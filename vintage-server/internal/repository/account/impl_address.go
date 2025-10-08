package repository

import (
	"context"
	"database/sql"
	"fmt"
	"vintage-server/internal/model"

	"github.com/google/uuid"
)

// UpdateAddress implements account.AccountRepository.
func (r *sqlAccountRepository) UpdateAddress(ctx context.Context, accountID uuid.UUID, address model.Address) (model.Address, error) {
	ctx, cancel := context.WithTimeout(ctx, DefaultQueryTimeout)
	defer cancel()

	query := `
		UPDATE addresses SET
			label = :label,
			province_id = :province_id,
			regency_id = :regency_id,
			district_id = :district_id,
			village_id = :village_id,
			recipient_name = :recipient_name,
			recipient_phone = :recipient_phone,
			street = :street,
			postal_code = :postal_code,
			updated_at = :updated_at
		WHERE id = :id AND account_id = :account_id
		RETURNING *
	`

	var updatedAddress model.Address
	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return model.Address{}, err
	}
	defer stmt.Close()

	if err := stmt.GetContext(ctx, &updatedAddress, address); err != nil {
		return model.Address{}, err
	}

	return updatedAddress, nil
}

func (r *sqlAccountRepository) SaveAddress(ctx context.Context, address model.Address) (model.Address, error) {
	ctx, cancel := context.WithTimeout(ctx, DefaultQueryTimeout)
	defer cancel()

	var savedAddress model.Address
	query := `
		INSERT INTO addresses (
			account_id, village_id, district_id, regency_id, province_id, label,
			recipient_name, recipient_phone, street, postal_code, is_primary, created_at, updated_at
		)
		VALUES (
			:account_id, :village_id, :district_id, :regency_id, :province_id, :label,
			:recipient_name, :recipient_phone, :street, :postal_code, :is_primary, :created_at, :updated_at
		)
		RETURNING *
	`

	// Prepare statement pakai DBTX
	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return model.Address{}, err
	}
	defer stmt.Close()

	if err := stmt.GetContext(ctx, &savedAddress, address); err != nil {
		return model.Address{}, err
	}

	return savedAddress, nil
}

// Transactional SetPrimaryAddress: semua logic digabung, pakai repo langsung
func (r *sqlAccountRepository) SetPrimaryAddress(ctx context.Context, accountID uuid.UUID, addressID int64) error {
	// 1. cek apakah address ada
	var exists bool
	queryExists := `SELECT EXISTS(SELECT 1 FROM addresses WHERE id = $1 AND account_id = $2)`
	if err := r.db.GetContext(ctx, &exists, queryExists, addressID, accountID); err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("address not found")
	}

	// 2. ambil current primary
	var currentID int64
	queryCurrent := `SELECT id FROM addresses WHERE account_id = $1 AND is_primary = true LIMIT 1`
	err := r.db.GetContext(ctx, &currentID, queryCurrent, accountID)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	// 3. unset primary lama jika ada
	if currentID != 0 && currentID != addressID {
		if _, err := r.db.ExecContext(ctx, `UPDATE addresses SET is_primary = false WHERE id = $1 AND account_id = $2`, currentID, accountID); err != nil {
			return err
		}
	}

	// 4. set primary baru
	if _, err := r.db.ExecContext(ctx, `UPDATE addresses SET is_primary = true WHERE id = $1 AND account_id = $2`, addressID, accountID); err != nil {
		return err
	}

	return nil
}

// FindAddressByIDAndAccountID implements account.AccountRepository.
func (r *sqlAccountRepository) FindAddressByIDAndAccountID(ctx context.Context, addressID int64, accountID uuid.UUID) (model.Address, error) {
	ctx, cancel := context.WithTimeout(ctx, DefaultQueryTimeout)
	defer cancel()

	var address model.Address
	query := `SELECT * FROM addresses WHERE id = $1 AND account_id = $2`

	err := r.db.GetContext(ctx, &address, query, addressID, accountID)
	if err != nil {
		return model.Address{}, err
	}

	return address, nil
}

// FindAddressesByAccountID implements account.AccountRepository.
func (r *sqlAccountRepository) FindAddressesByAccountID(ctx context.Context, accountID uuid.UUID) ([]model.Address, error) {
	ctx, cancel := context.WithTimeout(ctx, DefaultQueryTimeout)
	defer cancel()

	var addresses []model.Address
	query := `SELECT * FROM addresses WHERE account_id = $1 ORDER BY is_primary DESC, updated_at DESC`

	err := r.db.SelectContext(ctx, &addresses, query, accountID)
	if err != nil {
		return nil, err
	}

	return addresses, nil
}
