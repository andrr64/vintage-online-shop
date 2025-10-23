package handler

import (
	"log"
	dto "vintage-server/internal/shared"
	"vintage-server/pkg/helper"
	"vintage-server/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AddToWishlist implements AccountHandler.
func (a *accountHandler) AddToWishlist(c *gin.Context) {
	accountid, err := helper.ExtractAccountID(c)
	if err != nil {
		response.ErrorUnauthorized(c)
		return
	}
	productID, err := helper.GetParamUUID(c, "product-id")
	if err != nil {
		response.ErrorBadRequest(c)
		return
	}
	err = a.services.Wishlist.AddToWishlist(c.Request.Context(), accountid, productID)
	if err != nil {
		helper.HandleError(c, err)
		return
	}
	response.SuccessWD_OK(c)
}

// GetWishlist implements AccountHandler.
func (a *accountHandler) GetWishlist(c *gin.Context) {
	accountID, err := helper.ExtractAccountID(c)
	if err != nil {
		response.ErrorUnauthorized(c)
		return
	}
	page := helper.GetQueryInt(c, "page", 1)  // default halaman 1
	size := helper.GetQueryInt(c, "size", 10) // default page size 10\
	keyword := c.Query("keyword")             // kosong jika tidak ada
	log.Printf("GetWishlist: accountID=%s, page=%d, size=%d, keyword=%s", accountID, page, size, keyword)
	productIDS, total, err := a.services.Wishlist.GetWishlist(c.Request.Context(), accountID, page, size, keyword)
	if err != nil {
		helper.HandleError(c, err)
		return
	}
	res := dto.Pageable[uuid.UUID]{
		Items:     productIDS,
		Total:     total,
		Page:      page,
		PageSize:  size,
		TotalPage: (total + size - 1) / size,
	}

	response.SuccessOK(c, res)
}

// RemoveFromWishlist implements AccountHandler.
func (a *accountHandler) RemoveFromWishlist(c *gin.Context) {
	accountid, err := helper.ExtractAccountID(c)
	if err != nil {
		response.ErrorUnauthorized(c)
		return
	}
	productID, err := helper.GetParamUUID(c, "product-id")
	if err != nil {
		response.ErrorBadRequest(c)
		return
	}
	err = a.services.Wishlist.RemoveFromWishlist(c.Request.Context(), accountid, productID)
	if err != nil {
		helper.HandleError(c, err)
		return
	}
	response.SuccessWD_OK(c)
}
