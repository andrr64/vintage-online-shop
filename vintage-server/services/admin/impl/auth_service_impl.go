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

func NewAdminAuthService(repo repositories.AdminRepository) adminService.AuthService {
	return &authService{repo: repo}
}

// Method Login sekarang sesuai interface
func (s *authService) Login(input dto.InputAdminLoginDTO) (string, error) {
	// Contoh implementasi dummy
	if input.Username == "admin" && input.Password == "admin123" {
		randomString := "token-xyz"
		return randomString, nil
	}

	return "", fmt.Errorf("invalid credentials")
}

func (s *authService) Register(input dto.InputAdminRegisterDTO) (dto.ResponseAdminRegister, error) {
	username := input.Username
	password := input.Password

	hashedPassword, error := helpers.GeneratePasswordHash(password)
	if error != nil {
		return dto.ResponseAdminRegister{}, error
	}

	result, error := s.repo.CreateAccount(username, hashedPassword)

	if error != nil {
		return dto.ResponseAdminRegister{}, nil
	}

	return dto.ResponseAdminRegister{
		Username: result.Username,
		Status:   true,
	}, nil
}
