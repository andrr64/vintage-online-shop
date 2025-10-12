package dto

import "vintage-server/internal/account/domain"

type LoginRequest struct {
	Identifier string `json:"identifier" binding:"required"` // Bisa email atau username
	Password   string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken string  `json:"access_token"`
	Account     Account `json:"account"`
}

type Account struct {
	Username  string  `json:"username"`
	Firstname string  `json:"firstname"`
	Lastname  string  `json:"lastname"`
	AvatarURL *string `json:"avatar_url"`
}

func CreateLoginResponse(token string, acc domain.Account) LoginResponse {
	return LoginResponse{
		AccessToken: token,
		Account:     ConvertDomainToAccount(acc),
	}
}

func ConvertDomainToAccount(acc domain.Account) Account {
	lastname := ""
	if acc.Lastname != nil {
		lastname = *acc.Lastname
	}

	return Account{
		Username:  acc.Username,
		Firstname: acc.Firstname,
		Lastname:  lastname,
		AvatarURL: acc.AvatarURL,
	}
}
