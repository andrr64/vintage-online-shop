package services

import (
	"vintage-server/dto"
	"vintage-server/models"
)

// interface
type UserService interface {
	Register(input dto.RegisterUserDTO) (dto.ResponseRegisterUserDTO, error)
	Login(input dto.LoginUserDTO) (string, error) // return JWT token
	FindByID(userId uint) (*models.User, error)
}