package product

import (
	"context"
	"vintage-server/internal/model"
	"vintage-server/pkg/apperror"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq" // <-- Import penting
)

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) CreateCategory(ctx context.Context, data model.ProductCategory) error {
	query := `INSERT INTO product_categories (name) VALUES ($1)`

	_, err := r.db.ExecContext(ctx, query, data.Name)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return apperror.New(apperror.ErrCodeConflict, "category with this name already exists")
		}
		return apperror.HandleDBError(err, "failed to create category")
	}
	return nil
}

func (r *repository) FindAllCategories(ctx context.Context) ([]model.ProductCategory, error) {
	query := `SELECT id, name, created_at, updated_at FROM product_categories ORDER BY created_at DESC`
	var result []model.ProductCategory

	err := r.db.SelectContext(ctx, &result, query)
	if err != nil {
		return nil, apperror.HandleDBError(err, "failed to find all categories")
	}
	return result, nil
}

func (r *repository) FindById(ctx context.Context, id int) (model.ProductCategory, error) {
	var result model.ProductCategory
	query := `SELECT id, name, created_at, updated_at FROM product_categories WHERE id = $1`

	err := r.db.GetContext(ctx, &result, query, id)
	if err != nil {
		return model.ProductCategory{}, apperror.HandleDBError(err, "failed to find product category by id")
	}
	return result, nil
}

func (r *repository) UpdateCategory(ctx context.Context, data model.ProductCategory) error {
	query := `UPDATE product_categories SET name = $1, updated_at = NOW() WHERE id = $2`

	result, err := r.db.ExecContext(ctx, query, data.Name, data.ID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return apperror.New(apperror.ErrCodeConflict, "category name is already in use by another category")
		}
		return apperror.HandleDBError(err, "failed to update category")
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return apperror.New(apperror.ErrCodeNotFound, "category not found for update")
	}
	return nil
}

func (r *repository) CountProductsByCategory(ctx context.Context, categoryID int) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM products WHERE category_id = $1`
	err := r.db.GetContext(ctx, &count, query, categoryID)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *repository) DeleteCategory(ctx context.Context, categoryID int) error {
	query := `DELETE FROM product_categories WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, categoryID)
	return err
}