package service

import (
	"context"
	"fmt"
	"vintage-server/internal/product/domain"
	"vintage-server/internal/product/repository"
)

type ProductConditionService interface {
	CreateCondition(ctx context.Context, data domain.ProductCondition) (domain.ProductCondition, error)
	FindConditions(ctx context.Context, id *int16) ([]domain.ProductCondition, error)
	UpdateCondition(ctx context.Context, id int16, data domain.ProductCondition) (domain.ProductCondition, error)
	DeleteCondition(ctx context.Context, id int16) error
}

type productConditionService struct {
	store repository.ProductStore
}

// CreateCondition implements ProductConditionService.
func (s *productConditionService) CreateCondition(ctx context.Context, data domain.ProductCondition) (domain.ProductCondition, error) {
	cond, err := s.store.GetProductConditionRepository().CreateCondition(ctx, data)
	if err != nil {
		return domain.ProductCondition{}, fmt.Errorf("service_error: gagal membuat product condition: %w", err)
	}
	return cond, nil
}

// DeleteCondition implements ProductConditionService.
func (s *productConditionService) DeleteCondition(ctx context.Context, id int16) error {
	err := s.store.GetProductConditionRepository().DeleteCondition(ctx, id)
	if err != nil {
		return fmt.Errorf("service_error: gagal menghapus product condition %d: %w", id, err)
	}
	return nil
}

// FindConditions implements ProductConditionService.
func (s *productConditionService) FindConditions(ctx context.Context, id *int16) ([]domain.ProductCondition, error) {
	conds, err := s.store.GetProductConditionRepository().FindProductConditions(ctx, id)
	if err != nil {
		if id != nil {
			return nil, fmt.Errorf("service_error: gagal mengambil product condition dengan id %d: %w", *id, err)
		}
		return nil, fmt.Errorf("service_error: gagal mengambil semua product conditions: %w", err)
	}

	// Kalau mencari ID tertentu tapi tidak ditemukan
	if id != nil && len(conds) == 0 {
		return nil, fmt.Errorf("service_error: product condition dengan id %d tidak ditemukan", *id)
	}

	return conds, nil
}

// UpdateCondition implements ProductConditionService.
func (s *productConditionService) UpdateCondition(ctx context.Context, id int16, data domain.ProductCondition) (domain.ProductCondition, error) {
	data.ID = id
	err := s.store.GetProductConditionRepository().UpdateCondition(ctx, data)
	if err != nil {
		return domain.ProductCondition{}, fmt.Errorf("service_error: gagal memperbarui product condition %d: %w", id, err)
	}

	// Ambil lagi data terbaru dari DB
	updated, err := s.FindConditions(ctx, &id)
	if err != nil {
		return domain.ProductCondition{}, fmt.Errorf("service_error: gagal mengambil data terbaru product condition %d: %w", id, err)
	}

	return updated[0], nil
}

// Constructor
func NewProductConditionService(store repository.ProductStore) ProductConditionService {
	return &productConditionService{
		store: store,
	}
}
