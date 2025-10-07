package shop

import (
	"vintage-server/internal/domain/shop"
	"vintage-server/pkg/helper"
	"vintage-server/pkg/response"
	"vintage-server/pkg/security"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateShop(c *gin.Context) {
	_, err := security.CheckAuthAndRole(c, "seller")
	if err != nil {
		response.ErrorForbidden(c)
		return
	}

	var req shop.CreateShop
	if !helper.CheckBody(c, &req) {
		return
	}

	response.SuccessWD_OK(c)
}

func (h *Handler) UpdateShop(c *gin.Context) {
	response.SuccessWD_OK(c)
}
