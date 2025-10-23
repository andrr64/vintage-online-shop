package service

import (
	"context"
	"vintage-server/internal/product/domain"
)

type ProductConditionService interface {
	// -- PRODUCT CONDITION MANAGEMENT --
	CreateCondition(ctx context.Context, data domain.ProductCondition) (domain.ProductCondition, error)
	FindAllConditions(ctx context.Context) (domain.ProductCondition, error)
	FindConditionByID(ctx context.Context, id int16) (domain.ProductCondition, error)
	UpdateCondition(ctx context.Context, id int16, data domain.ProductCondition) (domain.ProductCondition, error)
	DeleteCondition(ctx context.Context, id int16) error
}
