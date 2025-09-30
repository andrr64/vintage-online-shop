package account

import (
	"time"
	"vintage-server/internal/model"

	"github.com/google/uuid"
)

// WishlistItemDetail adalah DTO untuk menampilkan wishlist beserta detail produk.
// Ini didefinisikan di sini agar Repository tahu bentuk data apa yang harus dikembalikan.
type WishlistItemDetail struct {
	ProductID       int64     `json:"product_id" `
	ProductName     string    `json:"product_name" `
	ProductPrice    int64     `json:"product_price" `
	ProductImageURL string    `json:"product_image_url"`
	AddedAt         time.Time `json:"added_at" `
}

type UserAddress struct {
	ID             int64     `json:"id"`
	Label          string    `json:"label"`
	DistrictID     string    `json:"district_id"`
	RegencyID      string    `json:"regency_id"`
	ProvinceID     string    `json:"province_id"`
	VillageID      string    `json:"village_id"`
	RecipientName  string    `json:"recipient_name"`
	RecipientPhone string    `json:"recipient_phone"`
	Street         string    `json:"street"`
	PostalCode     string    `json:"postal_code"`
	IsPrimary      bool      `json:"is_primary"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// -- REQUEST --
type RegisterRequest struct {
	Username  string  `json:"username" binding:"required"`
	Firstname string  `json:"firstname" binding:"required"`
	Lastname  *string `json:"lastname"`
	Email     string  `json:"email" binding:"required,email"`
	Password  string  `json:"password" binding:"required,min=8"`
}

type UpdateProfileRequest struct {
	Firstname string  `json:"firstname" binding:"required"`
	Lastname  *string `json:"lastname"`
}

type AddAddressRequest struct {
	DistrictID     string `json:"district_id"`
	RegencyID      string `json:"regency_id"`
	ProvinceID     string `json:"province_id"`
	VillageID      string `json:"village_id"`
	Label          string `json:"label"`
	RecipientName  string `json:"recipient_name"`
	RecipientPhone string `json:"recipient_phone"`
	Street         string `json:"street"`
	PostalCode     string `json:"postal_code"`
}

type LoginRequest struct {
	Identifier string `json:"identifier" binding:"required"` // bisa username / email
	Password   string `json:"password" binding:"required"`
}

// -- RESPONSE --
type UserProfileResponse struct {
	Username  string  `json:"username"`
	Firstname string  `json:"firstname"`
	Lastname  *string `json:"lastname"`
	Email     string  `json:"email"`
	AvatarURL *string `json:"avatar_url"`
}

type LoginResponse struct {
	AccessToken string              `json:"access_token"`
	UserProfile UserProfileResponse `json:"user"`
}

type AddressIdentifier struct {
	AddressID int64 `json:"address_id" binding:"required"`
}

// -- FUNCTIONS --

func ConvertAccountToUserProfileResponse(acc *model.Account) UserProfileResponse {
	if acc == nil {
		return UserProfileResponse{}
	}
	return UserProfileResponse{
		Username:  acc.Username,
		Firstname: acc.Firstname,
		Lastname:  acc.Lastname,
		Email:     acc.Email,
		AvatarURL: acc.AvatarURL,
	}
}

func NewAddressFromRequest(accountID uuid.UUID, req AddAddressRequest, isPrimary bool) model.Address {
	now := time.Now()

	return model.Address{
		AccountID:      accountID,
		DistrictID:     req.DistrictID,
		RegencyID:      req.RegencyID,
		ProvinceID:     req.ProvinceID,
		VillageID:      req.VillageID,
		Label:          req.Label,
		RecipientName:  req.RecipientName,
		RecipientPhone: req.RecipientPhone,
		Street:         req.Street,
		PostalCode:     req.PostalCode,
		IsPrimary:      isPrimary,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

func ConvertAddressToDTO(addr model.Address) UserAddress {
	return UserAddress{
		ID:             addr.ID,
		Label:          addr.Label,
		DistrictID:     addr.DistrictID,
		RegencyID:      addr.RegencyID,
		ProvinceID:     addr.ProvinceID,
		VillageID:      addr.VillageID,
		RecipientName:  addr.RecipientName,
		RecipientPhone: addr.RecipientPhone,
		Street:         addr.Street,
		PostalCode:     addr.PostalCode,
		IsPrimary:      addr.IsPrimary,
		CreatedAt:      addr.CreatedAt,
		UpdatedAt:      addr.UpdatedAt,
	}
}

func ConvertAddressesToDTO(addresses []model.Address) []UserAddress {
	result := make([]UserAddress, len(addresses))
	for i, addr := range addresses {
		result[i] = ConvertAddressToDTO(addr)
	}
	return result
}
