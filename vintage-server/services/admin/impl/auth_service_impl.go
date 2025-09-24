package impl

import (
	"errors"
	"fmt"
	"vintage-server/dto"
	"vintage-server/helpers"
	"vintage-server/repositories"
	"vintage-server/security/jwt"
	adminService "vintage-server/services/admin"
)

type authService struct {
	repo repositories.AdminRepository
}

// NewAdminAuthService buat instance service
func NewAdminAuthService(repo repositories.AdminRepository) adminService.AuthService {
	return &authService{repo: repo}
}

// Login sesuai interface AuthService
func (s *authService) Login(input dto.InputAdminLoginDTO) (dto.ResponseAdminLogin, error) {

	admin, err := s.repo.FindByUsername(input.Username)

	if err != nil {
		return dto.ResponseAdminLogin{}, err
	}

	if admin == nil {
		return dto.ResponseAdminLogin{}, fmt.Errorf("user not found")
	}

	if !helpers.ComparePassword(input.Password, admin.Password) {
		return dto.ResponseAdminLogin{}, fmt.Errorf("username or password is wrong")
	}

	return dto.ResponseAdminLogin{
		Username: admin.Username,
		Token:    jwt.CreateAdminAccessToken(admin.ID, admin.Username),
	}, nil
}

// Register buat mendaftarkan admin baru
func (s *authService) Register(input dto.InputAdminRegisterDTO) (dto.ResponseAdminRegister, error) {
	username := input.Username
	password := input.Password

	// Hash password
	hashedPassword, err := helpers.GeneratePasswordHash(password)
	if err != nil {
		return dto.ResponseAdminRegister{}, err
	}

	exists, err := s.repo.IsExists(username)
	if err != nil {
		return dto.ResponseAdminRegister{}, err
	}
	if exists {
		return dto.ResponseAdminRegister{}, errors.New("username already exists")
	}

	// Simpan ke repository
	result, err := s.repo.CreateAccount(username, hashedPassword)
	if err != nil {
		return dto.ResponseAdminRegister{}, err
	}

	return dto.ResponseAdminRegister{
		Username: result.Username,
		Status:   true,
	}, nil
}
