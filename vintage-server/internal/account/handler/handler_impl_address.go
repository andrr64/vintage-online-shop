// handlerr_impl_address.go
package handler

import (
	"log"
	"vintage-server/internal/account/dto"
	"vintage-server/pkg/helper"
	"vintage-server/pkg/response"

	"github.com/gin-gonic/gin"
)

// SetPrimaryAddress implements AccountHandler.
func (a *accountHandler) SetPrimaryAddress(c *gin.Context) {
	accountID, err := helper.ExtractAccountID(c)

	if err != nil {
		response.ErrorUnauthorized(c)
		return
	}
	addressID, err := helper.GetParamInt64(c, "address-id")
	if err != nil {
		response.ErrorBadRequest(c, "invalid address")
	}

	err = a.services.Address.SetPrimaryAddress(c.Request.Context(), accountID, addressID)
	if err != nil {
		helper.HandleError(c, err)
		return
	}
	response.SuccessWD_OK(c)

}

// UpdateAddress implements AccountHandler.
func (a *accountHandler) UpdateAddress(c *gin.Context) {
	accountID, err := helper.ExtractAccountID(c)
	if err != nil {
		response.ErrorUnauthorized(c)
		return
	}
	var req dto.AddressDTO
	if !helper.CheckBodyJSON(c, &req) {
		response.ErrorBadRequest(c)
		return
	}
	updatedAddress, err := a.services.Address.UpdateAddress(c.Request.Context(), dto.ConvertAddressDTOToDomain(req, accountID))
	if err != nil {
		helper.HandleError(c, err)
		return
	}
	response.SuccessOK(c, updatedAddress)
}

// CreateAddress implements AccountHandler.
func (a *accountHandler) CreateAddress(c *gin.Context) {
	accountID, err := helper.ExtractAccountID(c)

	if err != nil {
		response.ErrorUnauthorized(c)
		return
	}

	var req dto.AddressDTO
	if !helper.CheckBodyJSON(c, &req) {
		response.ErrorBadRequest(c)
		return
	}
	addressDomain := dto.ConvertAddressDTOToDomain(req, accountID)
	err = a.services.Address.AddAddress(c.Request.Context(), accountID, addressDomain)
	if err != nil {
		helper.HandleError(c, err)
		return
	}
	response.SuccessWD_Created(c)
}

// DeleteAddress implements AccountHandler.
func (a *accountHandler) DeleteAddress(c *gin.Context) {
	accountID, err := helper.ExtractAccountID(c)

	if err != nil {
		response.ErrorUnauthorized(c)
		return
	}

	addressID, err := helper.GetParamInt64(c, "address-id")
	if err != nil {
		response.ErrorBadRequest(c, "invalid address")
	}
	err = a.services.Address.DeleteAddress(c.Request.Context(), accountID, addressID)
	if err != nil {
		helper.HandleError(c, err)
		return
	}
	response.SuccessWD_OK(c)
}

// GetAddresses implements AccountHandler.
func (a *accountHandler) GetAddresses(c *gin.Context) {
	accountID, err := helper.ExtractAccountID(c)
	var addresses []dto.AddressDTO

	if err != nil {
		response.ErrorUnauthorized(c)
		return
	}
	addressID := helper.GetQueryInt64(c, "address-id", -1)

	if addressID == -1 {
		addresses, err = a.services.Address.GetAddresses(c.Request.Context(), accountID, nil)
	} else {
		addresses, err = a.services.Address.GetAddresses(c.Request.Context(), accountID, &addressID)
	}
	if err != nil {
		helper.HandleError(c, err)
		return
	}
	log.Printf("%v", len(addresses))
	response.SuccessOK(c, addresses)
}
