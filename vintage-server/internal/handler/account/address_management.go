package account

import (
	"net/http"
	"strconv"
	"vintage-server/internal/domain/account"
	"vintage-server/pkg/helper"
	"vintage-server/pkg/response"

	"github.com/gin-gonic/gin"
)

// -- ADDRESS MANAGEMENT --

func (h *handler) CreateAddress(c *gin.Context) {
	accountID, err := helper.CheckAuthAndRole(c, "customer")
	if err != nil {
		helper.HandleError(c, err)
		return
	}

	var req account.AddAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	res, err := h.svc.AddAddress(c.Request.Context(), accountID, req)
	if err != nil {
		helper.HandleError(c, err)
		return
	}
	response.Success(c, http.StatusCreated, res)
}

func (h *handler) GetAddresses(c *gin.Context) {
	accountID, err := helper.CheckAuthAndRole(c, "customer")
	if err != nil {
		helper.HandleError(c, err)
		return
	}

	addressIDStr := c.Param("addressId") // Ambil dari path, e.g., /addresses/123
	if addressIDStr == "" {
		// Ambil semua alamat jika tidak ada ID
		addresses, err := h.svc.GetAddressesByUserID(c.Request.Context(), accountID)
		if err != nil {
			helper.HandleError(c, err)
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
			helper.HandleError(c, err)
			return
		}
		response.Success(c, http.StatusOK, address)
	}
}

func (h *handler) UpdateAddress(c *gin.Context) {
	accountID, err := helper.CheckAuthAndRole(c, "customer")
	if err != nil {
		helper.HandleError(c, err)
		return
	}

	var req account.UserAddress
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	res, err := h.svc.UpdateAddress(c.Request.Context(), accountID, req)
	if err != nil {
		helper.HandleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, res)
}

func (h *handler) DeleteAddress(c *gin.Context) {
	accountID, err := helper.CheckAuthAndRole(c, "customer")
	if err != nil {
		helper.HandleError(c, err)
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
		helper.HandleError(c, err)
		return // Pastikan ada return setelah error
	}
	response.SuccessWithoutData(c, http.StatusOK, "address deleted successfully")
}

func (h *handler) SetPrimaryAddress(c *gin.Context) {
	accountID, err := helper.CheckAuthAndRole(c, "customer")
	if err != nil {
		helper.HandleError(c, err)
		return
	}

	var req account.AddressIdentifier
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	err = h.svc.SetPrimaryAddress(c.Request.Context(), accountID, req.AddressID)
	if err != nil {
		helper.HandleError(c, err)
		return
	}
	response.SuccessWithoutData(c, http.StatusOK, "primary address updated successfully")
}
