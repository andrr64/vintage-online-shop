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

func Register(c *gin.Context) {
	var input dto.RegisterUserDTO
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

func Login(c *gin.Context) {
	var input dto.LoginUserDTO
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

func UpdateProfile(c *gin.Context) {
	currentUserI, exists := c.Get("currentUser")
	if !exists {
		response_helper.Failed[string](c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	var input dto.BodyUpdateAccountDTO
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
