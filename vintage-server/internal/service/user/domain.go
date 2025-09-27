package user

// File: internal/service/user/domain.go

import (
	"context"
	"vintage-server/internal/model" // Sesuaikan dengan path proyekmu
)

// =================================================================================
// KONTRAK UNTUK SERVICE (Logika Bisnis) 🧠
// =================================================================================
type Service interface {
	// --- User & Authentication ---
	// Usecase: CustomerRegister
	Register(ctx context.Context, req RegisterRequest) (UserProfileResponse, error)

	// Usecase: CustomerLogin, SellerLogin, AdminLogin
	Login(ctx context.Context, req LoginRequest) (LoginResponse, error)

	// Usecase: AdminManage Users
	DeactivateUser(ctx context.Context, userID int64, reason string) error
	GetUserProfile(ctx context.Context, userID int64) (model.Account, error)

	// --- Address Management ---
	// Usecase: CustomerManage Addresses
	AddAddress(ctx context.Context, userID int64, address model.Address) (model.Address, error)
	GetAddressesByUserID(ctx context.Context, userID int64) ([]model.Address, error)
	UpdateAddress(ctx context.Context, userID, addressID int64, address model.Address) (model.Address, error)
	DeleteAddress(ctx context.Context, userID, addressID int64) error

	// Usecase: CustomerSet Primary Address
	SetPrimaryAddress(ctx context.Context, userID int64, addressID int64) error

	// --- Wishlist Management ---
	// Usecase: CustomerAdd/View/Remove Wishlist
	AddToWishlist(ctx context.Context, userID, productID int64) error
	GetWishlistByUserID(ctx context.Context, userID int64) ([]WishlistItemDetail, error)
	RemoveFromWishlist(ctx context.Context, userID, productID int64) error
}

// =================================================================================
// KONTRAK UNTUK REPOSITORY (Akses Database) 🚚
// =================================================================================
// Repository mendefinisikan semua interaksi ke database yang dibutuhkan oleh Service.
type Repository interface {
	// --- Account ---
	FindAccountByID(ctx context.Context, id int64) (model.Account, error)
	FindAccountByEmail(ctx context.Context, email string) (model.Account, error)
	FindAccountByUsername(ctx context.Context, username string) (model.Account, error)
	SaveAccount(ctx context.Context, account model.Account) (model.Account, error)
	UpdateAccount(ctx context.Context, account model.Account) error

	// --- Address ---
	SaveAddress(ctx context.Context, address model.Address) (model.Address, error)
	FindAddressesByAccountID(ctx context.Context, accountID int64) ([]model.Address, error)
	FindAddressByIDAndAccountID(ctx context.Context, addressID, accountID int64) (model.Address, error)
	UpdateAddress(ctx context.Context, address model.Address) (model.Address, error)
	DeleteAddress(ctx context.Context, addressID, accountID int64) error
	// TransactionSetPrimaryAddress akan menangani 2 query (unset old, set new) dalam satu transaksi DB
	TransactionSetPrimaryAddress(ctx context.Context, accountID, addressID int64) error

	// --- Wishlist ---
	SaveWishlistItem(ctx context.Context, item model.Wishlist) error
	FindWishlistByAccountID(ctx context.Context, accountID int64) ([]WishlistItemDetail, error)
	DeleteWishlistItem(ctx context.Context, accountID, productID int64) error
	CheckWishlistItemExists(ctx context.Context, accountID, productID int64) (bool, error)
	IsUsernameUsed(ctx context.Context, username string) (bool, error)
}
