package admin

import (
	"net/http"
	"vintage-server/dto"
	"vintage-server/helpers/response_helper"

	"github.com/gin-gonic/gin"

	jwt "vintage-server/security/jwt"

	repo "vintage-server/repositories/impl"
	adminService "vintage-server/services/admin/impl"
)

var authService = adminService.NewAdminAuthService(
	repo.NewAdminRepository(),
)

// Login godoc
// @Summary Admin login
// @Description Login for admin users
// @Tags Admin/Auth
// @Accept json
// @Produce json
// @Param input body dto.InputAdminLoginDTO true "Login credentials"
// @Success 200 {object} dto.CommonResponse[dto.ResponseAdminLogin]
// @Failure 400 {object} dto.CommonResponse[string]
// @Router /admin/auth/login [post]
func Login(c *gin.Context) {
	var input dto.InputAdminLoginDTO
	if err := c.ShouldBind(&input); err != nil {
		response_helper.Failed[any](c, http.StatusUnauthorized, "Username or Password is Wrong", nil)
		return
	}
	// TESTING
	response := dto.ResponseAdminLogin{
		Username: "test",
		Token:    jwt.CreateToken(1, "test"),
	}
	response_helper.Success(c, &response, "OK")
}

// Register godoc
// @Summary Admin register
// @Description Register new admin account
// @Tags Admin/Auth
// @Accept json
// @Produce json
// @Param input body dto.InputAdminRegisterDTO true "Register data"
// @Success 200 {object} dto.CommonResponse[dto.ResponseAdminLogin]
// @Failure 400 {object} dto.CommonResponse[string]
// @Router /admin/auth/register [post]
func Register(c *gin.Context) {
	var input dto.InputAdminRegisterDTO

	if err := c.ShouldBind(&input); err != nil {
		response_helper.Failed[any](c, http.StatusBadRequest, "Invalid Input", nil)
		return
	}
	result, error := authService.Register(input)
	if error != nil {
		response_helper.Failed[any](c, http.StatusBadRequest, error.Error(), nil)
		return
	}
	response_helper.Success(c, &result, "OK")
}
