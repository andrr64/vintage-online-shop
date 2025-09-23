package impl

import (
	"fmt"
	"vintage-server/dto"
	"vintage-server/helpers"
	"vintage-server/repositories"
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
func (s *authService) Login(input dto.InputAdminLoginDTO) (string, error) {
	// Dummy login
	if input.Username == "admin" && input.Password == "admin123" {
		randomToken := "token-xyz"
		return randomToken, nil
	}

	return "", fmt.Errorf("invalid credentials")
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
