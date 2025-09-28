package user

import (
	"time"

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

type LoginRequest struct {
	Identifier string `json:"identifier" binding:"required"` // bisa username / email
	Password   string `json:"password" binding:"required"`
}

type UserProfileResponse struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Firstname string    `json:"firstname"`
	Lastname  *string   `json:"lastname"`
	Email     string    `json:"email"`
	AvatarURL *string   `json:"avatar_url"`
}

type LoginResponse struct {
	AccessToken string              `json:"access_token"`
	UserProfile UserProfileResponse `json:"user"`
}
