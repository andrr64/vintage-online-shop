package account

// File: internal/service/user/domain.go

import (
	"context"
	"mime/multipart"
	"vintage-server/internal/model" // Sesuaikan dengan path proyekmu

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// AccountHandler adalah kontrak atau interface yang mendefinisikan semua
// fungsionalitas yang harus disediakan oleh handler akun.
type AccountHandler interface {
	// Authentication & Registration
	RegisterCustomer(c *gin.Context)
	LoginCustomer(c *gin.Context)
	LoginAdmin(c *gin.Context)
	LoginSeller(c *gin.Context)
	Logout(c *gin.Context)

	// Profile Management
	UpdateProfile(c *gin.Context)
	UpdateAvatar(c *gin.Context)

	// Address Management
	CreateAddress(c *gin.Context)
	GetAddresses(c *gin.Context)
	UpdateAddress(c *gin.Context)
	DeleteAddress(c *gin.Context)
	SetPrimaryAddress(c *gin.Context)
}

// =================================================================================
// KONTRAK UNTUK SERVICE (Logika Bisnis) 🧠
// =================================================================================
type AccountService interface {
	// --- User & Authentication ---
	// Usecase: CustomerRegister
	RegisterCustomer(ctx context.Context, req RegisterRequest) (UserProfileResponse, error)

	UpdateProfile(ctx context.Context, userid uuid.UUID, req UpdateProfileRequest) (UserProfileResponse, error)
	UpdateAvatar(ctx context.Context, userID uuid.UUID, file multipart.File) (UserProfileResponse, error)

	LoginAdmin(ctx context.Context, req LoginRequest) (LoginResponse, error)
	LoginSeller(ctx context.Context, req LoginRequest) (LoginResponse, error)
	LoginCustomer(ctx context.Context, req LoginRequest) (LoginResponse, error)

	Logout(ctx context.Context, userId uuid.UUID) (string, error)

	// Usecase: AdminManage Users
	DeactivateUser(ctx context.Context, userID int64, reason string) error
	GetUserProfile(ctx context.Context, userID int64) (model.Account, error)

	// --- Address Management ---
	// Usecase: CustomerManage Addresses
	AddAddress(ctx context.Context, userID uuid.UUID, req AddAddressRequest) (UserAddress, error)
	GetAddressesByUserID(ctx context.Context, userID uuid.UUID) ([]UserAddress, error)
	UpdateAddress(ctx context.Context, userID uuid.UUID, address UserAddress) (UserAddress, error)
	DeleteAddress(ctx context.Context, userID uuid.UUID, addressID int64) error
	GetAddressByID(ctx context.Context, userID uuid.UUID, addressId int64) (UserAddress, error)

	// Usecase: CustomerSet Primary Address
	SetPrimaryAddress(ctx context.Context, userID uuid.UUID, addressID int64) error

	// --- Wishlist Management ---
	// Usecase: CustomerAdd/View/Remove Wishlist
	AddToWishlist(ctx context.Context, userID, productID int64) error
	GetWishlistByUserID(ctx context.Context, userID int64) ([]WishlistItemDetail, error)
	RemoveFromWishlist(ctx context.Context, userID int64, productID int64) error
}

// =================================================================================
// KONTRAK UNTUK REPOSITORY (Akses Database) 🚚
// =================================================================================
// AccountRepository mendefinisikan semua interaksi ke database yang dibutuhkan oleh Service.
type AccountRepository interface {
	WithTx(tx *sqlx.Tx) AccountRepository

	// --- Account ---
	FindAccountByID(ctx context.Context, id uuid.UUID) (model.Account, error)
	FindAccountByEmailWithRole(ctx context.Context, email string, roleName string) (model.Account, error)
	FindAccountByUsernameWithRole(ctx context.Context, username string, roleName string) (model.Account, error)

	FindAccountByUsername(ctx context.Context, username string) (model.Account, error)
	FindAccountByEmail(ctx context.Context, email string) (model.Account, error)

	InsertAccount(ctx context.Context, account model.Account) (model.Account, error)
	GetRoleIDByName(ctx context.Context, roleName string) (int64, error)
	InsertAccountRole(ctx context.Context, accountID uuid.UUID, roleID int64) error

	UpdateAccount(ctx context.Context, account model.Account) (model.Account, error)
	UpdateAvatar(ctx context.Context, avatarUrl string, id uuid.UUID) (model.Account, error)

	// --- Address ---
	SaveAddress(ctx context.Context, address model.Address) (model.Address, error)
	FindAddressesByAccountID(ctx context.Context, accountID uuid.UUID) ([]model.Address, error)
	FindAddressByIDAndAccountID(ctx context.Context, addressID int64, accountID uuid.UUID) (model.Address, error)
	UpdateAddress(ctx context.Context, accound_id uuid.UUID, address model.Address) (model.Address, error)
	DeleteAddress(ctx context.Context, addressID int64, accountID uuid.UUID) error

	SetPrimaryAddress(ctx context.Context, accountID uuid.UUID, addressID int64) error

	// --- Wishlist ---
	SaveWishlistItem(ctx context.Context, item model.Wishlist) error
	FindWishlistByAccountID(ctx context.Context, accountID int64) ([]WishlistItemDetail, error)
	DeleteWishlistItem(ctx context.Context, accountID, productID int64) error
	CheckWishlistItemExists(ctx context.Context, accountID, productID int64) (bool, error)
	IsUsernameUsed(ctx context.Context, username string) (bool, error)
}
