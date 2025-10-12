package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	product "vintage-server/internal/domain"
	"vintage-server/internal/model"
	"vintage-server/internal/shared/db"
	"vintage-server/pkg/apperror"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type productRepository struct {
	db *sqlx.DB
	tx *sqlx.Tx
}

// FindProductByID implements product.ProductRepository.
func (r *productRepository) FindProductByID(ctx context.Context, productID uuid.UUID) (model.Product, error) {
	// 🔹 Ambil data utama produk
	query := `
		SELECT 
			p.id, p.shop_id, p.condition_id, p.category_id, p.brand_id, p.size_id,
			p.name, p.summary, p.description, p.price, p.stock,
			p.created_at, p.updated_at
		FROM products p
		WHERE p.id = $1
	`
	var product model.Product
	if err := r.db.GetContext(ctx, &product, query, productID); err != nil {
		return model.Product{}, err
	}

	// 🔹 Ambil relasi: Brand
	if product.BrandID != nil {
		var brand model.Brand
		err := r.db.GetContext(ctx, &brand, `SELECT id, name, logo_url, created_at, updated_at FROM brands WHERE id = $1`, *product.BrandID)
		if err == nil {
			product.Brand = &brand
		}
	}

	// 🔹 Ambil relasi: Category
	var category model.ProductCategory
	err := r.db.GetContext(ctx, &category, `SELECT id, name, created_at, updated_at FROM product_categories WHERE id = $1`, product.CategoryID)
	if err == nil {
		product.Category = &category
	}

	// 🔹 Ambil relasi: Condition
	var condition model.ProductCondition
	err = r.db.GetContext(ctx, &condition, `SELECT id, name, created_at, updated_at FROM product_conditions WHERE id = $1`, product.ConditionID)
	if err == nil {
		product.Condition = &condition
	}

	// 🔹 Ambil relasi: Size
	if product.SizeID != nil {
		var size model.ProductSize
		err := r.db.GetContext(ctx, &size, `SELECT id, name FROM product_size WHERE id = $1`, *product.SizeID)
		if err == nil {
			product.Size = &size
		}
	}

	// 🔹 Ambil relasi: Shop
	var shop model.Shop
	err = r.db.GetContext(ctx, &shop, `
		SELECT id, account_id, name, summary, description, active, created_at, updated_at
		FROM shop WHERE id = $1
	`, product.ShopID)
	if err == nil {
		product.Shop = &shop
	}

	// 🔹 Ambil relasi: Images
	var images []model.ProductImage
	err = r.db.SelectContext(ctx, &images, `
		SELECT id, product_id, image_index, url, created_at
		FROM product_images
		WHERE product_id = $1
		ORDER BY image_index ASC
	`, product.ID)
	if err == nil {
		product.Images = images
	}

	return product, nil
}

