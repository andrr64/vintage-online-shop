package product

import (
	"context"
	"vintage-server/internal/model"

	"github.com/gin-gonic/gin"
)

type Service interface {
	CreateCategory(ctx context.Context, req ProductCategory) error
	FindAllCategories(ctx context.Context) ([]ProductCategory, error)
	FindById(ctx context.Context, id int) (ProductCategory, error)
}

type Repository interface {
	CreateCategory(ctx context.Context, data ProductCategory) error 
	FindAllCategories(ctx context.Context) ([]model.ProductCategory, error)
	FindById(ctx context.Context, id int) (model.ProductCategory, error)
}

type ProductHandler interface {
	// Manajemen Kategori
	CreateCategory(c *gin.Context)
	ReadCategories(c *gin.Context)
	UpdateCategory(c *gin.Context)
	DeleteCategory(c *gin.Context)
}
