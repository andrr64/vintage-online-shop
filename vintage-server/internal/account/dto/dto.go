package dto

import (
	"vintage-server/internal/account/domain"

	"github.com/google/uuid"
)

type LoginRequest struct {
	Identifier string `json:"identifier" binding:"required"` // Bisa email atau username
	Password   string `json:"password" binding:"required"`
}

// -- REQUEST --
type RegisterRequest struct {
	Username  string  `json:"username" binding:"required"`
	Firstname string  `json:"firstname" binding:"required"`
	Lastname  *string `json:"lastname"`
	Email     string  `json:"email" binding:"required,email"`
	Password  string  `json:"password" binding:"required,min=8"`
}

type LoginResponse struct {
	AccessToken string          `json:"access_token"`
	Account     AccountResponse `json:"account"`
}

type AccountResponse struct {
	Username  string  `json:"username"`
	Firstname string  `json:"firstname"`
	Lastname  string  `json:"lastname"`
	AvatarURL *string `json:"avatar_url"`
}

type AddressDTO struct {
	ID             int64  `json:"address_id"`
	DistrictID     string `json:"district_id" binding:"required"`
	RegencyID      string `json:"regency_id" binding:"required"`
	ProvinceID     string `json:"province_id" binding:"required"`
	VillageID      string `json:"village_id" binding:"required"`
	Label          string `json:"label" binding:"required"`
	RecipientName  string `json:"recipient_name" binding:"required"`
	RecipientPhone string `json:"recipient_phone" binding:"required"`
	Street         string `json:"street" binding:"required"`
	PostalCode     string `json:"postal_code" binding:"required"`
}

func CreateLoginResponse(token string, acc domain.Account) LoginResponse {
	return LoginResponse{
		AccessToken: token,
		Account:     ConvertAccountnDomainToDTO(acc),
	}
}

func ConvertAddressDTOToDomain(adrs AddressDTO, accountID uuid.UUID) domain.Address {
	return domain.Address{
		ID:             adrs.ID,
		AccountID:      accountID,
		DistrictID:     adrs.DistrictID,
		RegencyID:      adrs.RegencyID,
		ProvinceID:     adrs.ProvinceID,
		VillageID:      adrs.VillageID,
		Label:          adrs.Label,
		RecipientName:  adrs.RecipientName,
		RecipientPhone: adrs.RecipientPhone,
		Street:         adrs.Street,
		PostalCode:     adrs.PostalCode,
		IsPrimary:      false, // default false saat pembuatan awal
	}
}

func ConvertAddressDomainToDTO(address domain.Address) AddressDTO {
	return AddressDTO{
		ID:             address.ID,
		DistrictID:     address.DistrictID,
		RegencyID:      address.RegencyID,
		ProvinceID:     address.ProvinceID,
		VillageID:      address.VillageID,
		Label:          address.Label,
		RecipientName:  address.RecipientName,
		RecipientPhone: address.RecipientPhone,
		Street:         address.Street,
		PostalCode:     address.PostalCode,
	}
}

func ConvertAccountnDomainToDTO(acc domain.Account) AccountResponse {
	lastname := ""
	if acc.Lastname != nil {
		lastname = *acc.Lastname
	}

	return AccountResponse{
		Username:  acc.Username,
		Firstname: acc.Firstname,
		Lastname:  lastname,
		AvatarURL: acc.AvatarURL,
	}
}
