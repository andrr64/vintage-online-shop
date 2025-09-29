package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"time"
	"vintage-server/internal/model"
	"vintage-server/pkg/apperror"
	"vintage-server/pkg/auth"
	"vintage-server/pkg/hash"
	"vintage-server/pkg/uploader"

	"strings"

	"github.com/google/uuid"
)

// service adalah struct yang akan mengimplementasikan interface Service dari domain.go
type service struct {
	repo     Repository
	jwt      *auth.JWTService
	uploader uploader.Uploader // <-- Tambahkan dependency ke uploader
}

// NewService adalah constructor untuk service
func NewService(repo Repository, jwtSecret string, uploader uploader.Uploader) Service {
	return &service{
		repo:     repo,
		jwt:      auth.NewJWTService(jwtSecret),
		uploader: uploader,
	}
}

// ----- DONE ---------------
func (s *service) RegisterCustomer(ctx context.Context, req RegisterRequest) (UserProfileResponse, error) {
	// 1. Validasi duplikasi data dalam satu transaksi untuk konsistensi
	_, err := s.repo.FindAccountByUsername(ctx, req.Username)
	if err != sql.ErrNoRows { // Jika BUKAN error "tidak ketemu"
		if err == nil { // Jika tidak ada error sama sekali, artinya username ditemukan
			// KEMBALIKAN PESAN ERROR AMBIGU
			return UserProfileResponse{}, apperror.New(apperror.ErrCodeConflict, "Username or Email already taken")
		}
		// Untuk error database lainnya
		log.Printf("Error finding account by username: %v", err)
		return UserProfileResponse{}, apperror.New(apperror.ErrCodeInternal, "an internal error occurred")
	}

	_, err = s.repo.FindAccountByEmail(ctx, req.Email)
	if err != sql.ErrNoRows { // Jika BUKAN error "tidak ketemu"
		if err == nil { // Jika tidak ada error sama sekali, artinya email ditemukan
			// KEMBALIKAN PESAN ERROR AMBIGU YANG SAMA
			return UserProfileResponse{}, apperror.New(apperror.ErrCodeConflict, "Username or Email already taken")
		}
		// Untuk error database lainnya
		log.Printf("Error finding account by email: %v", err)
		return UserProfileResponse{}, apperror.New(apperror.ErrCodeInternal, "an internal error occurred")
	}

	// 2. Hash password
	hashedPassword, err := hash.Generate(req.Password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return UserProfileResponse{}, apperror.New(apperror.ErrCodeInternal, "failed to process registration")
	}

	// 3. Buat entitas akun baru
	newAccount := model.Account{
		Username:  req.Username,
		Email:     req.Email,
		Password:  hashedPassword,
		Firstname: req.Firstname,
		Lastname:  req.Lastname,
		Active:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 4. Simpan ke repository
	createdAcc, err := s.repo.SaveAccount(ctx, newAccount, "customer")
	if err != nil {
		log.Printf("Error saving account: %v", err)
		return UserProfileResponse{}, apperror.New(apperror.ErrCodeInternal, "failed to create account")
	}

	// 5. Kembalikan response sukses menggunakan DTO
	response := UserProfileResponse{
		Username:  createdAcc.Username,
		Firstname: createdAcc.Firstname,
		Lastname:  createdAcc.Lastname,
		Email:     createdAcc.Email,
		AvatarURL: createdAcc.AvatarURL,
	}

	return response, nil
}

func (s *service) LoginCustomer(ctx context.Context, req LoginRequest) (LoginResponse, error) {
	var acc model.Account
	var err error

	roleName := "customer"

	if strings.Contains(req.Identifier, "@") {
		acc, err = s.repo.FindAccountByEmailWithRole(ctx, req.Identifier, roleName)
	} else {
		acc, err = s.repo.FindAccountByUsernameWithRole(ctx, req.Identifier, roleName)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return LoginResponse{}, apperror.New(apperror.ErrCodeUnauthorized, "invalid data")
		}
		log.Printf("Error finding account: %v", err)
		return LoginResponse{}, apperror.New(apperror.ErrCodeInternal, "an internal error occurred")
	}

	if err := hash.Verify(acc.Password, req.Password); err != nil {
		return LoginResponse{}, apperror.New(apperror.ErrCodeUnauthorized, "invalid data")
	}

	token, err := s.jwt.GenerateToken(acc.ID, roleName)
	if err != nil {
		log.Printf("Error generating token: %v", err)
		return LoginResponse{}, apperror.New(apperror.ErrCodeInternal, "an internal error occurred")
	}

	return LoginResponse{
		AccessToken: token,
		UserProfile: ConvertAccountToUserProfileResponse(&acc),
	}, nil
}

