package service

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"strings"
	"time"
	"vintage-server/internal/domain/account"
	"vintage-server/internal/model"
	"vintage-server/pkg/apperror"

	"github.com/google/uuid"
)

// DeleteAddress deletes a user's address.
func (s *accountService) DeleteAddress(ctx context.Context, accountID uuid.UUID, addressID int64) error {
	if err := s.store.DeleteAddress(ctx, addressID, accountID); err != nil {
		log.Printf("Error deleting address: %v", err)
		return apperror.New(apperror.ErrCodeInternal, "could not delete address")
	}
	return nil
}

// GetAddressByID retrieves a single address by its ID for a specific user.
func (s *accountService) GetAddressByID(ctx context.Context, accountID uuid.UUID, addressID int64) (account.UserAddress, error) {
	addr, err := s.store.FindAddressByIDAndAccountID(ctx, addressID, accountID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return account.UserAddress{}, apperror.New(apperror.ErrCodeNotFound, "address not found")
		}
		log.Printf("Error getting address by ID: %v", err)
		return account.UserAddress{}, apperror.New(apperror.ErrCodeInternal, "an internal error occurred")
	}
	return account.ConvertAddressToDTO(addr), nil
}

func (s *accountService) GetAddressesByUserID(ctx context.Context, accountID uuid.UUID) ([]account.UserAddress, error) {
	addrs, err := s.store.FindAddressesByAccountID(ctx, accountID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []account.UserAddress{}, nil
		}
		log.Printf("Error finding addresses by account ID: %v", err)
		return nil, apperror.New(apperror.ErrCodeInternal, "an internal error occurred")
	}

	return account.ConvertAddressesToDTO(addrs), nil
}


// SetPrimaryAddress sets a specific address as the user's primary one.
func (s *accountService) SetPrimaryAddress(ctx context.Context, accountID uuid.UUID, addressID int64) error {
	if err := s.store.SetPrimaryAddress(ctx, accountID, addressID); err != nil {
		log.Printf("Error setting primary address: %v", err)
		return apperror.New(apperror.ErrCodeInternal, "could not set primary address")
	}
	return nil
}

// UpdateAddress updates an existing user address.
func (s *accountService) UpdateAddress(ctx context.Context, accountID uuid.UUID, req account.UserAddress) (account.UserAddress, error) {
	addr, err := s.store.FindAddressByIDAndAccountID(ctx, req.ID, accountID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return account.UserAddress{}, apperror.New(apperror.ErrCodeNotFound, "address not found")
		}
		log.Printf("Error finding address to update: %v", err)
		return account.UserAddress{}, apperror.New(apperror.ErrCodeInternal, "an internal error occurred")
	}

	addr.ProvinceID = req.ProvinceID
	addr.RegencyID = req.RegencyID
	addr.DistrictID = req.DistrictID
	addr.VillageID = req.VillageID
	addr.RecipientName = req.RecipientName
	addr.RecipientPhone = req.RecipientPhone
	addr.Label = req.Label
	addr.PostalCode = req.PostalCode
	addr.Street = req.Street
	addr.UpdatedAt = time.Now()

	updatedAddr, err := s.store.UpdateAddress(ctx, accountID, addr)
	if err != nil {
		log.Printf("Error updating address: %v", err)
		return account.UserAddress{}, apperror.New(apperror.ErrCodeInternal, "could not update address")
	}

	return account.ConvertAddressToDTO(updatedAddr), nil
}

// AddAddress implements account.AccountService.
func (s *accountService) AddAddress(ctx context.Context, userID uuid.UUID, req account.AddAddressRequest) (account.UserAddress, error) {
	address := account.NewAddressFromRequest(userID, req, false)

	// Optional: gunakan ExecTx jika ingin banyak operasi address dalam satu transaksi
	var savedAddress model.Address
	err := s.store.ExecTx(ctx, func(repo account.AccountRepository) error {
		var err error
		savedAddress, err = repo.SaveAddress(ctx, address)
		return err
	})
	if err != nil {
		if strings.Contains(err.Error(), "violate") {
			return account.UserAddress{}, apperror.New(apperror.ErrCodeBadRequest, "invalid address data provided")
		}
		return account.UserAddress{}, apperror.New(apperror.ErrCodeInternal, "could not save address")
	}

	return account.ConvertAddressToDTO(savedAddress), nil
}
