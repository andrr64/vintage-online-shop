package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"time"
	"vintage-server/internal/domain/account"
	"vintage-server/pkg/apperror"

	"github.com/google/uuid"
)

// UpdateProfile pakai repo biasa karena hanya update satu row
func (s *accountService) UpdateProfile(ctx context.Context, accountID uuid.UUID, req account.UpdateProfileRequest) (account.UserProfileResponse, error) {
	acc, err := s.store.FindAccountByID(ctx, accountID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return account.UserProfileResponse{}, apperror.New(apperror.ErrCodeNotFound, "account not found")
		}
		log.Printf("Error finding account by ID: %v", err)
		return account.UserProfileResponse{}, apperror.New(apperror.ErrCodeInternal, "an internal error occurred")
	}

	acc.Firstname = req.Firstname
	acc.Lastname = req.Lastname
	acc.UpdatedAt = time.Now()

	if _, err := s.store.UpdateAccount(ctx, acc); err != nil {
		log.Printf("Error updating account: %v", err)
		return account.UserProfileResponse{}, apperror.New(apperror.ErrCodeInternal, "could not update profile")
	}

	return account.ConvertAccountToUserProfileResponse(&acc), nil
}

// UpdateAvatar pakai repo biasa tapi tetap rollback file jika error
func (s *accountService) UpdateAvatar(ctx context.Context, accountID uuid.UUID, file multipart.File) (account.UserProfileResponse, error) {
	acc, err := s.store.FindAccountByID(ctx, accountID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return account.UserProfileResponse{}, apperror.New(apperror.ErrCodeNotFound, "account not found")
		}
		log.Printf("Error finding account for avatar update: %v", err)
		return account.UserProfileResponse{}, apperror.New(apperror.ErrCodeInternal, "an internal error occurred")
	}

	oldURL := acc.AvatarURL
	filename := fmt.Sprintf("%s-avatar", uuid.NewString())

	newURL, err := s.uploader.Upload(ctx, file, filename)
	if err != nil {
		log.Printf("Error uploading avatar: %v", err)
		return account.UserProfileResponse{}, apperror.New(apperror.ErrCodeInternal, "could not process avatar update")
	}

	acc.AvatarURL = &newURL
	acc.UpdatedAt = time.Now()

	if _, err := s.store.UpdateAccount(ctx, acc); err != nil {
		log.Printf("Error updating account with new avatar URL: %v", err)
		_ = s.uploader.DeleteByURL(ctx, newURL) // rollback upload
		return account.UserProfileResponse{}, apperror.New(apperror.ErrCodeInternal, "could not save avatar information")
	}

	if oldURL != nil {
		_ = s.uploader.DeleteByURL(ctx, *oldURL) // best effort hapus avatar lama
	}

	return account.ConvertAccountToUserProfileResponse(&acc), nil
}

