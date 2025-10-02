package product

import (
	"context"
	"vintage-server/internal/model"

	"github.com/gin-gonic/gin"
)

type Service interface {
	// -- BRAND MANAGEMENT --
	CreateCategory(ctx context.Context, req ProductCategory) error
	UpdateCategory(ctx context.Context, req ProductCategory) error
	FindAllCategories(ctx context.Context) ([]ProductCategory, error)
	FindById(ctx context.Context, id int) (ProductCategory, error)
	DeleteCategory(ctx context.Context, id int) error

	// -- BRAND MANAGEMENT --
	CreateBrand(ctx context.Context, req CreateBrandRequest) (model.Brand, error)
	FindAllBrands(ctx context.Context) ([]model.Brand, error)
	FindBrandByID(ctx context.Context, id int) (model.Brand, error)
	UpdateBrand(ctx context.Context, id int, req UpdateBrandRequest) error
	DeleteBrand(ctx context.Context, id int) error
}

type Repository interface {
	// -- CATEGORY MANAGEMENT --
	CreateCategory(ctx context.Context, data model.ProductCategory) error
	FindAllCategories(ctx context.Context) ([]model.ProductCategory, error)
	FindById(ctx context.Context, id int) (model.ProductCategory, error)
	UpdateCategory(ctx context.Context, data model.ProductCategory) error
	DeleteCategory(ctx context.Context, categoryID int) error
	CountProductsByCategory(ctx context.Context, categoryID int) (int, error)

	// -- BRAND MANAGEMENT --
	CreateBrand(ctx context.Context, data model.Brand) (model.Brand, error)
	FindAllBrands(ctx context.Context) ([]model.Brand, error)
	FindBrandByID(ctx context.Context, id int) (model.Brand, error)
	UpdateBrand(ctx context.Context, data model.Brand) error
	DeleteBrand(ctx context.Context, id int) error
	CountProductsByBrand(ctx context.Context, brandID int) (int, error)
}

type ProductHandler interface {
	// Manajemen Kategori
	CreateCategory(c *gin.Context)
	ReadCategories(c *gin.Context)
	UpdateCategory(c *gin.Context)
	DeleteCategory(c *gin.Context)

	// -- brand management --
	CreateBrand(c *gin.Context)
	ReadBrand(c *gin.Context)
	UpdateBrand(c *gin.Context)
	DeleteBrand(c *gin.Context)
}
