package product

import (
	"context"
	"database/sql"
	"errors"
	"strings"
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

// --- CATEGORY MANAGEMENT ---
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

// --- BRAND MANAGEMENT ---
func (r *repository) CreateBrand(ctx context.Context, data model.Brand) (model.Brand, error) {
	var createdBrand model.Brand
	query := `INSERT INTO brands (name, logo_url) VALUES ($1, $2) RETURNING id, name, logo_url, created_at, updated_at`

	err := r.db.QueryRowxContext(ctx, query, data.Name, data.LogoURL).StructScan(&createdBrand)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return model.Brand{}, apperror.New(apperror.ErrCodeConflict, "brand with this name already exists")
		}
		return model.Brand{}, apperror.HandleDBError(err, "failed to create brand")
	}
	return createdBrand, nil
}

func (r *repository) FindAllBrands(ctx context.Context) ([]model.Brand, error) {
	var brands []model.Brand
	query := `SELECT id, name, logo_url, created_at, updated_at FROM brands ORDER BY name ASC`

	err := r.db.SelectContext(ctx, &brands, query)
	if err != nil {
		return nil, apperror.HandleDBError(err, "failed to find all brands")
	}
	return brands, nil
}

func (r *repository) FindBrandByID(ctx context.Context, id int) (model.Brand, error) {
	var brand model.Brand
	query := `SELECT id, name, logo_url, created_at, updated_at FROM brands WHERE id = $1`

	err := r.db.GetContext(ctx, &brand, query, id)
	if err != nil {
		// Cek dulu apakah errornya karena tidak ketemu
		if errors.Is(err, sql.ErrNoRows) {
			// Jika ya, kembalikan error aslinya agar bisa dikenali service
			return model.Brand{}, err
		}
		// Jika error lain (koneksi putus, dll), baru bungkus dengan apperror
		return model.Brand{}, apperror.HandleDBError(err, "failed to find brand by id")
	}
	return brand, nil
}

func (r *repository) UpdateBrand(ctx context.Context, data model.Brand) error {
	query := `UPDATE brands SET name = $1, logo_url = $2 WHERE id = $3`

	result, err := r.db.ExecContext(ctx, query, data.Name, data.LogoURL, data.ID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return apperror.New(apperror.ErrCodeConflict, "brand name is already in use by another brand")
		}
		return apperror.HandleDBError(err, "failed to update brand")
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return apperror.New(apperror.ErrCodeNotFound, "brand not found for update")
	}
	return nil
}

func (r *repository) DeleteBrand(ctx context.Context, id int) error {
	query := `DELETE FROM brands WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return apperror.HandleDBError(err, "failed to delete brand")
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return apperror.New(apperror.ErrCodeNotFound, "brand with this id not found")
	}
	return nil
}

func (r *repository) CountProductsByBrand(ctx context.Context, brandID int) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM products WHERE brand_id = $1`
	err := r.db.GetContext(ctx, &count, query, brandID)
	if err != nil {
		return 0, apperror.HandleDBError(err, "failed to count products by brand")
	}
	return count, nil
}

// -- PRODUCT CONDITION MANAGEMENT --
func (r *repository) CreateCondition(ctx context.Context, data model.ProductCondition) (model.ProductCondition, error) {
	var createdCondition model.ProductCondition
	query := `INSERT INTO product_conditions (name) VALUES ($1) RETURNING id, name, created_at, updated_at`

	err := r.db.QueryRowxContext(ctx, query, data.Name).StructScan(&createdCondition)
	if err != nil {
		if strings.Contains(err.Error(), "product_condition_name_lower_idx"){
			return model.ProductCondition{}, apperror.New(apperror.ErrCodeConflict, "Condition already exists, try another condition name")
		}
		return model.ProductCondition{}, apperror.HandleDBError(err, "failed to create product condition")
	}

	return createdCondition, nil
}

func (r *repository) FindAllConditions(ctx context.Context) ([]model.ProductCondition, error) {
	var conditions []model.ProductCondition
	query := `SELECT id, name, created_at, updated_at FROM product_conditions ORDER BY name ASC`

	err := r.db.SelectContext(ctx, &conditions, query)
	if err != nil {
		return nil, apperror.HandleDBError(err, "failed to find all product conditions")
	}

	return conditions, nil
}

func (r *repository) FindConditionByID(ctx context.Context, id int16) (model.ProductCondition, error) {
	var condition model.ProductCondition
	query := `SELECT id, name, created_at, updated_at FROM product_conditions WHERE id = $1`

	err := r.db.GetContext(ctx, &condition, query, id)
	if err != nil {
		// Cek spesifik untuk ErrNoRows agar service bisa menanganinya
		if errors.Is(err, sql.ErrNoRows) {
			return model.ProductCondition{}, err
		}
		return model.ProductCondition{}, apperror.HandleDBError(err, "failed to find product condition by id")
	}

	return condition, nil
}

func (r *repository) UpdateCondition(ctx context.Context, data model.ProductCondition) (model.ProductCondition, error) {
	var updatedCondition model.ProductCondition
	query := `UPDATE product_conditions SET name = $1, updated_at = NOW() WHERE id = $2 RETURNING id, name, created_at, updated_at`

	err := r.db.QueryRowxContext(ctx, query, data.Name, data.ID).StructScan(&updatedCondition)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.ProductCondition{}, err // Kembalikan ErrNoRows jika ID tidak ditemukan
		}
		return model.ProductCondition{}, apperror.HandleDBError(err, "failed to update product condition")
	}
	return updatedCondition, nil
}

func (r *repository) DeleteCondition(ctx context.Context, id int16) error {
	query := `DELETE FROM product_conditions WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return apperror.HandleDBError(err, "failed to delete product condition")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return apperror.HandleDBError(err, "failed to check rows affected on delete condition")
	}

	// Jika tidak ada baris yang terhapus, berarti ID tidak ditemukan
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *repository) CountProductsByCondition(ctx context.Context, conditionID int16) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM products WHERE condition_id = $1` // Asumsi nama kolom FK adalah 'condition_id'

	err := r.db.GetContext(ctx, &count, query, conditionID)
	if err != nil {
		return 0, apperror.HandleDBError(err, "failed to count products by condition")
	}

	return count, nil
}