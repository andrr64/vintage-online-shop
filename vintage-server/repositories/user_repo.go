package repositories

import (
	"vintage-server/models"
)

type UserRepository interface {
	FindByUsernameOrEmail(username, email string) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	Create(user *models.User) error
}
