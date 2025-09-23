package admin

import "vintage-server/dto"


type AuthService interface {
	Login(input dto.InputAdminLoginDTO) (string, error)
	Register(input dto.InputAdminRegisterDTO) (dto.ResponseAdminRegister, error)
}