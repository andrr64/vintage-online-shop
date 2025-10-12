package service

import (
	"context"
	"vintage-server/internal/account/dto"
)

type LoginService interface {
	LoginAs(ctx context.Context, role string, req dto.LoginRequest) (dto.LoginResponse, error)
}
