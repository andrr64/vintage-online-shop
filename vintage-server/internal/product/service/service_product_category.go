package service

import (
	"context"
	"fmt"
	"vintage-server/internal/product/domain"
	"vintage-server/internal/product/repository"
)

type ProductCategoryService interface {
	CreateCategory(ctx context.Context, data domain.ProductCategory) (domain.ProductCategory, error)
	FindCategories(ctx context.Context, id *int) ([]domain.ProductCategory, error)
	UpdateCategory(ctx context.Context, data domain.ProductCategory) (domain.ProductCategory, error)
	DeleteCategory(ctx context.Context, categoryID int) error
	CountProductsByCategory(ctx context.Context, categoryID int) (int64, error)
}

type productCategorySvc struct {
	store repository.ProductStore
}

// CreateCategory implements ProductCategoryService.
func (s *productCategorySvc) CreateCategory(ctx context.Context, data domain.ProductCategory) (domain.ProductCategory, error) {
	category, err := s.store.GetProductCategoryRepo().CreateCategory(ctx, data)
	if err != nil {
		return domain.ProductCategory{}, fmt.Errorf("service_error: failed to create product category: %w", err)
	}
	return category, nil
}

// FindCategories implements ProductCategoryService.
func (s *productCategorySvc) FindCategories(ctx context.Context, id *int) ([]domain.ProductCategory, error) {
	categories, err := s.store.GetProductCategoryRepo().FindAllCategories(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("service_error: failed to fetch product categories: %w", err)
	}

	if id != nil && len(categories) == 0 {
		return categories, fmt.Errorf("service_error: product category with id %d not found", *id)
	}

	return categories, nil
}

// UpdateCategory implements ProductCategoryService.
func (s *productCategorySvc) UpdateCategory(ctx context.Context, data domain.ProductCategory) (domain.ProductCategory, error) {
	// Update the category
	err := s.store.GetProductCategoryRepo().UpdateCategory(ctx, data)
	if err != nil {
		return domain.ProductCategory{}, fmt.Errorf("service_error: failed to update product category %d: %w", data.ID, err)
	}

	// Return the updated data from DB
	updatedList, err := s.FindCategories(ctx, &data.ID)
	if err != nil {
		return domain.ProductCategory{}, fmt.Errorf("service_error: failed to retrieve updated category %d: %w", data.ID, err)
	}

	return updatedList[0], nil
}

// DeleteCategory implements ProductCategoryService.
func (s *productCategorySvc) DeleteCategory(ctx context.Context, categoryID int) error {
	err := s.store.GetProductCategoryRepo().DeleteCategory(ctx, categoryID)
	if err != nil {
		return fmt.Errorf("service_error: failed to delete product category %d: %w", categoryID, err)
	}
	return nil
}

// CountProductsByCategory implements ProductCategoryService.
func (s *productCategorySvc) CountProductsByCategory(ctx context.Context, categoryID int) (int64, error) {
	count, err := s.store.GetProductCategoryRepo().CountProductsByCategory(ctx, categoryID)
	if err != nil {
		return 0, fmt.Errorf("service_error: failed to count products by category %d: %w", categoryID, err)
	}
	return count, nil
}

// Constructor
func NewProductCategorySvc(store repository.ProductStore) ProductCategoryService {
	return &productCategorySvc{
		store: store,
	}
}
