package service

import (
	"context"
	"strings"
	"vintage-server/internal/account/domain"
	"vintage-server/internal/account/dto"
	"vintage-server/internal/account/repository"
	"vintage-server/pkg/auth"
	serviceerror "vintage-server/pkg/service_error"
	"vintage-server/pkg/uploader"
)

type loginServiceImpl struct {
	store    repository.AccountStore
	jwt      *auth.JWTService
	uploader uploader.Uploader
}

// RegisterAs implements RegisterService.
func (l *loginServiceImpl) RegisterAs(c context.Context, role string, req dto.RegisterRequest) {
	panic("unimplemented")
}

// LoginAs implements LoginService.
func (l *loginServiceImpl) LoginAs(ctx context.Context, role string, req dto.LoginRequest) (dto.LoginResponse, error) {
	var account domain.Account
	var err error

	if strings.Contains(req.Identifier, "@") {
		email, e := domain.NewEmail(req.Identifier)
		if e != nil {
			return dto.LoginResponse{}, serviceerror.New(serviceerror.ErrBadRequest, "Invalid email.", e)
		}
		account, err = l.store.GetAccountRepo().FindAccountByEmailAndRoleString(ctx, email, role)
	} else {
		username, e := domain.NewUsername(req.Identifier)
		if e != nil {
			return dto.LoginResponse{}, serviceerror.New(serviceerror.ErrBadRequest, "Invalid username.", e)
		}
		account, err = l.store.GetAccountRepo().FindAccountByUsernameAndRoleString(ctx, username, role)
	}

	if err != nil {
		return dto.LoginResponse{}, err
	}

	if !account.ComparePassword(req.Password) {
		return dto.LoginResponse{}, serviceerror.New(serviceerror.ErrUnauthorized, "Identifier or password is wrong.", nil)
	}

	token, err := l.jwt.GenerateToken(account.ID, role)
	if err != nil {
		return dto.LoginResponse{}, serviceerror.New(serviceerror.ErrInternal, "Internal server error.", err)
	}

	return dto.CreateLoginResponse(token, account), nil
}

func NewLoginService(store repository.AccountStore, jwtSecret string, uploader uploader.Uploader) LoginService {
	return &loginServiceImpl{
		store:    store,
		jwt:      auth.NewJWTService(jwtSecret),
		uploader: uploader,
	}
}
