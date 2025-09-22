package services

import (
	"vintage-server/dto"
)

// interface
type UserService interface {
	Register(input dto.RegisterUserDTO) (dto.ResponseRegisterUserDTO, error)
	Login(input dto.LoginUserDTO) (string, error) // return JWT token
}
