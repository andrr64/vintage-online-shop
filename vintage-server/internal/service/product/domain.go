package product

import (
	"context"
	"vintage-server/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Service interface {
	// -- CATEGORY MANAGEMENT --
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

	// -- PRODUCT CONDITION MANAGEMENT --
	CreateCondition(ctx context.Context, req ProductConditionRequest) (model.ProductCondition, error)
	FindAllConditions(ctx context.Context) ([]model.ProductCondition, error)
	FindConditionByID(ctx context.Context, id int16) (model.ProductCondition, error)
	UpdateCondition(ctx context.Context, id int16, req ProductConditionRequest) (model.ProductCondition, error)
	DeleteCondition(ctx context.Context, id int16) error

	// -- PRODUCT MANAGEMENT
	CreateProduct(ctx context.Context, accountID uuid.UUID, request CreateProductRequest) 
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

	// -- PRODUCT CONDITION MANAGEMENT --
	CreateCondition(ctx context.Context, data model.ProductCondition) (model.ProductCondition, error)
	FindAllConditions(ctx context.Context) ([]model.ProductCondition, error)
	FindConditionByID(ctx context.Context, id int16) (model.ProductCondition, error)
	UpdateCondition(ctx context.Context, data model.ProductCondition) (model.ProductCondition, error)
	DeleteCondition(ctx context.Context, id int16) error
	CountProductsByCondition(ctx context.Context, conditionID int16) (int, error)

	// -- SHOP MANAGEMENT --
	FindShopByAccountID(ctx context.Context, accountID uuid.UUID) (model.Shop, error)
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

	// -- PRODUCT CONDITION MANAGEMENT --
	CreateCondition(c *gin.Context)
	ReadConditions(c *gin.Context)
	UpdateCondition(c *gin.Context)
	DeleteCondition(c *gin.Context)

	// -- PRODUCT MANAGEMENT --
	CreateProduct(c *gin.Context)
}
