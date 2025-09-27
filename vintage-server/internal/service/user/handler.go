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

// Register adalah handler untuk use case pendaftaran customer
func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest

	// 1. Bind & Validasi request body JSON ke DTO RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	// 2. Panggil service untuk menjalankan logika bisnis
	userProfile, err := h.svc.Register(c.Request.Context(), req)
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

// Login adalah handler untuk use case login
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	loginResponse, err := h.svc.Login(c.Request.Context(), req)
	if err != nil {
		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			response.Error(c, 404, appErr.Message)
		} else {
			response.Error(c, 404, "An unexpected error occurred")
		}
		return
	}

	response.Success(c, http.StatusOK, loginResponse)
}
