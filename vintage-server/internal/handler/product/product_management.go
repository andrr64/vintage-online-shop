package product

import (
	"vintage-server/internal/domain/product"
	"vintage-server/pkg/helper"
	"vintage-server/pkg/response"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateProduct(c *gin.Context) {
	// Daftar role yang diizinkan
	roles := []string{"seller"}

	// Cek autentikasi dan role user
	_, err := helper.CheckAuthAndRole(c, roles...)
	if err != nil {
		response.ErrorForbiddenRoles(c, roles...)
		return
	}

	// validasi token
	accountID, _, err := helper.ExtractAccountInfoFromToken(c)

	if err != nil {
		response.ErrorInternalServer(c)
		return
	}

	var request product.CreateProductRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response.ErrorBadRequest(c)
		return
	}

	h.svc.CreateProduct(c.Request.Context(), accountID, request)

	response.SuccessWD_Created(c)
}