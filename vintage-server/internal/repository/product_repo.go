package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"
	"vintage-server/internal/domain/product"
	"vintage-server/internal/model"
	"vintage-server/internal/shared/db"
	"vintage-server/pkg/apperror"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type productRepository struct {
	db *sqlx.DB
	tx *sqlx.Tx
}

const DefaultQueryTimeout = 3 * time.Second

// -- PRIVATE FUNCTION
func (r *productRepository) GetQuerier() db.DBTX {
	if r.tx != nil {
		return r.tx
	}
	return r.db
}

func (r *productRepository) WithTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, DefaultQueryTimeout)
}

// -- PRIVATE FUNCTION --
func (r *productRepository) GetCount(ctx context.Context, query string, args ...interface{}) (int, error) {
	ctx, cancel := r.WithTimeout(ctx)
	defer cancel()

	var count int
	if err := r.GetQuerier().GetContext(ctx, &count, query, args...); err != nil {
		return 0, apperror.HandleDBError(err, "count query failed")
	}
	return count, nil
}

// -- COUNT FUNCTION --
// CountProductsByBrand implements product.ProductRepository.
func (r *productRepository) CountProductsByBrand(ctx context.Context, brandID int) (int, error) {
	query := "SELECT COUNT(*) FROM products WHERE brand_id = $1"
	return r.GetCount(ctx, query, brandID)
}

// CountProductsByCategory implements product.ProductRepository.
func (r *productRepository) CountProductsByCategory(ctx context.Context, categoryID int) (int, error) {
	query := "SELECT COUNT(*) FROM products WHERE category_id = $1"
	return r.GetCount(ctx, query, categoryID)
}

// CountProductsByCondition implements product.ProductRepository.
func (r *productRepository) CountProductsByCondition(ctx context.Context, conditionID int16) (int, error) {
	query := "SELECT COUNT(*) FROM products WHERE condition_id = $1"
	return r.GetCount(ctx, query, conditionID)
}

// -- CREATE FUNCTION --
// CreateBrand implements product.ProductRepository.
func (r *productRepository) CreateBrand(ctx context.Context, data model.Brand) (model.Brand, error) {
	ctx, cancel := r.WithTimeout(ctx)
	defer cancel()

	var brand model.Brand
	query := `INSERT INTO brands (name, logo_url) VALUES ($1, $2) RETURNING id, name, logo_url, created_at, updated_at`
	err := r.db.QueryRowxContext(ctx, query, data.Name, data.LogoURL).StructScan(&brand)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return model.Brand{}, apperror.New(apperror.ErrCodeConflict, "brand with this name already exists")
		}
		return model.Brand{}, apperror.HandleDBError(err, "failed to create brand")
	}
	return brand, nil
}

// CreateCategory implements product.ProductRepository.
func (r *productRepository) CreateCategory(ctx context.Context, data model.ProductCategory) (model.ProductCategory, error) {
	ctx, cancel := r.WithTimeout(ctx)
	defer cancel()

	var category model.ProductCategory

	query := "INSERT INTO product_categories (name) VALUES ($1)"

	err := r.db.QueryRowxContext(ctx, query, data.Name).StructScan(&category)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return model.ProductCategory{}, apperror.New(apperror.ErrCodeConflict, "category with this name already exists")
		}
		return model.ProductCategory{}, apperror.HandleDBError(err, "failed to create category")
	}
	return category, nil
}

// CreateCondition implements product.ProductRepository.
func (r *productRepository) CreateCondition(ctx context.Context, data model.ProductCondition) (model.ProductCondition, error) {
	ctx, cancel := r.WithTimeout(ctx)
	defer cancel()

	var condition model.ProductCondition

	query := "INSERT INTO product_condition (name) VALUES ($1)"
	err := r.db.QueryRowxContext(ctx, query, data.Name).StructScan(&condition)
	if err != nil {
		if strings.Contains(err.Error(), "product_condition_name_lower_idx") {
			return model.ProductCondition{}, apperror.New(apperror.ErrCodeConflict, "Condition already exists, try another condition name")
		}
		return model.ProductCondition{}, apperror.HandleDBError(err, "failed to create product condition")
	}
	return condition, nil
}

// CreateProduct implements product.ProductRepository.
func (r *productRepository) CreateProduct(ctx context.Context, product model.Product) (model.Product, error) {
	panic("unimplemented")
}

// CreateProductImages implements product.ProductRepository.
func (r *productRepository) CreateProductImages(ctx context.Context, images []model.ProductImage) error {
	panic("unimplemented")
}

