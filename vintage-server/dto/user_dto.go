package dto

// DTO untuk register user
type RegisterUserDTO struct {
	Username string `json:"username" binding:"required,min=4"`
	Fullname string `json:"fullname" binding:"required,min=4"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=10"`
}

// DTO untuk login
type LoginUserDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// DTO untuk response setelah register
type ResponseRegisterUserDTO struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Fullname string `json:"fullname"`
	Email    string `json:"email"`
}
