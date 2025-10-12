package account

import (
	"net/http"
	"vintage-server/internal/domain/account"
	"vintage-server/pkg/helper"
	"vintage-server/pkg/response"
	"vintage-server/pkg/utils"

	"github.com/gin-gonic/gin"
)

func (h *handler) UpdateProfile(c *gin.Context) {
	// Gunakan helper untuk otentikasi dan otorisasi role "customer"
	accountID, err := helper.CheckAuthAndRole(c, "customer")
	if err != nil {
		helper.HandleError(c, err)
		return
	}

	var req account.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	res, err := h.svc.UpdateProfile(c.Request.Context(), accountID, req)
	if err != nil {
		helper.HandleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, res)
}

func (h *handler) UpdateAvatar(c *gin.Context) {
	// Hanya perlu memastikan user sudah login, role tidak spesifik
	accountID, err := helper.CheckAuthAndRole(c)
	if err != nil {
		helper.HandleError(c, err)
		return
	}

	fileHeader, err := c.FormFile("avatar")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "avatar file is required")
		return
	}

	// Perbaiki logika: SizeIsOk harusnya SizeIsNotOk atau if !SizeIsOk
	if !utils.SizeIsOk(fileHeader, utils.Megabytes(2)) {
		response.Error(c, http.StatusBadRequest, "file size must be less than 2MB")
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to open uploaded file")
		return
	}
	defer file.Close()

	res, err := h.svc.UpdateAvatar(c.Request.Context(), accountID, file)
	if err != nil {
		helper.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, res)
}