// -- DELETE FUNCTION --
// DeleteBrand implements product.ProductRepository.
func (r *productRepository) DeleteBrand(ctx context.Context, id int) error {
	ctx, cancel := r.WithTimeout(ctx)
	defer cancel()

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

// DeleteCategory implements product.ProductRepository.
func (r *productRepository) DeleteCategory(ctx context.Context, categoryID int) error {
	ctx, cancel := r.WithTimeout(ctx)
	defer cancel()
	query := "DELETE FROM product_categories WHERE id = $1"

	result, err := r.db.ExecContext(ctx, query, categoryID)
	if err != nil {
		return apperror.HandleDBError(err, "failed to delete category")
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return apperror.New(apperror.ErrCodeNotFound, "category not found")
	}
	return nil
}

// DeleteCondition implements product.ProductRepository.
func (r *productRepository) DeleteCondition(ctx context.Context, id int16) error {
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
		return apperror.New(apperror.ErrCodeNotFound, "category not found")
	}
	return nil
}

// FindAllBrands implements product.ProductRepository.
func (r *productRepository) FindAllBrands(ctx context.Context) ([]model.Brand, error) {
	ctx, cancel := r.WithTimeout(ctx)
	defer cancel()

	var brands []model.Brand

	query := "SELECT * FROM brands ORDER BY name ASC"
	err := r.db.SelectContext(ctx, &brands, query)
	if err != nil {
		return nil, apperror.HandleDBError(err, "failed to find all brands")
	}
	return brands, nil
}

// FindAllCategories implements product.ProductRepository.
func (r *productRepository) FindAllCategories(ctx context.Context) ([]model.ProductCategory, error) {
	ctx, cancel := r.WithTimeout(ctx)
	defer cancel()

	var categories []model.ProductCategory
	query := `SELECT * FROM product_categories ORDER BY name`
	err := r.db.SelectContext(ctx, &categories, query)
	if err != nil {
		return nil, apperror.HandleDBError(err, "failed to find all brands")
	}
	return categories, nil
}

// FindAllConditions implements product.ProductRepository.
func (r *productRepository) FindAllConditions(ctx context.Context) ([]model.ProductCondition, error) {
	ctx, cancel := r.WithTimeout(ctx)
	defer cancel()

	var conditions []model.ProductCondition

	query := "SELECT * FROM product_conditions ORDER BY name"
	err := r.db.SelectContext(ctx, &conditions, query)
	if err != nil {
		return nil, apperror.HandleDBError(err, "failed to find all product conditions")
	}
	return conditions, nil
}

// FindBrandByID implements product.ProductRepository.
func (r *productRepository) FindBrandByID(ctx context.Context, id int) (model.Brand, error) {
	var brand model.Brand
	query := `SELECT * FROM brands WHERE id = $1`

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

// FindCategoryById implements product.ProductRepository.
func (r *productRepository) FindCategoryById(ctx context.Context, id int) (model.ProductCategory, error) {
	ctx, cancel := r.WithTimeout(ctx)
	defer cancel()

	var result model.ProductCategory
	query := `
		SELECT *
		FROM product_categories
		WHERE id = $1
	`
	err := r.db.GetContext(ctx, &result, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.ProductCategory{}, apperror.New(apperror.ErrCodeNotFound, "category not found")
		}
		return model.ProductCategory{}, apperror.HandleDBError(err, "failed to find product category by id")
	}
	return result, nil
}

// FindConditionByID implements product.ProductRepository.
func (r *productRepository) FindConditionByID(ctx context.Context, id int16) (model.ProductCondition, error) {
	var condition model.ProductCondition
	query := `SELECT * FROM product_conditions WHERE id = $1`

	err := r.db.GetContext(ctx, &condition, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.ProductCondition{}, err
		}
		return model.ProductCondition{}, apperror.HandleDBError(err, "failed to find product condition by id")
	}

	return condition, nil
}

// FindShopByAccountID implements product.ProductRepository.
func (r *productRepository) FindShopByAccountID(ctx context.Context, accountID uuid.UUID) (model.Shop, error) {
	var shop model.Shop
	query := "SELECT * FROM shop WHERE account_id = $1"

	if err := r.db.GetContext(ctx, &shop, query); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Shop{}, err
		}
		return model.Shop{}, apperror.HandleDBError(err, "failed to find product condition by id")
	}
	return shop, nil
}

// UpdateBrand implements product.ProductRepository.
func (r *productRepository) UpdateBrand(ctx context.Context, data model.Brand) error {
	panic("unimplemented")
}

// UpdateCategory implements product.ProductRepository.
func (r *productRepository) UpdateCategory(ctx context.Context, data model.ProductCategory) error {
	panic("unimplemented")
}

// UpdateCondition implements product.ProductRepository.
func (r *productRepository) UpdateCondition(ctx context.Context, data model.ProductCondition) (model.ProductCondition, error) {
	panic("unimplemented")
}

// WithTx implements product.ProductRepository.
func (r *productRepository) WithTx(tx *sqlx.Tx) product.ProductRepository {
	return &productRepository{db: r.db, tx: tx}
}

func NewProductRepository(db *sqlx.DB) product.ProductRepository {
	return &productRepository{
		db: db,
	}
}
