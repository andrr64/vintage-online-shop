package account

import (
	"vintage-server/pkg/helper"
	"vintage-server/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AddToWishlist implements account.AccountHandler.
func (h *handler) AddToWishlist(c *gin.Context) {
	accountID, err := helper.ExtractAccountID(c)
	if err != nil {
		helper.HandleError(c, err)
		return
	}
	productID, err := uuid.Parse(c.Param("product-id"))
	if err != nil {
		response.ErrorBadRequest(c, "Invalid product ID")
		return
	}
	err = h.svc.AddToWishlist(c.Request.Context(), accountID, productID)

	if err != nil {
		helper.HandleError(c, err)
		return
	}
	response.SuccessOK(c, "OK")
}

// GetWishlist implements account.AccountHandler.
func (h *handler) GetWishlist(c *gin.Context) {
	panic("unimplemented")
}

// RemoveFromWishlist implements account.AccountHandler.
func (h *handler) RemoveFromWishlist(c *gin.Context) {
	panic("unimplemented")
}
