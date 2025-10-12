package account

import (
	"errors"
	"net/http"
	"vintage-server/internal/domain/account"
	"vintage-server/pkg/apperror"
	"vintage-server/pkg/helper"
	"vintage-server/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *handler) Logout(c *gin.Context) {
	accountIDv, exists := c.Get("accountID")
	if !exists {
		// Ini adalah kasus aneh, seharusnya tidak terjadi jika middleware dipasang benar
		response.Error(c, http.StatusUnauthorized, "User ID not found in context")
		return
	}

	accountID, ok := accountIDv.(uuid.UUID)
	if !ok {
		// Error ini mengindikasikan ada masalah programming, bukan input user
		response.Error(c, http.StatusInternalServerError, "Invalid user ID format in context")
		return
	}

	_, err := h.svc.Logout(c.Request.Context(), accountID)

	if err != nil {
		response.Error(c, http.StatusInternalServerError, "internal server error")
		return
	}

	c.SetCookie("access_token", "", -1, "/", "", false, true)
	response.Success(c, http.StatusOK, "Logged out successfully")
}

func (h *handler) LoginAdmin(c *gin.Context) {
	var req account.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorBadRequest(c)
		return
	}

	loginResponse, err := h.svc.LoginAdmin(c.Request.Context(), req)

	if err != nil {
		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			response.Error(c, appErr.Code, appErr.Message)
		} else {
			response.Error(c, http.StatusInternalServerError, "An unexpected error occurred")
		}
		return
	}
	c.SetCookie(
		"access_token",
		loginResponse.AccessToken,
		3600*72, // 3 hari
		"/",     // path
		"",      // domain (atau kosong "")
		false,   // secure (true kalau https)
		true,    // httpOnly biar gak bisa diakses JS
	)

	response.Success(c, http.StatusOK, loginResponse)
}

func (h *handler) LoginSeller(c *gin.Context) {
	var req account.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}
	loginResponse, err := h.svc.LoginSeller(c.Request.Context(), req)

	if err != nil {
		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			response.Error(c, appErr.Code, appErr.Message)
		} else {
			response.Error(c, http.StatusInternalServerError, "An unexpected error occurred")
		}
		return
	}
	c.SetCookie(
		"access_token",
		loginResponse.AccessToken,
		3600*72, // 3 hari
		"/",     // path
		"",      // domain (atau kosong "")
		false,   // secure (true kalau https)
		true,    // httpOnly biar gak bisa diakses JS
	)

	response.Success(c, http.StatusOK, loginResponse)
}

// RegisterCustomer adalah handler untuk use case pendaftaran customer
func (h *handler) RegisterCustomer(c *gin.Context) {
	var req account.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	userProfile, err := h.svc.RegisterCustomer(c.Request.Context(), req)
	if err != nil {
		helper.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, userProfile)
}

// LoginCustomer adalah handler untuk use case login
func (h *handler) LoginCustomer(c *gin.Context) {
	var req account.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	loginResponse, err := h.svc.LoginCustomer(c.Request.Context(), req)
	if err != nil {
		helper.HandleError(c, err)
		return
	}

	c.SetCookie("access_token", loginResponse.AccessToken, 3600*24, "/", "", false, true)
	response.Success(c, http.StatusOK, loginResponse)
}
