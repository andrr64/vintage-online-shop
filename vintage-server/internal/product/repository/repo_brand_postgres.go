package repository

import (
	"context"
	"fmt"
	"vintage-server/internal/product/domain"

	"github.com/jmoiron/sqlx"
)

// CreateBrand implements BrandRepository.
func (b *brandRepoPostgres) CreateBrand(ctx context.Context, brand domain.Brand) (domain.Brand, error) {
	query := `
		INSERT INTO brands (name, logo_url, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW())
		RETURNING id, name, logo_url, created_at, updated_at
	`
	row := b.db.QueryRowxContext(ctx, query, brand.Name, brand.LogoURL)
	if err := row.Scan(&brand.ID, &brand.Name, &brand.LogoURL, &brand.CreatedAt, &brand.UpdatedAt); err != nil {
		return domain.Brand{}, err
	}
	return brand, nil
}

// ReadBrands implements BrandRepository.
// Jika brandID != nil, maka ambil brand dengan ID tersebut.
// Jika brandID == nil, maka ambil semua brand.
func (b *brandRepoPostgres) ReadBrands(ctx context.Context, brandID *int) ([]domain.Brand, error) {
	var (
		rows *sqlx.Rows
		err  error
	)

	if brandID != nil {
		query := `SELECT id, name, logo_url, created_at, updated_at FROM brands WHERE id = $1`
		rows, err = b.db.QueryxContext(ctx, query, *brandID)
	} else {
		query := `SELECT id, name, logo_url, created_at, updated_at FROM brands ORDER BY id ASC`
		rows, err = b.db.QueryxContext(ctx, query)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var brands []domain.Brand = make([]domain.Brand, 0)
	for rows.Next() {
		var brand domain.Brand
		if err := rows.Scan(
			&brand.ID,
			&brand.Name,
			&brand.LogoURL,
			&brand.CreatedAt,
			&brand.UpdatedAt,
		); err != nil {
			return nil, err
		}
		brands = append(brands, brand)
	}

	return brands, nil
}

// UpdateBrand implements BrandRepository.
func (b *brandRepoPostgres) UpdateBrand(ctx context.Context, brand domain.Brand) (domain.Brand, error) {
	query := `
		UPDATE brands
		SET name = $1, logo_url = $2, updated_at = NOW()
		WHERE id = $3
		RETURNING id, name, logo_url, created_at, updated_at
	`
	row := b.db.QueryRowxContext(ctx, query, brand.Name, brand.LogoURL, brand.ID)
	if err := row.Scan(&brand.ID, &brand.Name, &brand.LogoURL, &brand.CreatedAt, &brand.UpdatedAt); err != nil {
		return domain.Brand{}, err
	}
	return brand, nil
}

// DeleteBrand implements BrandRepository.
func (b *brandRepoPostgres) DeleteBrand(ctx context.Context, brandID int) error {
	query := `DELETE FROM brands WHERE id = $1`
	res, err := b.db.ExecContext(ctx, query, brandID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("brand dengan id %d tidak ditemukan", brandID)
	}

	return nil
}

func (b *brandRepoPostgres) CountProducts(c context.Context, brand domain.Brand) (int64, error) {
	var count int64
	query := `
		SELECT COUNT(*) 
		FROM products 
		WHERE brand_id = $1
	`
	err := b.db.QueryRowContext(c, query, brand.ID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count products by brand: %w", err)
	}
	return count, nil
}
