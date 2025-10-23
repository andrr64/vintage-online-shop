package service

import (
	"context"
	"vintage-server/internal/account/dto"
)

type RegisterService interface {
	RegisterAs(c context.Context, role string, req dto.RegisterRequest) (dto.AccountResponse, error)
}
