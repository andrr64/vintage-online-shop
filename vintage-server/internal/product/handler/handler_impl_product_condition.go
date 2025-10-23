package handler

import (
	"vintage-server/internal/product/domain"
	"vintage-server/internal/product/dto"
	"vintage-server/pkg/helper"
	"vintage-server/pkg/response"

	"github.com/gin-gonic/gin"
)

func (h *productHandler) CreateCondition(c *gin.Context) {
	var req dto.ProductConditionBase

	if !helper.CheckBodyJSON(c, &req) {
		response.ErrorBadRequest(c)
		return
	}

	prodConditionDomain := domain.ProductCondition{
		Name: req.Name,
	}

	response.SuccessCreated(c, prodConditionDomain)
}

func (h *productHandler) ReadConditions(c *gin.Context) {

}
func (h *productHandler) UpdateCondition(c *gin.Context) {

}
func (h *productHandler) DeleteCondition(c *gin.Context) {

}
