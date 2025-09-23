package repositories

import (
	"vintage-server/models"
)

type AdminRepository interface {
	IsExists(username string) (bool, error)
	CreateAccount(username string, hPassword string) (*models.Admin, error)
}
