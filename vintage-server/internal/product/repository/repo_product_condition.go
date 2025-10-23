package repository

import (
	"context"
	"vintage-server/internal/product/domain"
	"vintage-server/internal/shared/db"
)

type ProductConditionRepository interface {
	CreateCondition(ctx context.Context, data domain.ProductCondition) (domain.ProductCondition, error)
	FindProductConditions(ctx context.Context, id *int16) ([]domain.ProductCondition, error)
	UpdateCondition(ctx context.Context, data domain.ProductCondition) error
	DeleteCondition(ctx context.Context, id int16) error
	CountProductsByCondition(ctx context.Context, conditionID int16) (int64, error)
}

type productConditionPostgres struct {
	db db.DBTX
}

func NewProductConditionRepoPostgres(db db.DBTX) ProductConditionRepository {
	return &productConditionPostgres{
		db: db,
	}
}
