package admin

import "vintage-server/dto"

type AuthService interface {
	Login(input dto.InputAdminLoginDTO) (dto.ResponseAdminLogin, error)
	Register(input dto.InputAdminRegisterDTO) (dto.ResponseAdminRegister, error)
}
