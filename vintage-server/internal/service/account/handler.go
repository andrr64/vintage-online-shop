package user

import (
	"errors"
	"net/http"
	"vintage-server/pkg/apperror" // Path ke package error kustom kita
	"vintage-server/pkg/response" // Path ke package error kustom kita

	"github.com/gin-gonic/gin"
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
		3600*72, // 3 hari
		"/",     // path
		"",      // domain (atau kosong "")
		false,   // secure (true kalau https)
		true,    // httpOnly biar gak bisa diakses JS
	)

	response.Success(c, http.StatusOK, loginResponse)
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
