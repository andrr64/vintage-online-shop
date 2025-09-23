package user

import (
	"vintage-server/repositories"
	"vintage-server/services"
)

var userService = services.NewUserService(
	repositories.NewUserRepository(),
)
