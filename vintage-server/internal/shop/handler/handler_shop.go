package handler

import (
	"vintage-server/internal/shop/domain"
	"vintage-server/internal/shop/dto"
	"vintage-server/pkg/helper"
	"vintage-server/pkg/response"

	"github.com/gin-gonic/gin"
)

// CreateShop implements ShopHandler.
func (s *shopHandler) CreateShop(c *gin.Context) {
	accId, err := helper.ExtractAccountID(c)
	if err != nil {
		response.ErrorUnauthorized(c)
		return
	}
	var req dto.ShopBase
	if !helper.BindJSON(c, &req){
		return
	}
	data := domain.Shop{
		AccountID: accId,
		Name: req.Name,
		Description: &req.Description,
		Summary: &req.Summary,
	}
	ctx := c.Request.Context()
	res, err := s.svc.Shop.CreateShop(ctx, data)
	if err != nil {
		helper.HandleError(c, err)
		return
	}
	response.SuccessOK(c, res)
}

// UpdateShop implements ShopHandler.
func (s *shopHandler) UpdateShop(c *gin.Context) {
	panic("unimplemented")
}
