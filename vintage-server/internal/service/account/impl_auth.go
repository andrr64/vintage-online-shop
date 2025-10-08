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
	"vintage-server/pkg/hash"

	"github.com/google/uuid"
)

func (s *accountService) LoginCustomer(ctx context.Context, req account.LoginRequest) (account.LoginResponse, error) {
	roleName := "customer"
	var acc model.Account
	var err error

	if strings.Contains(req.Identifier, "@") {
		acc, err = s.store.FindAccountByEmailWithRole(ctx, req.Identifier, roleName)
	} else {
		acc, err = s.store.FindAccountByUsernameWithRole(ctx, req.Identifier, roleName)
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return account.LoginResponse{}, apperror.New(apperror.ErrCodeUnauthorized, "invalid credentials")
		}
		log.Printf("Error finding account for login: %v", err)
		return account.LoginResponse{}, apperror.New(apperror.ErrCodeInternal, "an internal error occurred")
	}

	if err := hash.Verify(acc.Password, req.Password); err != nil {
		return account.LoginResponse{}, apperror.New(apperror.ErrCodeUnauthorized, "invalid credentials")
	}

	token, err := s.jwt.GenerateToken(acc.ID, roleName)
	if err != nil {
		log.Printf("Error generating token: %v", err)
		return account.LoginResponse{}, apperror.New(apperror.ErrCodeInternal, "an internal error occurred")
	}

	return account.LoginResponse{
		AccessToken: token,
		UserProfile: account.ConvertAccountToUserProfileResponse(&acc),
	}, nil
}

// LoginAdmin authenticates an admin user using email or username
func (s *accountService) LoginAdmin(ctx context.Context, req account.LoginRequest) (account.LoginResponse, error) {
	role := "admin"

	// 1. Cari akun admin berdasarkan email atau username
	acc, err := func() (model.Account, error) {
		if strings.Contains(req.Identifier, "@") {
			return s.store.FindAccountByEmailWithRole(ctx, req.Identifier, role)
		}
		return s.store.FindAccountByUsernameWithRole(ctx, req.Identifier, role)
	}()

	// 2. Tangani error saat pencarian
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return account.LoginResponse{}, apperror.New(
				apperror.ErrCodeUnauthorized,
				"invalid credentials",
			)
		}
		log.Printf("Error finding admin account: %v", err)
		return account.LoginResponse{}, apperror.New(
			apperror.ErrCodeInternal,
			"an internal error occurred",
		)
	}

	// 3. Verifikasi password
	if err := hash.Verify(acc.Password, req.Password); err != nil {
		return account.LoginResponse{}, apperror.New(
			apperror.ErrCodeUnauthorized,
			"invalid credentials",
		)
	}

	// 4. Generate token JWT
	token, err := s.jwt.GenerateToken(acc.ID, role)
	if err != nil {
		log.Printf("Error generating admin token: %v", err)
		return account.LoginResponse{}, apperror.New(
			apperror.ErrCodeInternal,
			"failed to generate token",
		)
	}

	// 5. Return response
	return account.LoginResponse{
		AccessToken: token,
		UserProfile: account.ConvertAccountToUserProfileResponse(&acc),
	}, nil
}

// LoginSeller authenticates a seller user using email or username
func (s *accountService) LoginSeller(ctx context.Context, req account.LoginRequest) (account.LoginResponse, error) {
	role := "seller"

	// 1. Cari akun seller berdasarkan email atau username
	acc, err := func() (model.Account, error) {
		if strings.Contains(req.Identifier, "@") {
			return s.store.FindAccountByEmailWithRole(ctx, req.Identifier, role)
		}
		return s.store.FindAccountByUsernameWithRole(ctx, req.Identifier, role)
	}()

	// 2. Tangani error saat pencarian
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return account.LoginResponse{}, apperror.New(
				apperror.ErrCodeUnauthorized,
				"invalid credentials",
			)
		}
		log.Printf("Error finding seller account: %v", err)
		return account.LoginResponse{}, apperror.New(
			apperror.ErrCodeInternal,
			"an internal error occurred",
		)
	}

	// 3. Verifikasi password
	if err := hash.Verify(acc.Password, req.Password); err != nil {
		return account.LoginResponse{}, apperror.New(
			apperror.ErrCodeUnauthorized,
			"invalid credentials",
		)
	}

	// 4. Generate token JWT
	token, err := s.jwt.GenerateToken(acc.ID, role)
	if err != nil {
		log.Printf("Error generating seller token: %v", err)
		return account.LoginResponse{}, apperror.New(
			apperror.ErrCodeInternal,
			"failed to generate token",
		)
	}

	// 5. Return response
	return account.LoginResponse{
		AccessToken: token,
		UserProfile: account.ConvertAccountToUserProfileResponse(&acc),
	}, nil
}

// Logout tetap sederhana
func (s *accountService) Logout(ctx context.Context, userID uuid.UUID) (string, error) {
	// implementasi kedepan: invalidasi token, logging, dll
	return "OK", nil
}

// RegisterCustomer handles new user registration using ExecTx
func (s *accountService) RegisterCustomer(ctx context.Context, req account.RegisterRequest) (account.UserProfileResponse, error) {
	var createdAcc model.Account

	err := s.store.ExecTx(ctx, func(repoInTx account.AccountRepository) error {
		// 1. cek username/email
		if _, err := repoInTx.FindAccountByUsername(ctx, req.Username); err == nil {
			return apperror.New(apperror.ErrCodeConflict, "username already taken")
		} else if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return err
		}

		if _, err := repoInTx.FindAccountByEmail(ctx, req.Email); err == nil {
			return apperror.New(apperror.ErrCodeConflict, "email already taken")
		} else if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return err
		}

		// 2. hash password
		hashedPassword, err := hash.Generate(req.Password)
		if err != nil {
			return err
		}

		// 3. insert account
		newAccount := model.Account{
			Username:  req.Username,
			Email:     req.Email,
			Password:  hashedPassword,
			Firstname: req.Firstname,
			Lastname:  req.Lastname,
			Active:    true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		createdAcc, err = repoInTx.InsertAccount(ctx, newAccount)
		if err != nil {
			return err
		}

		// 4. ambil roleID & insert account_role
		roleID, err := repoInTx.GetRoleIDByName(ctx, "customer")
		if err != nil {
			return err
		}
		if err := repoInTx.InsertAccountRole(ctx, createdAcc.ID, roleID); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		if apperr, ok := err.(*apperror.AppError); ok {
			return account.UserProfileResponse{}, apperr
		}
		return account.UserProfileResponse{}, apperror.New(apperror.ErrCodeInternal, "failed to create account")
	}
	return account.ConvertAccountToUserProfileResponse(&createdAcc), nil
}
