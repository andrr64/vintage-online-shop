package user

import (
	"context"
	"database/sql"
	"log"
	"time"
	"vintage-server/internal/model"
	"vintage-server/pkg/apperror"
	"vintage-server/pkg/auth"
	"vintage-server/pkg/hash"

	"strings"
)

// service adalah struct yang akan mengimplementasikan interface Service dari domain.go
type service struct {
	repo Repository
	jwt  *auth.JWTService
}

// NewService adalah constructor untuk service
func NewService(repo Repository, jwtSecret string) Service {
	return &service{
		repo: repo,
		jwt:  auth.NewJWTService(jwtSecret),
	}
}

// --- User & Authentication ---
// File: internal/service/user/service.go

// File: internal/service/user/service.go

func (s *service) Register(ctx context.Context, req RegisterRequest) (UserProfileResponse, error) {
	// 1. Validasi duplikasi data dalam satu transaksi untuk konsistensi
	// Kita cek keduanya, tapi pesan error yang kita kembalikan akan sama.
	_, err := s.repo.FindAccountByUsername(ctx, req.Username)
	if err != sql.ErrNoRows { // Jika BUKAN error "tidak ketemu"
		if err == nil { // Jika tidak ada error sama sekali, artinya username ditemukan
			// KEMBALIKAN PESAN ERROR AMBIGU
			return UserProfileResponse{}, apperror.New(apperror.ErrCodeConflict, "username or email already exists")
		}
		// Untuk error database lainnya
		log.Printf("Error finding account by username: %v", err)
		return UserProfileResponse{}, apperror.New(apperror.ErrCodeInternal, "an internal error occurred")
	}

	_, err = s.repo.FindAccountByEmail(ctx, req.Email)
	if err != sql.ErrNoRows { // Jika BUKAN error "tidak ketemu"
		if err == nil { // Jika tidak ada error sama sekali, artinya email ditemukan
			// KEMBALIKAN PESAN ERROR AMBIGU YANG SAMA
			return UserProfileResponse{}, apperror.New(apperror.ErrCodeConflict, "username or email already exists")
		}
		// Untuk error database lainnya
		log.Printf("Error finding account by email: %v", err)
		return UserProfileResponse{}, apperror.New(apperror.ErrCodeInternal, "an internal error occurred")
	}

	// Jika kita sampai di sini, artinya username dan email tersedia.

	// 2. Hash password
	hashedPassword, err := hash.Generate(req.Password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return UserProfileResponse{}, apperror.New(apperror.ErrCodeInternal, "failed to process registration")
	}

	// 3. Buat entitas akun baru
	newAccount := model.Account{
		Username:  req.Username,
		Email:     req.Email, // <-- PERBAIKAN: Ambil alamat memori dari string
		Password:  hashedPassword,
		Role:      model.RoleCustomer, // Role default: Customer
		Firstname: &req.Firstname,
		Lastname:  req.Lastname,
		Active:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 4. Simpan ke repository
	createdAcc, err := s.repo.SaveAccount(ctx, newAccount)
	if err != nil {
		log.Printf("Error saving account: %v", err)
		return UserProfileResponse{}, apperror.New(apperror.ErrCodeInternal, "failed to create account")
	}

	// 5. Kembalikan response sukses menggunakan DTO
	response := UserProfileResponse{
		ID:        createdAcc.ID,
		Username:  createdAcc.Username,
		Firstname: createdAcc.Firstname,
		Lastname:  createdAcc.Lastname,
		Email:     createdAcc.Email,
		AvatarURL: createdAcc.AvatarURL,
	}

	return response, nil
}

// File: internal/service/user/service.go

// Login melakukan autentikasi user dengan username/email dan password
func (s *service) Login(ctx context.Context, req LoginRequest) (LoginResponse, error) {
	var acc model.Account
	var err error

	// 1. Cari akun berdasarkan identifier (email atau username)
	if strings.Contains(req.Identifier, "@") {
		acc, err = s.repo.FindAccountByEmail(ctx, req.Identifier)
	} else {
		acc, err = s.repo.FindAccountByUsername(ctx, req.Identifier)
	}

	// 2. Handle error pencarian akun (user tidak ada atau error DB)
	// Jika user tidak ditemukan (sql.ErrNoRows), kita tetap anggap sebagai Unauthorized.
	// Ini untuk mencegah user enumeration.
	if err != nil {
		if err == sql.ErrNoRows {
			return LoginResponse{}, apperror.New(apperror.ErrCodeUnauthorized, "invalid data")
		}
		// Untuk error database lainnya, log error asli dan kembalikan error internal.
		log.Printf("Error finding account: %v", err) // Contoh logging
		return LoginResponse{}, apperror.New(apperror.ErrCodeInternal, "an internal error occurred")
	}

	// 3. Verifikasi password
	err = hash.Verify(acc.Password, req.Password)
	if err != nil {
		// Jika password salah, kembalikan error Unauthorized yang SAMA.
		return LoginResponse{}, apperror.New(apperror.ErrCodeUnauthorized, "invalid data")
	}

	// 4. Buat token JWT jika semua berhasil
	token, err := s.jwt.GenerateToken(acc.ID, acc.Role)
	if err != nil {
		log.Printf("Error generating token: %v", err) // Contoh logging
		return LoginResponse{}, apperror.New(apperror.ErrCodeInternal, "an internal error occurred")
	}

	// 5. Kembalikan response sukses dengan DTO
	return LoginResponse{
		AccessToken: token,
		UserProfile: UserProfileResponse{
			ID:        acc.ID,
			Username:  acc.Username,
			AvatarURL: acc.AvatarURL,
			Email:     acc.Email,
			Firstname: acc.Firstname,
			Lastname:  acc.Lastname,
		},
	}, nil
}

func (s *service) DeactivateUser(ctx context.Context, userID int64, reason string) error {
	// TODO: Implement business logic
	return nil
}

func (s *service) GetUserProfile(ctx context.Context, userID int64) (model.Account, error) {
	// TODO: Implement business logic
	return model.Account{}, nil
}

// --- Address Management ---

func (s *service) AddAddress(ctx context.Context, userID int64, req model.Address) (model.Address, error) {
	// TODO: Implement business logic
	return model.Address{}, nil
}

func (s *service) GetAddressesByUserID(ctx context.Context, userID int64) ([]model.Address, error) {
	// TODO: Implement business logic
	return nil, nil
}

func (s *service) UpdateAddress(ctx context.Context, userID, addressID int64, req model.Address) (model.Address, error) {
	// TODO: Implement business logic
	return model.Address{}, nil
}

func (s *service) DeleteAddress(ctx context.Context, userID, addressID int64) error {
	// TODO: Implement business logic
	return nil
}

func (s *service) SetPrimaryAddress(ctx context.Context, userID, addressID int64) error {
	// TODO: Implement business logic
	return nil
}

// --- Wishlist Management ---

func (s *service) AddToWishlist(ctx context.Context, userID, productID int64) error {
	// TODO: Implement business logic
	return nil
}

func (s *service) GetWishlistByUserID(ctx context.Context, userID int64) ([]WishlistItemDetail, error) {
	// TODO: Implement business logic
	return nil, nil
}

func (s *service) RemoveFromWishlist(ctx context.Context, userID, productID int64) error {
	// TODO: Implement business logic
	return nil
}
