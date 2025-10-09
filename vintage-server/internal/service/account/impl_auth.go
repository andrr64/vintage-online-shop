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

// loginWithRole adalah fungsi umum untuk login berdasarkan role.
func (s *accountService) loginWithRole(ctx context.Context, req account.LoginRequest, roleName string) (account.LoginResponse, error) {
	var acc model.Account
	var err error

	// Cek apakah identifier berupa email atau username
	if strings.Contains(req.Identifier, "@") {
		acc, err = s.store.FindAccountByEmailWithRole(ctx, req.Identifier, roleName)
	} else {
		acc, err = s.store.FindAccountByUsernameWithRole(ctx, req.Identifier, roleName)
	}

	// Tangani error pencarian akun
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("[WARN] Akun dengan role '%s' tidak ditemukan untuk identifier '%s'", roleName, req.Identifier)
			return account.LoginResponse{}, apperror.New(apperror.ErrCodeUnauthorized, "invalid credentials")
		}
		log.Printf("[ERROR] Gagal mencari akun '%s': %v", roleName, err)
		return account.LoginResponse{}, apperror.New(apperror.ErrCodeInternal, "an internal error occurred")
	}

	// Verifikasi password
	if err := hash.Verify(acc.Password, req.Password); err != nil {
		log.Printf("[WARN] Password salah untuk akun '%s'", acc.Email)
		return account.LoginResponse{}, apperror.New(apperror.ErrCodeUnauthorized, "invalid credentials")
	}

	// Generate JWT token
	token, err := s.jwt.GenerateToken(acc.ID, roleName)
	if err != nil {
		log.Printf("[ERROR] Gagal generate token untuk akun '%s': %v", acc.Email, err)
		return account.LoginResponse{}, apperror.New(apperror.ErrCodeInternal, "failed to generate token")
	}

	// Return hasil login
	return account.LoginResponse{
		AccessToken: token,
		UserProfile: account.ConvertAccountToUserProfileResponse(&acc),
	}, nil
}

// LoginCustomer authenticates a customer user
func (s *accountService) LoginCustomer(ctx context.Context, req account.LoginRequest) (account.LoginResponse, error) {
	return s.loginWithRole(ctx, req, "customer")
}

// LoginAdmin authenticates an admin user
func (s *accountService) LoginAdmin(ctx context.Context, req account.LoginRequest) (account.LoginResponse, error) {
	return s.loginWithRole(ctx, req, "admin")
}

// LoginSeller authenticates a seller user
func (s *accountService) LoginSeller(ctx context.Context, req account.LoginRequest) (account.LoginResponse, error) {
	return s.loginWithRole(ctx, req, "seller")
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
