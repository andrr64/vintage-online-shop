package service

import (
	"context"
	"fmt"
	"vintage-server/internal/account/domain"
	"vintage-server/internal/account/dto"
	"vintage-server/internal/account/repository"

	"github.com/google/uuid"
)

type addressServiceImpl struct {
	store repository.AccountStore
}

// UpdateAddress implements AddressService.
func (a *addressServiceImpl) UpdateAddress(ctx context.Context, address domain.Address) (dto.AddressDTO, error) {
	panic("unimplemented")
}

// SetPrimaryAddress implements AddressService.
func (a *addressServiceImpl) SetPrimaryAddress(ctx context.Context, accountID uuid.UUID, addressID int64) error {
	err := a.store.GetAddressRepo().SetPrimaryAddress(ctx, accountID, addressID)
	if err != nil {
		// handle: errror like logging etc
	}
	return err
}

// GetAddresses implements AddressService.
func (a *addressServiceImpl) GetAddresses(ctx context.Context, accountID uuid.UUID, addressID *int64) ([]dto.AddressDTO, error) {
	addresses, err := a.store.GetAddressRepo().GetAddresses(ctx, accountID, addressID)
	if err != nil {
		return nil, fmt.Errorf("failed to get addresses: %w", err)
	}

	if len(addresses) == 0 {
		// inisialisasi slice kosong agar JSON jadi [] bukan null
		return []dto.AddressDTO{}, nil
	}

	addressesDTO := make([]dto.AddressDTO, 0, len(addresses))
	for _, addr := range addresses {
		addressesDTO = append(addressesDTO, dto.ConvertAddressDomainToDTO(addr))
	}
	return addressesDTO, nil
}

// AddAddress implements AddressService.
func (a *addressServiceImpl) AddAddress(ctx context.Context, accountID uuid.UUID, address domain.Address) error {
	if err := a.store.GetAddressRepo().AddAddress(ctx, address); err != nil {
		///TODO: handle error
		return err
	}
	return nil
}

// DeleteAddress implements AddressService.
func (a *addressServiceImpl) DeleteAddress(ctx context.Context, accountID uuid.UUID, addressID int64) error {
	if err := a.store.GetAddressRepo().RemoveAddress(ctx, accountID, addressID); err != nil {
		///TODO: handle error
		return err
	}
	return nil
}

func NewAddressService(store repository.AccountStore) AddressService {
	return &addressServiceImpl{
		store: store,
	}
}
