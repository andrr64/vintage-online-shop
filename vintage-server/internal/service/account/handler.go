package account

import (
	"errors"
	"net/http"
	"strconv"
	"vintage-server/pkg/apperror"
	"vintage-server/pkg/response"
	"vintage-server/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler adalah struct yang memegang dependency ke Service
type Handler struct {
	svc Service
}

// Perhatikan bahwa return type-nya adalah interface, bukan struct-nya langsung.
func NewHandler(svc Service) AccountHandler {
	return &Handler{svc: svc}
}


// handleError adalah helper internal untuk menangani error dari service secara konsisten
func (h *Handler) handleError(c *gin.Context, err error) {
	var appErr *apperror.AppError
	if errors.As(err, &appErr) {
		response.Error(c, appErr.Code, appErr.Message)
	} else {
		// Sembunyikan detail error internal dari client
		response.Error(c, http.StatusInternalServerError, "an unexpected internal error occurred")
	}
}

// RegisterCustomer adalah handler untuk use case pendaftaran customer
func (h *Handler) RegisterCustomer(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	userProfile, err := h.svc.RegisterCustomer(c.Request.Context(), req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, userProfile)
}

// LoginCustomer adalah handler untuk use case login
func (h *Handler) LoginCustomer(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	loginResponse, err := h.svc.LoginCustomer(c.Request.Context(), req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.SetCookie("access_token", loginResponse.AccessToken, 3600*24, "/", "", false, true)
	response.Success(c, http.StatusOK, loginResponse)
}

func (h *Handler) UpdateProfile(c *gin.Context) {
	// Gunakan helper untuk otentikasi dan otorisasi role "customer"
	accountID, err := checkAuthAndRole(c, "customer")
	if err != nil {
		h.handleError(c, err)
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	res, err := h.svc.UpdateProfile(c.Request.Context(), accountID, req)
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, res)
}

func (h *Handler) UpdateAvatar(c *gin.Context) {
	// Hanya perlu memastikan user sudah login, role tidak spesifik
	accountID, err := checkAuthAndRole(c)
	if err != nil {
		h.handleError(c, err)
		return
	}

	fileHeader, err := c.FormFile("avatar")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "avatar file is required")
		return
	}

	// Perbaiki logika: SizeIsOk harusnya SizeIsNotOk atau if !SizeIsOk
	if !utils.SizeIsOk(fileHeader, utils.BytesToMegaBytes(2)) {
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
		h.handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, res)
}

// -- ADDRESS MANAGEMENT --

func (h *Handler) CreateAddress(c *gin.Context) {
	accountID, err := checkAuthAndRole(c, "customer")
	if err != nil {
		h.handleError(c, err)
		return
	}

	var req AddAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	res, err := h.svc.AddAddress(c.Request.Context(), accountID, req)
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.Success(c, http.StatusCreated, res)
}

func (h *Handler) GetAddresses(c *gin.Context) {
	accountID, err := checkAuthAndRole(c, "customer")
	if err != nil {
		h.handleError(c, err)
		return
	}

	addressIDStr := c.Param("addressId") // Ambil dari path, e.g., /addresses/123
	if addressIDStr == "" {
		// Ambil semua alamat jika tidak ada ID
		addresses, err := h.svc.GetAddressesByUserID(c.Request.Context(), accountID)
		if err != nil {
			h.handleError(c, err)
			return
		}
		response.Success(c, http.StatusOK, addresses)
	} else {
		// Ambil alamat spesifik jika ada ID
		addressID, err := strconv.ParseInt(addressIDStr, 10, 64)
		if err != nil {
			response.Error(c, http.StatusBadRequest, "invalid address ID format")
			return
		}

		address, err := h.svc.GetAddressByID(c.Request.Context(), accountID, addressID)
		if err != nil {
			h.handleError(c, err)
			return
		}
		response.Success(c, http.StatusOK, address)
	}
}

func (h *Handler) UpdateAddress(c *gin.Context) {
	accountID, err := checkAuthAndRole(c, "customer")
	if err != nil {
		h.handleError(c, err)
		return
	}

	var req UserAddress
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	res, err := h.svc.UpdateAddress(c.Request.Context(), accountID, req)
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, res)
}

func (h *Handler) DeleteAddress(c *gin.Context) {
	accountID, err := checkAuthAndRole(c, "customer")
	if err != nil {
		h.handleError(c, err)
		return
	}

	addressIDStr := c.Param("addressId")
	addressID, err := strconv.ParseInt(addressIDStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid address ID format")
		return
	}

	err = h.svc.DeleteAddress(c.Request.Context(), accountID, addressID)
	if err != nil {
		h.handleError(c, err)
		return // Pastikan ada return setelah error
	}
	response.SuccessWithoutData(c, http.StatusOK, "address deleted successfully")
}

func (h *Handler) SetPrimaryAddress(c *gin.Context) {
	accountID, err := checkAuthAndRole(c, "customer")
	if err != nil {
		h.handleError(c, err)
		return
	}

	var req AddressIdentifier
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	err = h.svc.SetPrimaryAddress(c.Request.Context(), accountID, req.AddressID)
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.SuccessWithoutData(c, http.StatusOK, "primary address updated successfully")
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

func (h *Handler) LoginSeller(c *gin.Context) {
		var req LoginRequest

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