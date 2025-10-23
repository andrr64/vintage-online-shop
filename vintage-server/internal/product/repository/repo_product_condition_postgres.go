package repository

import (
	"context"
	"fmt"
	"time"
	"vintage-server/internal/product/domain"

	"github.com/jmoiron/sqlx"
)

// CountProductsByCondition implements ProductConditionRepository.
func (r *productConditionPostgres) CountProductsByCondition(ctx context.Context, conditionID int16) (int64, error) {
	var count int64

	query := `
		SELECT COUNT(*)
		FROM products
		WHERE condition_id = $1
	`

	err := r.db.QueryRowxContext(ctx, query, conditionID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("db_error: gagal menghitung produk dengan condition_id %d: %w", conditionID, err)
	}

	return count, nil
}

// CreateCondition implements ProductConditionRepository.
func (r *productConditionPostgres) CreateCondition(ctx context.Context, data domain.ProductCondition) (domain.ProductCondition, error) {
	var saved domain.ProductCondition

	query := `
		INSERT INTO product_conditions (name)
		VALUES ($1)
		RETURNING id, name, created_at, updated_at
	`

	err := r.db.QueryRowxContext(ctx, query, data.Name).Scan(
		&saved.ID,
		&saved.Name,
		&saved.CreatedAt,
		&saved.UpdatedAt,
	)
	if err != nil {
		return domain.ProductCondition{}, fmt.Errorf("db_error: gagal membuat product condition: %w", err)
	}

	return saved, nil
}

// DeleteCondition implements ProductConditionRepository.
func (r *productConditionPostgres) DeleteCondition(ctx context.Context, id int16) error {
	var count int64

	countQuery := `
		SELECT COUNT(id)
		FROM products
		WHERE condition_id = $1
	`
	err := r.db.QueryRowxContext(ctx, countQuery, id).Scan(&count)
	if err != nil {
		return fmt.Errorf("db_error: gagal memeriksa penggunaan condition %d: %w", id, err)
	}

	if count > 0 {
		return fmt.Errorf("db_error: tidak bisa menghapus condition %d karena masih digunakan oleh %d produk", id, count)
	}

	deleteQuery := `
		DELETE FROM product_conditions 
		WHERE id = $1
	`
	_, err = r.db.ExecContext(ctx, deleteQuery, id)
	if err != nil {
		return fmt.Errorf("db_error: gagal menghapus product condition %d: %w", id, err)
	}
	return nil
}

// FindProductConditions implements ProductConditionRepository.
func (r *productConditionPostgres) FindProductConditions(ctx context.Context, id *int16) ([]domain.ProductCondition, error) {
	conditions := make([]domain.ProductCondition, 0)

	var (
		rows *sqlx.Rows
		err  error
	)

	if id != nil {
		query := `
			SELECT id, name, created_at, updated_at
			FROM product_conditions
			WHERE id = $1
			ORDER BY id ASC
		`
		rows, err = r.db.QueryxContext(ctx, query, *id)
	} else {
		query := `
			SELECT id, name, created_at, updated_at
			FROM product_conditions
			ORDER BY name ASC
		`
		rows, err = r.db.QueryxContext(ctx, query)
	}

	if err != nil {
		return nil, fmt.Errorf("db_error: gagal mengambil data product conditions: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var d domain.ProductCondition
		if err := rows.Scan(
			&d.ID,
			&d.Name,
			&d.CreatedAt,
			&d.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("db_error: gagal memindai data product condition: %w", err)
		}
		conditions = append(conditions, d)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("db_error: error setelah membaca baris data: %w", err)
	}

	return conditions, nil
}

// UpdateCondition implements ProductConditionRepository.
func (r *productConditionPostgres) UpdateCondition(ctx context.Context, data domain.ProductCondition) error {
	query := `
		UPDATE product_conditions 
		SET name = $1, updated_at = $2
		WHERE id = $3
	`

	_, err := r.db.ExecContext(ctx, query, data.Name, time.Now(), data.ID)
	if err != nil {
		return fmt.Errorf("db_error: gagal memperbarui product condition %d: %w", data.ID, err)
	}

	return nil
}
