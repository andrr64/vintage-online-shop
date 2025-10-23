package service

import (
	"context"
	"vintage-server/internal/product/domain"
	"vintage-server/internal/product/repository"
)

type productConditionService struct {
	store repository.ProductStore
}

// CreateCondition implements ProductConditionService.
func (p *productConditionService) CreateCondition(ctx context.Context, data domain.ProductCondition) (domain.ProductCondition, error) {
	panic("unimplemented")
}

// DeleteCondition implements ProductConditionService.
func (p *productConditionService) DeleteCondition(ctx context.Context, id int16) error {
	panic("unimplemented")
}

// FindAllConditions implements ProductConditionService.
func (p *productConditionService) FindAllConditions(ctx context.Context) (domain.ProductCondition, error) {
	panic("unimplemented")
}

// FindConditionByID implements ProductConditionService.
func (p *productConditionService) FindConditionByID(ctx context.Context, id int16) (domain.ProductCondition, error) {
	panic("unimplemented")
}

// UpdateCondition implements ProductConditionService.
func (p *productConditionService) UpdateCondition(ctx context.Context, id int16, data domain.ProductCondition) (domain.ProductCondition, error) {
	panic("unimplemented")
}

func NewProductConditionService(store repository.ProductStore) ProductConditionService {
	return &productConditionService{
		store: store,
	}
}
