package service

import (
	"context"
	"vintage-server/internal/account/domain"
	"vintage-server/internal/account/dto"
	"vintage-server/internal/account/repository"
	"vintage-server/pkg/auth"
)

type registerServiceImpl struct {
	store repository.AccountStore
	jwt   *auth.JWTService
}

// RegisterAs implements RegisterService.
func (r *registerServiceImpl) RegisterAs(c context.Context, role string, req dto.RegisterRequest) (dto.AccountResponse, error) {
	var createdAccount domain.Account
	err := r.store.ExecTx(c, func(store repository.AccountStore) error {
		acc, err := store.GetAccountRepo().CreateAccount(c, createdAccount)
		if err != nil {
			return err
		}

		if err := store.GetAccountRepo().AssignRole(c, acc.ID, role); err != nil {
			return err
		}

		createdAccount = acc
		return nil
	})
	if err != nil {
		return dto.AccountResponse{}, err
	}
	return dto.ConvertAccountnDomainToDTO(createdAccount), nil
}

func NewRegisterService(store repository.AccountStore, jwtSecret string) RegisterService {
	return &registerServiceImpl{
		store: store,
		jwt:   auth.NewJWTService(jwtSecret),
	}
}
