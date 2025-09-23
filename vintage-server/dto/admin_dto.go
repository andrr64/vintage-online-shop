package dto

type InputAdminLoginDTO struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type InputAdminRegisterDTO struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ResponseAdminLogin struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

type ResponseAdminRegister struct {
	Username string `json:"username"`
	Status   bool   `json:"status"`
}
