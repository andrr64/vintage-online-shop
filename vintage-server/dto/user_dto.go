package dto

// DTO untuk register user
type InputRegisterDTO struct {
	Username string `json:"username" binding:"required,min=4"`
	Fullname string `json:"fullname" binding:"required,min=4"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=10"`
}

// DTO untuk login
type InputLoginDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type InputUpdateAccountDTO struct {
	Username string `json:"username"` // opsional
	Fullname string `json:"fullname"`
	Email    string `json:"email"` // opsional
}

type InputUpdatePasswordDTO struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

// DTO untuk response setelah register
type ResponseRegisterUserDTO struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Fullname string `json:"fullname"`
	Email    string `json:"email"`
}

type ResponseUserInfoDTO struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Fullname string `json:"fullname"`
}