// UpdateProductImageIndex implements product.ProductRepository.
func (r *productRepository) UpdateProductImageIndex(ctx context.Context, imageURL string, newIndex int16) error {
	const query = `
		UPDATE public.product_images
		SET image_index = $1
		WHERE url = $2
	`

	result, err := r.db.ExecContext(ctx, query, newIndex, imageURL)
	if err != nil {
		return fmt.Errorf("failed to update product image index: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("failed cant find images: %w", err)
	}

	return nil
}

// CreateProductImage implements product.ProductRepository.
func (r *productRepository) CreateProductImage(ctx context.Context, image model.ProductImage) (model.ProductImage, error) {
	const query = `
		INSERT INTO public.product_images (product_id, image_index, url)
		VALUES ($1, $2, $3)
		RETURNING id, product_id, image_index, url, created_at
	`

	var created model.ProductImage
	err := r.db.QueryRowContext(
		ctx,
		query,
		image.ProductID,
		image.ImageIndex,
		image.URL,
	).Scan(
		&created.ID,
		&created.ProductID,
		&created.ImageIndex,
		&created.URL,
		&created.CreatedAt,
	)
	return created, err
}

// FindProductImageByURL implements product.ProductRepository.
func (r *productRepository) FindProductImageByURL(ctx context.Context, imageURL string) (model.ProductImage, error) {
	var image model.ProductImage
	query := `SELECT * FROM product_images WHERE url = $1`
	err := r.db.GetContext(ctx, &image, query, imageURL)
	return image, err
}

// FindProductWithImagesByProductID implements product.ProductRepository.
// FindProductWithImagesByProductID implements product.ProductRepository.
func (r *productRepository) FindImagesByProductID(ctx context.Context, productID uuid.UUID) ([]model.ProductImage, error) {
	var images []model.ProductImage

	query := `
		SELECT id, product_id, url, image_index
		FROM product_images
		WHERE product_id = $1
		ORDER BY image_index ASC
	`

	err := r.db.SelectContext(ctx, &images, query, productID)
	if err != nil {
		return nil, err
	}

	return images, nil
}

// CreateProductSize implements product.ProductRepository.
// CreateProductSize implements product.ProductRepository.
func (r *productRepository) CreateProductSize(ctx context.Context, productSize model.ProductSize) (model.ProductSize, error) {
	var createdSize model.ProductSize
	query := `
		INSERT INTO product_size (name)
		VALUES ($1)
		RETURNING id, name
	`
	err := r.db.GetContext(ctx, &createdSize, query, productSize.Name)
	if err != nil {
		return model.ProductSize{}, err
	}
	return createdSize, nil
}

// -- PRIVATE FUNCTION
func (r *productRepository) GetQuerier() db.DBTX {
	if r.tx != nil {
		return r.tx
	}
	return r.db
}

// -- PRIVATE FUNCTION --
func (r *productRepository) GetCount(ctx context.Context, query string, args ...interface{}) (int, error) {

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
func (r *productRepository) CreateBrand(ctx context.Context, data model.Brand) (model.Brand, error) {
	var brand model.Brand
	query := `INSERT INTO brands (name, logo_url) VALUES ($1, $2)
			  RETURNING id, name, logo_url, created_at, updated_at`
	err := r.db.QueryRowxContext(ctx, query, data.Name, data.LogoURL).StructScan(&brand)
	return brand, err
}

func (r *productRepository) CreateCategory(ctx context.Context, data model.ProductCategory) (model.ProductCategory, error) {
	var category model.ProductCategory
	query := `INSERT INTO product_categories (name)
			  VALUES ($1) RETURNING id, name, created_at, updated_at`
	err := r.db.QueryRowxContext(ctx, query, data.Name).StructScan(&category)
	return category, err
}

func (r *productRepository) CreateCondition(ctx context.Context, data model.ProductCondition) (model.ProductCondition, error) {
	var condition model.ProductCondition
	query := `INSERT INTO product_conditions (name)
			  VALUES ($1) RETURNING id, name, created_at, updated_at`
	err := r.db.QueryRowxContext(ctx, query, data.Name).StructScan(&condition)
	return condition, err
}

// CreateProduct implements product.ProductRepository.
func (r *productRepository) CreateProduct(ctx context.Context, p model.Product) (model.Product, error) {
	query := `
		INSERT INTO public.products
			(shop_id, condition_id, category_id, size_id, brand_id, name, summary, description, price, stock)
		VALUES
			(:shop_id, :condition_id, :category_id, :size_id, :brand_id, :name, :summary, :description, :price, :stock)
		RETURNING *
	`

	// sqlx.NamedExec + Get bisa pakai NamedQuery untuk map ke struct
	rows, err := r.db.NamedQueryContext(ctx, query, p)
	if err != nil {
		return model.Product{}, err
	}
	defer rows.Close()

	if rows.Next() {
		var created model.Product
		if err := rows.StructScan(&created); err != nil {
			return model.Product{}, err
		}
		return created, nil
	}

	return model.Product{}, errors.New("failed to insert product")
}

// CreateProductImages implements product.ProductRepository.
// CreateProductImages implements product.ProductRepository.
func (r *productRepository) CreateProductImages(ctx context.Context, images []model.ProductImage) error {
	if len(images) == 0 {
		return nil // tidak ada yang diinsert
	}

	query := `
		INSERT INTO public.product_images (product_id, image_index, url)
		VALUES (:product_id, :image_index, :url)
	`

	// NamedExec dengan slice struct akan melakukan batch insert otomatis
	_, err := r.db.NamedExecContext(ctx, query, images)
	if err != nil {
		return err
	}

	return nil
}

// -- DELETE FUNCTION --
func (r *productRepository) DeleteBrand(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM brands WHERE id = $1`, id)
	return err
}

func (r *productRepository) DeleteCategory(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM product_categories WHERE id = $1`, id)
	return err
}

func (r *productRepository) DeleteCondition(ctx context.Context, id int16) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM product_conditions WHERE id = $1`, id)
	return err
}

// FindAllBrands implements product.ProductRepository.
func (r *productRepository) FindAllBrands(ctx context.Context) ([]model.Brand, error) {

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

func (r *productRepository) DeleteProductImage(ctx context.Context, imgID int64) error {
	query := "DELETE FROM product_images WHERE id = $1"
	_, err := r.db.ExecContext(ctx, query, imgID)
	return err
}

// FindProductByIDAndShop implements product.ProductRepository.
func (r *productRepository) FindProductByIDAndShop(ctx context.Context, productID, shopID uuid.UUID) (model.Product, error) {
	var product model.Product

	query := `
		SELECT 
			p.id, p.shop_id, p.condition_id, p.category_id, p.size_id, p.brand_id,
			p.name, p.summary, p.description, p.price, p.stock,
			p.created_at, p.updated_at
		FROM products p
		WHERE p.id = $1 AND p.shop_id = $2
	`

	err := r.GetQuerier().GetContext(ctx, &product, query, productID, shopID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Product{}, apperror.New(apperror.ErrCodeNotFound, "product not found or unauthorized")
		}
		return model.Product{}, apperror.HandleDBError(err, "failed to find product by id and shop")
	}

	return product, nil
}

func (r *productRepository) UpdateProduct(ctx context.Context, p model.Product) (model.Product, error) {
	query := `
		UPDATE products
		SET 
			name = $1,
			category_id = $2,
			condition_id = $3,
			brand_id = $4,
			size_id = $5,
			price = $6,
			stock = $7,
			summary = $8,
			description = $9,
			updated_at = NOW()
		WHERE id = $10
		RETURNING id, shop_id, condition_id, category_id, brand_id, size_id,
				  name, summary, description, price, stock,
				  created_at, updated_at
	`

	var updated model.Product
	err := r.GetQuerier().QueryRowContext(ctx, query,
		p.Name,
		p.CategoryID,
		p.ConditionID,
		p.BrandID,
		p.SizeID,
		p.Price,
		p.Stock,
		p.Summary,
		p.Description,
		p.ID,
	).Scan(
		&updated.ID,
		&updated.ShopID,
		&updated.ConditionID,
		&updated.CategoryID,
		&updated.BrandID,
		&updated.SizeID,
		&updated.Name,
		&updated.Summary,
		&updated.Description,
		&updated.Price,
		&updated.Stock,
		&updated.CreatedAt,
		&updated.UpdatedAt,
	)

	if err != nil {
		return model.Product{}, apperror.HandleDBError(err, "failed to update product")
	}

	return updated, nil
}

// FindConditionByID implements product.ProductRepository.
func (r *productRepository) FindConditionByID(ctx context.Context, id int16) (model.ProductCondition, error) {
	var condition model.ProductCondition
	query := `SELECT * FROM product_conditions WHERE id = $1`
	err := r.db.GetContext(ctx, &condition, query, id)
	return condition, err
}

// FindShopByAccountID implements product.ProductRepository.
func (r *productRepository) FindShopByAccountID(ctx context.Context, accountID uuid.UUID) (model.Shop, error) {
	var shop model.Shop
	query := "SELECT * FROM shop WHERE account_id = $1"
	err := r.db.GetContext(ctx, &shop, query, accountID) // <- kirim accountID sebagai parameter $1
	return shop, err
}

// UpdateBrand implements product.ProductRepository.
func (r *productRepository) UpdateBrand(ctx context.Context, brand model.Brand) error {
	query := `
		UPDATE brands
		SET name = $1, logo_url = $2, updated_at = NOW()
		WHERE id = $3
	`
	_, err := r.GetQuerier().ExecContext(ctx, query, brand.Name, brand.LogoURL, brand.ID)
	return err
}

// UpdateCategory implements product.ProductRepository.
// UpdateCategory implements product.ProductRepository.
func (r *productRepository) UpdateCategory(ctx context.Context, data model.ProductCategory) error {
	query := `UPDATE product_categories 
	          SET name = $1, updated_at = NOW()
	          WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, data.Name, data.ID)
	return err
}

// UpdateCondition implements product.ProductRepository.
func (r *productRepository) UpdateCondition(ctx context.Context, data model.ProductCondition) error {
	query := `
		UPDATE product_conditions
		SET name = $1, updated_at = NOW()
		WHERE id = $2D
	`
	_, err := r.db.ExecContext(ctx, query, data.Name, data.ID)
	return err
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
