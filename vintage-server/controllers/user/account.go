package user

import (
	"net/http"
	"vintage-server/dto"
	"vintage-server/helpers/response_helper"

	"github.com/gin-gonic/gin"
)

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
// @Router /user/account [get]
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
// @Router /user/account/password [put]
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
// @Router /user/account/profile [put]
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