func (s *service) Logout(ctx context.Context, userId uuid.UUID) (string, error) {
	// implementation logging, dst
	return "OK", nil
}

func (s *service) UpdateProfile(ctx context.Context, account_id uuid.UUID, req UpdateProfileRequest) (UserProfileResponse, error) {
	account, err := s.repo.FindAccountByID(ctx, account_id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return UserProfileResponse{}, apperror.New(
				apperror.ErrCodeUnauthorized,
				"Account Not Found",
			)
		}
	}

	account.Firstname = req.Firstname
	account.Lastname = req.Lastname
	account.UpdatedAt = time.Now()

	err = s.repo.UpdateAccount(ctx, account)

	if err != nil {
		return UserProfileResponse{}, apperror.New(
			apperror.ErrCodeInternal,
			"Something wrong happened when updating your data",
		)
	}

	return ConvertAccountToUserProfileResponse(&account), nil
}

func (s *service) UpdateAvatar(
	ctx context.Context,
	accountId uuid.UUID,
	file multipart.File,
) (UserProfileResponse, error) {
	// 1. Cari account
	acc, err := s.repo.FindAccountByID(ctx, accountId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return UserProfileResponse{}, apperror.New(404, "Account not found")
		}
		return UserProfileResponse{}, apperror.New(500, "Failed to fetch account data")
	}

	oldURL := acc.AvatarURL

	// 2. Upload ke Cloudinary
	filename := fmt.Sprintf("%s-avatar", uuid.NewString())
	newURL, err := s.uploader.Upload(ctx, file, filename)
	if err != nil {
		return UserProfileResponse{}, apperror.New(500, "Failed to upload avatar")
	}

	// 3. Update account dengan avatar baru
	acc.AvatarURL = &newURL
	acc.UpdatedAt = time.Now()

	if err := s.repo.UpdateAccount(ctx, acc); err != nil {
		// rollback → hapus file baru yang sudah diupload
		_ = s.uploader.Delete(ctx, newURL)
		return UserProfileResponse{}, apperror.New(500, "Failed to update avatar in database")
	}

	// 4. Jika ada avatar lama → hapus dari Cloudinary (best effort, error diabaikan)
	if oldURL != nil {
		_ = s.uploader.Delete(ctx, *oldURL)
	}

	// 5. Return profile baru
	return ConvertAccountToUserProfileResponse(&acc), nil
}

func (s *service) AddAddress(ctx context.Context, accountID uuid.UUID, req AddAddressRequest) (UserAddress, error) {
	address := NewAddressFromRequest(accountID, req, false)

	address, err := s.repo.SaveAddress(ctx, address)

	if err != nil {
		if strings.Contains(err.Error(), "violate") {
			return UserAddress{}, apperror.New(http.StatusBadRequest, "Bad request, check your data!")
		}
		return UserAddress{}, apperror.New(apperror.ErrCodeInternal, "Something wrong when adding your data")
	}

	address, err = s.repo.SaveAddress(ctx, address)

	if err != nil {
		return UserAddress{}, err
	}
	return ConvertAddressToDTO(address), nil
}

