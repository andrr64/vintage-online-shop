package admin

import (
	"net/http"
	"os"
	"vintage-server/dto"
	"vintage-server/helpers/response_helper"

	"github.com/gin-gonic/gin"


)


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
	if err := c.ShouldBindJSON(&input); err != nil {
		response_helper.Failed[any](c, http.StatusUnauthorized, "Username or Password is Wrong", nil)
		return
	}
	// TESTING
	result, err := authService.Login(input)

	if err != nil {
		response_helper.Failed[any](c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	c.SetCookie(
		"access_token",
		result.Token,
		3600*24,
		"/",
		"",
		false,
		true,
	)
	response_helper.Success(c, &result, "OK")
}

// Register godoc
// @Summary Admin register *
// @Description Register new admin account
// @Tags Admin/Auth
// @Accept json
// @Produce json
// @Param input body dto.InputAdminRegisterDTO true "Register data"
// @Param key path string true "Registration Key"
// @Success 200 {object} dto.CommonResponse[dto.ResponseAdminLogin]
// @Failure 400 {object} dto.CommonResponse[string]
// @Router /admin/auth/register/{key} [post]
func Register(c *gin.Context) {
	key := c.Param("key") // ambil value {key} dari path
	var input dto.InputAdminRegisterDTO

	if key != os.Getenv("DEV_KEY") {
		response_helper.Failed[any](c, http.StatusUnauthorized, "Unauthorizedd", nil)
		return
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		response_helper.Failed[any](c, http.StatusBadRequest, "Invalid Input", nil)
		return
	}
	result, err := authService.Register(input)
	if err != nil {
		response_helper.Failed[any](c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response_helper.Success(c, &result, "OK")
}

// Logout godoc
// @Summary Admin Logout
// @Description Clear JWT cookie and logout
// @Tags Admin/Auth
// @Accept json
// @Produce json
// @Success 200 {object} dto.CommonResponse[string]
// @Router /admin/auth/logout [post]
func Logout(c *gin.Context) {
	c.SetCookie(
		"access_token",
		"",
		-1,
		"/",
		"",
		true,
		true,
	)
	response_helper.Success[string](c, nil, "logged out successfully")
}
