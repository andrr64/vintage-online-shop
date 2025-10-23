package service

import (
	"context"
	"vintage-server/internal/account/domain"
	"vintage-server/internal/account/dto"

	"github.com/google/uuid"
)

type AddressService interface {
	AddAddress(ctx context.Context, accountID uuid.UUID, address domain.Address) error
	DeleteAddress(ctx context.Context, accountID uuid.UUID, addressID int64) error
	GetAddresses(ctx context.Context, accountID uuid.UUID, addressID *int64) ([]dto.AddressDTO, error)
	SetPrimaryAddress(ctx context.Context, accountID uuid.UUID, addressID int64) error
	UpdateAddress(ctx context.Context, address domain.Address) (dto.AddressDTO, error)
}
