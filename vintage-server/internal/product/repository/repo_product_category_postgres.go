package repository

import (
	"context"
	"fmt"
	"time"
	"vintage-server/internal/product/domain"

	"github.com/jmoiron/sqlx"
)

// CountProductsByCategory implements ProductCategoryRepository.
func (r *postgresImpl) CountProductsByCategory(ctx context.Context, categoryID int) (int64, error) {
	var count int64
	query := `
		SELECT COUNT(*)
		FROM products
		WHERE category_id = $1
	`
	err := r.db.QueryRowxContext(ctx, query, categoryID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("db_error: failed to count products with category_id %d: %w", categoryID, err)
	}
	return count, nil
}

// CreateCategory implements ProductCategoryRepository.
func (r *postgresImpl) CreateCategory(ctx context.Context, data domain.ProductCategory) (domain.ProductCategory, error) {
	var saved domain.ProductCategory
	query := `
		INSERT INTO product_categories (name)
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
		return domain.ProductCategory{}, fmt.Errorf("db_error: failed to save new product category: %w", err)
	}
	return saved, nil
}

// DeleteCategory implements ProductCategoryRepository.
func (r *postgresImpl) DeleteCategory(ctx context.Context, categoryID int) error {
	// check if any products are still using this category
	var count int64
	countQuery := `
		SELECT COUNT(*)
		FROM products
		WHERE category_id = $1
	`
	err := r.db.QueryRowxContext(ctx, countQuery, categoryID).Scan(&count)
	if err != nil {
		return fmt.Errorf("db_error: failed to check usage of category %d: %w", categoryID, err)
	}
	if count > 0 {
		return fmt.Errorf("db_error: cannot delete category %d because it is used by %d products", categoryID, count)
	}

	deleteQuery := `
		DELETE FROM product_categories
		WHERE id = $1
	`
	_, err = r.db.ExecContext(ctx, deleteQuery, categoryID)
	if err != nil {
		return fmt.Errorf("db_error: failed to delete category %d: %w", categoryID, err)
	}
	return nil
}

// FindAllCategories implements ProductCategoryRepository.
func (r *postgresImpl) FindAllCategories(ctx context.Context, id *int) ([]domain.ProductCategory, error) {
	var (
		rows *sqlx.Rows
		err  error
	)

	if id != nil {
		query := `
			SELECT id, name, created_at, updated_at
			FROM product_categories
			WHERE id = $1
			ORDER BY id ASC
		`
		rows, err = r.db.QueryxContext(ctx, query, *id)
	} else {
		query := `
			SELECT id, name, created_at, updated_at
			FROM product_categories
			ORDER BY id ASC
		`
		rows, err = r.db.QueryxContext(ctx, query)
	}

	if err != nil {
		return nil, fmt.Errorf("db_error: failed to fetch product categories: %w", err)
	}
	defer rows.Close()

	categories := make([]domain.ProductCategory, 0)
	for rows.Next() {
		var d domain.ProductCategory
		if err := rows.Scan(
			&d.ID,
			&d.Name,
			&d.CreatedAt,
			&d.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("db_error: failed to scan product category: %w", err)
		}
		categories = append(categories, d)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("db_error: error after reading rows: %w", err)
	}

	return categories, nil
}

// UpdateCategory implements ProductCategoryRepository.
func (r *postgresImpl) UpdateCategory(ctx context.Context, data domain.ProductCategory) error {
	query := `
		UPDATE product_categories
		SET name = $1, updated_at = $2
		WHERE id = $3
	`
	_, err := r.db.ExecContext(ctx, query, data.Name, time.Now(), data.ID)
	if err != nil {
		return fmt.Errorf("db_error: failed to update category %d: %w", data.ID, err)
	}
	return nil
}
