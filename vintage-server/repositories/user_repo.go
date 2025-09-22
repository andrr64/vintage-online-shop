package repositories

import (
	"vintage-server/models"
)

type UserRepository interface {
	FindByUsernameOrEmail(username, email string) (*models.User, error)
	FindByEmail(email string) (*models.User, error)

	IsUsernameExists(username string) (bool, error)
	IsEmailExists(username string) (bool, error)

	FindByID(id uint) (*models.User, error)
	Create(user *models.User) error
	Update(user *models.User) error
}
