package impl

import (
	"vintage-server/models"
	repo "vintage-server/repositories"
)

type adminRepo struct{}

func NewAdminRepository() repo.AdminRepository {
	return &adminRepo{}
}

func (r *adminRepo) CreateAccount(username string, hPassword string) (*models.Admin, error) {
	return &models.Admin{}, nil
}
