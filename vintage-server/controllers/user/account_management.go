package user

import (
	"net/http"
	"vintage-server/dto"
	"vintage-server/helpers/response_helper"
	"vintage-server/repositories"
	"vintage-server/services"

	"github.com/gin-gonic/gin"
)

var userService = services.NewUserService(
	repositories.NewUserRepository(),
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
// @Router /api/v1/users/register [post]
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
// @Router /api/v1/users/login [post]
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
// @Router /api/v1/users/logout [post]
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

// GetAccount godoc
// @Summary Get current user account info
// @Description Retrieve the current logged-in user's username, email, and fullname
// @Tags User/Account Management
// @Accept json
// @Produce json
// @Success 200 {object} dto.CommonResponse[dto.ResponseUserInfoDTO]
// @Failure 401 {object} dto.CommonResponse[string]
// @Failure 404 {object} dto.CommonResponse[string]
// @Security ApiKeyAuth
// @Router /api/v1/users/account [get]
func GetAccount(c *gin.Context) {
	currentUserI, exists := c.Get("currentUser")
	if !exists {
		response_helper.Failed[string](c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	currentUser := currentUserI.(map[string]any)
	userID := currentUser["id"].(uint)

	user, err := userService.FindByID(userID)
	if err != nil {
		response_helper.Failed[string](c, http.StatusNotFound, "User not found", nil)
		return
	}

	response := dto.ResponseUserInfoDTO{
		Username: user.Username,
		Email:    user.Email,
		Fullname: user.Fullname,
	}

	response_helper.Success(c, &response, "OK")
}

// UpdatePassword godoc
// @Summary Update current user password
// @Description Update password by providing old password and new password
// @Tags User/Account Management
// @Accept json
// @Produce json
// @Param input body dto.InputUpdatePasswordDTO true "Old and new password"
// @Success 200 {object} dto.CommonResponse[bool]
// @Failure 400 {object} dto.CommonResponse[string]
// @Failure 401 {object} dto.CommonResponse[string]
// @Failure 404 {object} dto.CommonResponse[string]
// @Security ApiKeyAuth
// @Router /api/v1/users/account/password [put]
func UpdatePassword(c *gin.Context) {
	currentUserI, exists := c.Get("currentUser")
	if !exists {
		response_helper.Failed[string](c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	var input dto.InputUpdatePasswordDTO

	if err := c.ShouldBindJSON(&input); err != nil {
		response_helper.Failed[string](c, http.StatusBadRequest, "Invalid input", nil)
		return
	}

	currentUser := currentUserI.(map[string]any)
	userID := currentUser["id"].(uint)

	http_code, err := userService.UpdatePassword(userID, input)
	if err != nil {
		response_helper.Failed[any](c, http_code, err.Error(), nil)
		return
	}
	response_helper.Success[any](c, nil, "OK")
}

// UpdateProfile godoc
// @Summary Update current user profile
// @Description Update username, email, and fullname of the current user
// @Tags User/Account Management
// @Accept json
// @Produce json
// @Param input body dto.InputUpdateAccountDTO true "Update account data"
// @Success 200 {object} dto.CommonResponse[dto.ResponseUserInfoDTO]
// @Failure 400 {object} dto.CommonResponse[string]
// @Failure 401 {object} dto.CommonResponse[string]
// @Failure 409 {object} dto.CommonResponse[string] "username/email already exists"
// @Security ApiKeyAuth
// @Router /api/v1/users/account/profile [put]
func UpdateProfile(c *gin.Context) {
	currentUserI, exists := c.Get("currentUser")
	if !exists {
		response_helper.Failed[string](c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	var input dto.InputUpdateAccountDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		response_helper.Failed[string](c, http.StatusBadRequest, "Invalid input", nil)
		return
	}
	currentUser := currentUserI.(map[string]any)
	userID := currentUser["id"].(uint)

	updatedUser, err := userService.UpdateAccount(userID, input)

	if err != nil {
		response_helper.Failed[any](c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response_helper.Success(c, &updatedUser, "OK")
}
