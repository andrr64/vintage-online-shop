package services

import (
	"vintage-server/dto"
	"vintage-server/models"
)

// interface
type UserService interface {
	Register(input dto.InputRegisterDTO) (dto.ResponseRegisterUserDTO, error)
	Login(input dto.InputLoginDTO) (string, error) // return JWT token
	FindByID(userId uint) (*models.User, error)
	UpdateAccount(userID uint, data dto.InputUpdateAccountDTO) (dto.ResponseUserInfoDTO, error)
	UpdatePassword(userId uint, data dto.InputUpdatePasswordDTO)  (int, error)
}
