package user

import (
	"errors"
	"net/http"
	"vintage-server/pkg/apperror" // Path ke package error kustom kita
	"vintage-server/pkg/helper"
	"vintage-server/pkg/response" // Path ke package error kustom kita
	"vintage-server/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler adalah struct yang memegang dependency ke Service
type Handler struct {
	svc Service
}

// NewHandler adalah constructor untuk handler
func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterCustomer adalah handler untuk use case pendaftaran customer
func (h *Handler) RegisterCustomer(c *gin.Context) {
	var req RegisterRequest

	// 1. Bind & Validasi request body JSON ke DTO RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	// 2. Panggil service untuk menjalankan logika bisnis
	userProfile, err := h.svc.RegisterCustomer(c.Request.Context(), req)
	if err != nil {
		// 3. Tangani error dari service menggunakan custom error kita
		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			response.Error(c, appErr.Code, appErr.Message)
		} else {
			response.Error(c, http.StatusInternalServerError, "An unexpected error occurred")
		}
		return
	}

	// 4. Sukses
	response.Success(c, http.StatusCreated, userProfile)
}

// LoginCustomer adalah handler untuk use case login
func (h *Handler) LoginCustomer(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	loginResponse, err := h.svc.LoginCustomer(c.Request.Context(), req)
	if err != nil {
		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			response.Error(c, 404, appErr.Message)
		} else {
			response.Error(c, 404, "An unexpected error occurred")
		}
		return
	}

	c.SetCookie(
		"access_token",
		loginResponse.AccessToken,
		3600*24, // 3 hari
		"/",     // path
		"",      // domain (atau kosong "")
		false,   // secure (true kalau https)
		true,    // httpOnly biar gak bisa diakses JS
	)
	response.Success(c, http.StatusOK, loginResponse)
}

func (h *Handler) UpdateProfile(c *gin.Context) {
	accountID, role, err := helper.ExtractAccountInfoFromToken(c)
	if err != nil || role == nil {
		response.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	if *role == "admin" {
		response.Error(c, http.StatusForbidden, "Admin is not allowed to update profile here")
		return
	}

	// Bind request body
	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request")
		return
	}

	res, err := h.svc.UpdateProfile(c.Request.Context(), accountID, req)
	if err != nil {
		response.Error(c, apperror.ErrCodeInternal, "Something wrong when we try to update your data")
		return
	}
	response.Success(c, http.StatusOK, res)
}

func (h *Handler) UpdateAvatar(c *gin.Context) {
	accountID, _, err := helper.ExtractAccountInfoFromToken(c)

	if err != nil {
		response.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// ambil file dari request
	fileHeader, err := c.FormFile("avatar")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Avatar file is required")
		return
	}

	// validasi ukuran file (max 2MB)
	if utils.SizeIsOk(fileHeader, utils.BytesToMegaBytes(2)) {
		response.Error(c, http.StatusBadRequest, "File size must be less than 2MB")
		return
	}

	// buka file
	file, err := fileHeader.Open()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to open uploaded file")
		return
	}
	defer file.Close()

	// panggil service
	res, err := h.svc.UpdateAvatar(c.Request.Context(), accountID, file)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, res)
}

func (h *Handler) CreateAddress(c *gin.Context) {
	var address AddAddressRequest
	if err := c.ShouldBindJSON(&address); err != nil {
		response.ErrorBadRequest(c)
		return
	}
	accountId, role, err := helper.ExtractAccountInfoFromToken(c)
	if err != nil {
		response.ErrorUnauthorized(c, err.Error())
		return
	}
	if *role != "customer" {
		response.ErrorForbidden(c)
		return
	}

	res, err := h.svc.AddAddress(c.Request.Context(), accountId, address)
	if err != nil {
		response.ErrorInternalServer(c, err.Error())
		return
	}
	response.Success(c, http.StatusCreated, res)
}

func (h *Handler) Logout(c *gin.Context) {
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

func (h *Handler) LoginAdmin(c *gin.Context) {
	var req LoginRequest

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
