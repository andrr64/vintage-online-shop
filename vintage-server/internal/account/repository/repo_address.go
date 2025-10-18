package repository

import (
	"context"
	"vintage-server/internal/account/domain"
	"vintage-server/internal/shared/db"

	"github.com/google/uuid"
)

type AddressRepository interface {
	WithTx(tx db.DBTX) AddressRepository

	AddAddress(ctx context.Context, address domain.Address) error
	RemoveAddress(ctx context.Context, accountID uuid.UUID, addressID int64) error
	GetAddresses(ctx context.Context, accountID uuid.UUID, addressID *int64) ([]domain.Address, error)
	SetPrimaryAddress(ctx context.Context, accountID uuid.UUID, addressID int64) error
	UpdateAddress(ctx context.Context, address domain.Address) (domain.Address, error)
}

func NewAddressRepositoryPostgres(db db.DBTX) AddressRepository {
	return &addressRepoPostgres{
		db: db,
	}
}

type addressRepoPostgres struct {
	db db.DBTX
}

