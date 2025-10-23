package repository

import (
	"context"
	"vintage-server/internal/product/domain"
	"vintage-server/internal/shared/db"
)

type ProductCategoryRepository interface {
	// -- CATEGORY MANAGEMENT --
	CreateCategory(ctx context.Context, data domain.ProductCategory) (domain.ProductCategory, error)
	FindAllCategories(ctx context.Context, id *int) ([]domain.ProductCategory, error)
	UpdateCategory(ctx context.Context, data domain.ProductCategory) error
	DeleteCategory(ctx context.Context, categoryID int) error
	CountProductsByCategory(ctx context.Context, categoryID int) (int64, error)
}

type postgresImpl struct {
	db db.DBTX
}

func NewProductCategoryRepoPostgres(db db.DBTX) ProductCategoryRepository {
	return &postgresImpl{
		db: db,
	}
}
