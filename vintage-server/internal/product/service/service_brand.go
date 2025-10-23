package service

import (
	"context"
	"mime/multipart"
	"vintage-server/internal/product/domain"
)

type BrandService interface {
	CreateBrand(c context.Context, brand domain.Brand, logo multipart.File) (domain.Brand, error)
	UpdateBrand(c context.Context, brand domain.Brand, logo *multipart.File) (domain.Brand, error)
	DeleteBrand(c context.Context, brandID int) error
	ReadBrands(c context.Context, brandID *int) ([]domain.Brand, error)
}
