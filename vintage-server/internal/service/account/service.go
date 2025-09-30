package account

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"strings"
	"time"
	"vintage-server/internal/model"
	"vintage-server/pkg/apperror"
	"vintage-server/pkg/auth"
	"vintage-server/pkg/hash"
	"vintage-server/pkg/uploader"

	"github.com/google/uuid"
)

// service is a struct that will implement the Service interface from domain.go
type service struct {
	repo     Repository
	jwt      *auth.JWTService
	uploader uploader.Uploader
}

// NewService is the constructor for the service
func NewService(repo Repository, jwtSecret string, uploader uploader.Uploader) Service {
	return &service{
		repo:     repo,
		jwt:      auth.NewJWTService(jwtSecret),
		uploader: uploader,
	}
}

// RegisterCustomer handles new user registration.
func (s *service) RegisterCustomer(ctx context.Context, req RegisterRequest) (UserProfileResponse, error) {
	_, err := s.repo.FindAccountByUsername(ctx, req.Username)
	if err != sql.ErrNoRows {
		if err == nil {
			return UserProfileResponse{}, apperror.New(apperror.ErrCodeConflict, "username or email already taken")
		}
		log.Printf("Error finding account by username: %v", err)
		return UserProfileResponse{}, apperror.New(apperror.ErrCodeInternal, "an internal error occurred")
	}

	_, err = s.repo.FindAccountByEmail(ctx, req.Email)
	if err != sql.ErrNoRows {
		if err == nil {
			return UserProfileResponse{}, apperror.New(apperror.ErrCodeConflict, "username or email already taken")
		}
		log.Printf("Error finding account by email: %v", err)
		return UserProfileResponse{}, apperror.New(apperror.ErrCodeInternal, "an internal error occurred")
	}

	hashedPassword, err := hash.Generate(req.Password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return UserProfileResponse{}, apperror.New(apperror.ErrCodeInternal, "failed to process registration")
	}

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

	createdAcc, err := s.repo.SaveAccount(ctx, newAccount, "admin")
	if err != nil {
		log.Printf("Error saving account: %v", err)
		return UserProfileResponse{}, apperror.New(apperror.ErrCodeInternal, "failed to create account")
	}

	return ConvertAccountToUserProfileResponse(&createdAcc), nil
}

// LoginCustomer handles user authentication.
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
		if errors.Is(err, sql.ErrNoRows) {
			return LoginResponse{}, apperror.New(apperror.ErrCodeUnauthorized, "invalid credentials")
		}
		log.Printf("Error finding account for login: %v", err)
		return LoginResponse{}, apperror.New(apperror.ErrCodeInternal, "an internal error occurred")
	}

	if err := hash.Verify(acc.Password, req.Password); err != nil {
		return LoginResponse{}, apperror.New(apperror.ErrCodeUnauthorized, "invalid credentials")
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

// Logout is a placeholder for user logout logic.
func (s *service) Logout(ctx context.Context, userId uuid.UUID) (string, error) {
	// Future implementation: token invalidation, logging, etc.
	return "OK", nil
}

// UpdateProfile updates a user's first and last name.
func (s *service) UpdateProfile(ctx context.Context, account_id uuid.UUID, req UpdateProfileRequest) (UserProfileResponse, error) {
	account, err := s.repo.FindAccountByID(ctx, account_id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return UserProfileResponse{}, apperror.New(apperror.ErrCodeNotFound, "account not found")
		}
		log.Printf("Error finding account by ID: %v", err)
		return UserProfileResponse{}, apperror.New(apperror.ErrCodeInternal, "an internal error occurred")
	}

	account.Firstname = req.Firstname
	account.Lastname = req.Lastname
	account.UpdatedAt = time.Now()

	if err := s.repo.UpdateAccount(ctx, account); err != nil {
		log.Printf("Error updating account: %v", err)
		return UserProfileResponse{}, apperror.New(apperror.ErrCodeInternal, "could not update profile")
	}

	return ConvertAccountToUserProfileResponse(&account), nil
}

// UpdateAvatar updates a user's avatar image.
func (s *service) UpdateAvatar(ctx context.Context, accountId uuid.UUID, file multipart.File) (UserProfileResponse, error) {
	acc, err := s.repo.FindAccountByID(ctx, accountId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return UserProfileResponse{}, apperror.New(apperror.ErrCodeNotFound, "account not found")
		}
		log.Printf("Error finding account for avatar update: %v", err)
		return UserProfileResponse{}, apperror.New(apperror.ErrCodeInternal, "an internal error occurred")
	}

	oldURL := acc.AvatarURL
	filename := fmt.Sprintf("%s-avatar", uuid.NewString())

	newURL, err := s.uploader.Upload(ctx, file, filename)
	if err != nil {
		log.Printf("Error uploading avatar: %v", err)
		return UserProfileResponse{}, apperror.New(apperror.ErrCodeInternal, "could not process avatar update")
	}

	acc.AvatarURL = &newURL
	acc.UpdatedAt = time.Now()

	if err := s.repo.UpdateAccount(ctx, acc); err != nil {
		log.Printf("Error updating account with new avatar URL: %v", err)
		_ = s.uploader.Delete(ctx, newURL) // Rollback: attempt to delete the newly uploaded file
		return UserProfileResponse{}, apperror.New(apperror.ErrCodeInternal, "could not save avatar information")
	}

	if oldURL != nil {
		_ = s.uploader.Delete(ctx, *oldURL) // Best effort deletion of the old avatar
	}

	return ConvertAccountToUserProfileResponse(&acc), nil
}

