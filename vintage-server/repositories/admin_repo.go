package repositories

import (
	"vintage-server/models"
)

type AdminRepository interface {
	CreateAccount(username string, hPassword string) (*models.Admin, error)
}
