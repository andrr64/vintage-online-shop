package shop

import (
	"vintage-server/internal/domain/shop"
	"vintage-server/pkg/helper"
	"vintage-server/pkg/response"
	"vintage-server/pkg/security"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc shop.ShopService
}

func NewHandler(svc shop.ShopService) shop.ShopHandler {
	return &Handler{svc: svc}
}

func (h *Handler) CreateShop(c *gin.Context) {
	accountID, err := security.CheckAuthAndRole(c, "customer")
	if err != nil {
		response.ErrorForbidden(c)
		return
	}

	var req shop.ShopRequest
	if !helper.CheckBody(c, &req) {
		return
	}
	_, err = h.svc.CreateShop(c, accountID, req)
	if err != nil {
		response.ErrorInternalServer(c, err.Error())
		return
	}
	response.SuccessWD_OK(c)
}

func (h *Handler) UpdateShop(c *gin.Context) {
	response.SuccessWD_OK(c)
}
