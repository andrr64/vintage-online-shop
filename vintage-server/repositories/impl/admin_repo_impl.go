package impl

import (
	"errors"
	"vintage-server/config"
	"vintage-server/models"
	repo "vintage-server/repositories"
)

type adminRepo struct{}

// NewAdminRepository buat instance adminRepo
func NewAdminRepository() repo.AdminRepository {
	return &adminRepo{}
}

// CreateAccount menyimpan akun admin baru ke database
func (r *adminRepo) CreateAccount(username string, hPassword string) (*models.Admin, error) {
	admin := &models.Admin{
		Username: username,
		Password: hPassword,
	}

	// Simpan ke DB
	if err := config.DB.Create(admin).Error; err != nil {
		return nil, err
	}

	return admin, nil
}

// IsExists mengecek apakah username sudah ada di database
func (r *adminRepo) IsExists(username string) (bool, error) {
	var admin models.Admin
	err := config.DB.Where("username = ?", username).First(&admin).Error
	if err != nil {
		if errors.Is(err, config.DB.Error) { // atau ganti dengan gorm.ErrRecordNotFound
			return false, nil
		}
		return false, err
	}
	return true, nil
}
