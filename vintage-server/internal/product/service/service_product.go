package service

import (
	"context"
	"mime/multipart"
	"vintage-server/internal/product/domain"
	"vintage-server/internal/product/repository"
	"vintage-server/pkg/uploader"

	"github.com/google/uuid"
)

type ProductService interface {
	CreateProduct(
		ctx context.Context,
		data domain.Product,
		accountID uuid.UUID,
		thumbnail multipart.File,
		pictures []multipart.File) (domain.Product, error)
	UpdateProduct(ctx context.Context, accountID uuid.UUID) (domain.Product, error)
	FindProductByID(ctx context.Context, productID uuid.UUID) (domain.Product, error)
}

type productSvc struct {
	store    repository.ProductStore
	uploader uploader.Uploader
}

// CreateProduct implements ProductService.
func (s *productSvc) CreateProduct(
	ctx context.Context,
	data domain.Product,
	accountID uuid.UUID,
	thumbnail multipart.File,
	pictures []multipart.File) (domain.Product, error) {

	panic("unimplemented")
}

// FindProductByID implements ProductService.
func (s *productSvc) FindProductByID(ctx context.Context, productID uuid.UUID) (domain.Product, error) {
	panic("unimplemented")
}

// UpdateProduct implements ProductService.
func (s *productSvc) UpdateProduct(ctx context.Context, accountID uuid.UUID) (domain.Product, error) {
	panic("unimplemented")
}

func NewProductService(store repository.ProductStore, uploader uploader.Uploader) ProductService {
	return &productSvc{
		store:    store,
		uploader: uploader,
	}
}
