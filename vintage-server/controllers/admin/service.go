package admin

import (
	repo "vintage-server/repositories/impl"
	adminService "vintage-server/services/admin/impl"
)

var authService = adminService.NewAdminAuthService(
	repo.NewAdminRepository(),
)
