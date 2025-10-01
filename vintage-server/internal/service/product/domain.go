package product

import (
	"context"
	"vintage-server/internal/model"

	"github.com/gin-gonic/gin"
)

type Service interface {
	CreateCategory(ctx context.Context, req ProductCategory) error
	UpdateCategory(ctx context.Context, req ProductCategory) error
	FindAllCategories(ctx context.Context) ([]ProductCategory, error)
	FindById(ctx context.Context, id int) (ProductCategory, error)
	DeleteCategory(ctx context.Context, id int) error
}

type Repository interface {
	CreateCategory(ctx context.Context, data model.ProductCategory) error
	FindAllCategories(ctx context.Context) ([]model.ProductCategory, error)
	FindById(ctx context.Context, id int) (model.ProductCategory, error)
	UpdateCategory(ctx context.Context, data model.ProductCategory) error
	CountProductsByCategory(ctx context.Context, categoryID int) (int, error)
	DeleteCategory(ctx context.Context, categoryID int) error
}

type ProductHandler interface {
	// Manajemen Kategori
	CreateCategory(c *gin.Context)
	ReadCategories(c *gin.Context)
	UpdateCategory(c *gin.Context)
	DeleteCategory(c *gin.Context)
}
