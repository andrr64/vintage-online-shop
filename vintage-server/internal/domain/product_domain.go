package product

import (
	"context"
	"vintage-server/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ProductRepository interface {
	WithTx(tx *sqlx.Tx) ProductRepository // <-- TAMBAHKAN INI

	// -- CATEGORY MANAGEMENT --
	CreateCategory(ctx context.Context, data model.ProductCategory) (model.ProductCategory, error)
	FindAllCategories(ctx context.Context) ([]model.ProductCategory, error)
	FindCategoryById(ctx context.Context, id int) (model.ProductCategory, error)
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
	UpdateCondition(ctx context.Context, data model.ProductCondition) error
	DeleteCondition(ctx context.Context, id int16) error
	CountProductsByCondition(ctx context.Context, conditionID int16) (int, error)

	// -- SHOP MANAGEMENT --
	FindShopByAccountID(ctx context.Context, accountID uuid.UUID) (model.Shop, error)

	// -- PRODUCT IMAGES MANAGEMENT --
	CreateProductImage(ctx context.Context, image model.ProductImage) (model.ProductImage, error)
	DeleteProductImage(ctx context.Context, imageID int64) error
	FindProductImageByURL(ctx context.Context, imageURL string) (model.ProductImage, error)
	FindImagesByProductID(ctx context.Context, productID uuid.UUID) ([]model.ProductImage, error)
	UpdateProductImageIndex(ctx context.Context, imageURL string, newIndex int16) error

	// -- PRODUCT MANAGEMENT
	CreateProduct(ctx context.Context, product model.Product) (model.Product, error)
	CreateProductImages(ctx context.Context, images []model.ProductImage) error
	UpdateProduct(ctx context.Context, p model.Product) (model.Product, error)
	FindProductByIDAndShop(ctx context.Context, productID, shopID uuid.UUID) (model.Product, error)
	FindProductByID(ctx context.Context, productID uuid.UUID) (model.Product, error)

	// -- SIZE
	CreateProductSize(ctx context.Context, productSize model.ProductSize) (model.ProductSize, error)
}

type ProductService interface {
	// -- CATEGORY MANAGEMENT --
	CreateCategory(ctx context.Context, req ProductCategoryDTO) error
	UpdateCategory(ctx context.Context, req ProductCategoryDTO) error
	FindAllCategories(ctx context.Context) ([]ProductCategoryDTO, error)
	FindById(ctx context.Context, id int) (ProductCategoryDTO, error)
	DeleteCategory(ctx context.Context, id int) error

	// -- BRAND MANAGEMENT --
	CreateBrand(ctx context.Context, req BrandRequest) (model.Brand, error)
	FindAllBrands(ctx context.Context) ([]ProductBrandDTO, error)
	FindBrandByID(ctx context.Context, id int) (model.Brand, error)
	UpdateBrand(ctx context.Context, id int, req BrandRequest) error
	DeleteBrand(ctx context.Context, id int) error

	// -- PRODUCT CONDITION MANAGEMENT --
	CreateCondition(ctx context.Context, req ProductConditionRequest) (model.ProductCondition, error)
	FindAllConditions(ctx context.Context) ([]model.ProductCondition, error)
	FindConditionByID(ctx context.Context, id int16) (model.ProductCondition, error)
	UpdateCondition(ctx context.Context, id int16, req ProductConditionRequest) (model.ProductCondition, error)
	DeleteCondition(ctx context.Context, id int16) error

	// -- PRODUCT MANAGEMENT
	CreateProduct(ctx context.Context, accountID uuid.UUID, request CreateProductRequest) (ProductDTO, error)
	UpdateProduct(ctx context.Context, accountID uuid.UUID, request UpdateProductDTO) (ProductDTO, error)
	FindProductByID(ctx context.Context, productID uuid.UUID) (ProductDTO, error)

	// -- SIZE
	CreateProductSize(ctx context.Context, request ProductConditionRequest) (ProductSizeDTO, error)
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
	UpdateProduct(c *gin.Context)
	GetProuctByID(c *gin.Context)

	// -- PRODUCT SIZE MANAGEMENT --
	CreateProductSize(c *gin.Context)
}
