package user

import (
	"net/http"
	"vintage-server/dto"
	"vintage-server/helpers/response_helper"

	"github.com/gin-gonic/gin"
)

// Register godoc
// @Summary Register a new user
// @Description Create a new user account with email, username, fullname, and password
// @Tags User/Account Management
// @Accept json
// @Produce json
// @Param input body dto.InputRegisterDTO true "Register user data"
// @Success 200 {object} dto.CommonResponse[dto.ResponseRegisterUserDTO]
// @Failure 400 {object} dto.CommonResponse[string]
// @Router /user/auth/register [post]
func Register(c *gin.Context) {
	var input dto.InputRegisterDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		response_helper.Failed[string](c, http.StatusBadRequest, "Invalid input", nil)
		return
	}

	userDTO, err := userService.Register(input)
	if err != nil {
		response_helper.Failed[string](c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response_helper.Success(c, &userDTO, "User registered successfully")
}

// Login godoc
// @Summary User login
// @Description Login with email and password. Returns JWT token and sets httpOnly cookie
// @Tags User/Account Management
// @Accept json
// @Produce json
// @Param input body dto.InputLoginDTO true "Login credentials"
// @Success 200 {object} dto.CommonResponse[map[string]string]
// @Failure 400 {object} dto.CommonResponse[string]
// @Failure 401 {object} dto.CommonResponse[string]
// @Router /user/auth/login [post]
func Login(c *gin.Context) {
	var input dto.InputLoginDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		response_helper.Failed[string](c, http.StatusBadRequest, "Invalid input", nil)
		return
	}

	token, err := userService.Login(input)
	if err != nil {
		response_helper.Failed[string](c, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	// Set cookie JWT
	c.SetCookie(
		"access_token",
		token,
		3600*24,
		"/",
		"",
		false,
		true,
	)

	response_helper.Success(c, &map[string]string{"token": token}, "Login successful")
}

// Logout godoc
// @Summary Logout user
// @Description Clear JWT cookie and logout user
// @Tags User/Account Management
// @Accept json
// @Produce json
// @Success 200 {object} dto.CommonResponse[string]
// @Router /user/auth/logout [post]
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
	response_helper.Success[string](c, nil, "User logged out successfully")
}
