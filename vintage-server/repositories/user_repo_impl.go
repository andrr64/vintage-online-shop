package repositories

import (
	"vintage-server/config"
	"vintage-server/models"
)

type userRepo struct{}

func NewUserRepository() UserRepository {
	return &userRepo{}
}

func (r *userRepo) FindByEmail(email string) (*models.User, error) {
	var user models.User
	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) FindByUsernameOrEmail(username, email string) (*models.User, error) {
	var user models.User
	if err := config.DB.Where("username = ? OR email = ?", username, email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) Create(user *models.User) error {
	return config.DB.Create(user).Error
}

func (r *userRepo) FindByID(id uint) (*models.User, error) {
	var user models.User

	if err := config.DB.First(&user, id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
