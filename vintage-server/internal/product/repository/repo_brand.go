package repository

import (
	"context"
	"vintage-server/internal/product/domain"
	"vintage-server/internal/shared/db"
)

type BrandRepository interface {
	CreateBrand(ctx context.Context, brand domain.Brand) (domain.Brand, error)
	UpdateBrand(ctx context.Context, brand domain.Brand) (domain.Brand, error)
	DeleteBrand(ctx context.Context, brandID int) error
	ReadBrands(ctx context.Context, brandID *int) ([]domain.Brand, error)

	CountProducts(ctx context.Context, brand domain.Brand) (int64, error)
}

type brandRepoPostgres struct {
	db db.DBTX
}


func NewBrandRepositoryPostgres(db db.DBTX) BrandRepository {
	return &brandRepoPostgres{
		db: db,
	}
}