func (s *service) UpdateAddress(ctx context.Context, accountId uuid.UUID, addressID int64, req UserAddress) (UserAddress, error) {
	old, err := s.repo.FindAddressByIDAndAccountID(ctx, addressID, accountId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return UserAddress{}, apperror.New(apperror.ErrCodeNotFound, "Data not found!")
		}
		return UserAddress{}, apperror.New(apperror.ErrCodeInternal, "Something wrong")
	}

	old.ProvinceID = req.ProvinceID
	old.RegencyID = req.RegencyID
	old.DistrictID = req.DistrictID
	old.VillageID = req.VillageID
	old.RecipientName = req.RecipientName
	old.RecipientPhone = req.RecipientPhone
	old.Label = req.Label
	old.RecipientPhone = req.RecipientPhone
	old.PostalCode = req.PostalCode
	old.Street = req.Street

	_, err = s.repo.UpdateAddress(ctx, accountId, old)

	if err != nil {
		return UserAddress{}, apperror.New(apperror.ErrCodeInternal, "Something wrong")
	}

	// TODO: Implement business logic
	return ConvertAddressToDTO(old), nil
}

func (s *service) GetAddressByID(ctx context.Context, accountID uuid.UUID, addressID int64) (UserAddress, error) {
	address, err := s.repo.FindAddressByIDAndAccountID(ctx, addressID, accountID)
	if err != nil {
		return UserAddress{}, apperror.New(http.StatusNotFound, "Not found")
	}
	return ConvertAddressToDTO(address), nil
}

func (s *service) GetAddressesByUserID(ctx context.Context, userID uuid.UUID) ([]UserAddress, error) {
	addresses, err := s.repo.FindAddressesByAccountID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.New(apperror.ErrCodeNotFound, "Not found")
		}
		return nil, err
	}
	return ConvertAddressesToDTO(addresses), nil
}

func (s *service) DeleteAddress(ctx context.Context, userID uuid.UUID, addressID int64) error {
	err := s.repo.DeleteAddress(ctx, addressID, userID)
	if err != nil {
		return apperror.New(apperror.ErrCodeInternal, "Something wong...")
	}
	return nil
}

func (s *service) SetPrimaryAddress(ctx context.Context, accountID uuid.UUID, addressID int64) error {
	err := s.repo.SetPrimaryAddress(ctx, accountID, addressID)
	if err != nil {
		return apperror.New(apperror.ErrCodeInternal, "Something wrong")
	}
	return nil
}

// ---- TO-DO ------
func (s *service) LoginAdmin(ctx context.Context, req LoginRequest) (LoginResponse, error) {
	var acc model.Account
	var err error

	role := "admin"
	if strings.Contains(req.Identifier, "@") {
		acc, err = s.repo.FindAccountByEmailWithRole(ctx, req.Identifier, role)
	} else {
		acc, err = s.repo.FindAccountByUsernameWithRole(ctx, req.Identifier, role)
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return LoginResponse{}, apperror.New(
				apperror.ErrCodeUnauthorized,
				"Invalid Data",
			)
		}
		log.Printf("Error finding account: %v", err) // Contoh logging
		return LoginResponse{}, apperror.New(apperror.ErrCodeInternal, "Invalid Data")
	}

	err = hash.Verify(acc.Password, req.Password)
	if err != nil {
		return LoginResponse{}, apperror.New(
			apperror.ErrCodeUnauthorized,
			"Invalid Data",
		)
	}
	access_token, err := s.jwt.GenerateToken(acc.ID, role)
	if err != nil {
		return LoginResponse{}, &apperror.AppError{
			Code:    apperror.ErrCodeInternal,
			Message: err.Error(),
		}
	}
	return LoginResponse{
		AccessToken: access_token,
		UserProfile: ConvertAccountToUserProfileResponse(&acc),
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