// AddAddress adds a new address for a user.
func (s *service) AddAddress(ctx context.Context, accountID uuid.UUID, req AddAddressRequest) (UserAddress, error) {
	address := NewAddressFromRequest(accountID, req, false)

	savedAddress, err := s.repo.SaveAddress(ctx, address)
	if err != nil {
		if strings.Contains(err.Error(), "violate") { // Note: Fragile check, better to use specific repo errors if possible
			return UserAddress{}, apperror.New(apperror.ErrCodeBadRequest, "invalid address data provided")
		}
		log.Printf("Error saving address: %v", err)
		return UserAddress{}, apperror.New(apperror.ErrCodeInternal, "could not save address")
	}

	return ConvertAddressToDTO(savedAddress), nil
}

// UpdateAddress updates an existing user address.
func (s *service) UpdateAddress(ctx context.Context, accountId uuid.UUID, req UserAddress) (UserAddress, error) {
	old, err := s.repo.FindAddressByIDAndAccountID(ctx, req.ID, accountId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return UserAddress{}, apperror.New(apperror.ErrCodeNotFound, "address not found")
		}
		log.Printf("Error finding address to update: %v", err)
		return UserAddress{}, apperror.New(apperror.ErrCodeInternal, "an internal error occurred")
	}

	old.ProvinceID = req.ProvinceID
	old.RegencyID = req.RegencyID
	old.DistrictID = req.DistrictID
	old.VillageID = req.VillageID
	old.RecipientName = req.RecipientName
	old.RecipientPhone = req.RecipientPhone
	old.Label = req.Label
	old.PostalCode = req.PostalCode
	old.Street = req.Street
	old.UpdatedAt = time.Now()

	updatedAddress, err := s.repo.UpdateAddress(ctx, accountId, old)
	if err != nil {
		log.Printf("Error updating address: %v", err)
		return UserAddress{}, apperror.New(apperror.ErrCodeInternal, "could not update address")
	}

	return ConvertAddressToDTO(updatedAddress), nil
}

// GetAddressByID retrieves a single address by its ID for a specific user.
func (s *service) GetAddressByID(ctx context.Context, accountID uuid.UUID, addressID int64) (UserAddress, error) {
	address, err := s.repo.FindAddressByIDAndAccountID(ctx, addressID, accountID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return UserAddress{}, apperror.New(apperror.ErrCodeNotFound, "address not found")
		}
		log.Printf("Error getting address by ID: %v", err)
		return UserAddress{}, apperror.New(apperror.ErrCodeInternal, "an internal error occurred")
	}
	return ConvertAddressToDTO(address), nil
}

// GetAddressesByUserID retrieves all addresses for a specific user.
func (s *service) GetAddressesByUserID(ctx context.Context, userID uuid.UUID) ([]UserAddress, error) {
	addresses, err := s.repo.FindAddressesByAccountID(ctx, userID)
	if err != nil {
		// sql.ErrNoRows is not an error here; it just means the user has no addresses yet.
		if errors.Is(err, sql.ErrNoRows) {
			return []UserAddress{}, nil // Return an empty slice instead of a not found error
		}
		log.Printf("Error finding addresses by account ID: %v", err)
		return nil, apperror.New(apperror.ErrCodeInternal, "an internal error occurred")
	}
	return ConvertAddressesToDTO(addresses), nil
}

// DeleteAddress deletes a user's address.
func (s *service) DeleteAddress(ctx context.Context, userID uuid.UUID, addressID int64) error {
	err := s.repo.DeleteAddress(ctx, addressID, userID)
	if err != nil {
		log.Printf("Error deleting address: %v", err)
		return apperror.New(apperror.ErrCodeInternal, "could not delete address")
	}
	return nil
}

// SetPrimaryAddress sets a specific address as the user's primary one.
func (s *service) SetPrimaryAddress(ctx context.Context, accountID uuid.UUID, addressID int64) error {
	err := s.repo.SetPrimaryAddress(ctx, accountID, addressID)
	if err != nil {
		// CRITICAL: Do not leak the raw error message to the client.
		log.Printf("Error setting primary address: %v", err)
		return apperror.New(apperror.ErrCodeInternal, "could not set primary address")
	}
	return nil
}

// Login as Seller
func (s *service) LoginSeller(ctx context.Context, req LoginRequest) (LoginResponse, error) {
	var acc model.Account
	var err error

	role := "seller"
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

// Login as Admin
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

// -------- TO-DO ---------
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
