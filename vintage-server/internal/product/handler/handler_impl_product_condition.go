package handler

import (
	"vintage-server/internal/product/domain"
	"vintage-server/internal/product/dto"
	"vintage-server/pkg/helper"
	"vintage-server/pkg/response"

	"github.com/gin-gonic/gin"
)

// CreateCondition handler
func (h *productHandler) CreateCondition(c *gin.Context) {
	var req dto.ProductConditionBase
	if !helper.CheckBodyJSON(c, &req) {
		response.ErrorBadRequest(c)
		return
	}

	data := domain.ProductCondition{
		Name: req.Name,
	}

	created, err := h.svc.ProductCondition.CreateCondition(c, data)
	if err != nil {
		helper.HandleError(c, err)
		return
	}

	response.SuccessCreated(c, created)
}

// ReadConditions handler
func (h *productHandler) ReadConditions(c *gin.Context) {
	idQuery := helper.GetQueryInt16(c, "id", -1)
	var id *int16
	if idQuery != -1 {
		tmp := int16(idQuery)
		id = &tmp
	}

	conds, err := h.svc.ProductCondition.FindConditions(c, id)
	if err != nil {
		helper.HandleError(c, err)
		return
	}

	response.SuccessOK(c, conds)
}

// UpdateCondition handler
func (h *productHandler) UpdateCondition(c *gin.Context) {
	idQuery, err := helper.GetParamInt16(c, "id")
	if err != nil {
		response.ErrorBadRequest(c)
		return
	}
	id := idQuery
	var req dto.ProductConditionBase
	if !helper.CheckBodyJSON(c, &req) {
		response.ErrorBadRequest(c)
		return
	}

	data := domain.ProductCondition{
		Name: req.Name,
	}

	updated, err := h.svc.ProductCondition.UpdateCondition(c, id, data)
	if err != nil {
		helper.HandleError(c, err)
		return
	}

	response.SuccessOK(c, updated)
}

// DeleteCondition handler
func (h *productHandler) DeleteCondition(c *gin.Context) {
	idQuery, err := helper.GetParamInt16(c, "id")
	if err != nil {
		response.ErrorBadRequest(c)
		return
	}
	id := int16(idQuery)
	err = h.svc.ProductCondition.DeleteCondition(c, id)
	if err != nil {
		helper.HandleError(c, err)
		return
	}

	response.SuccessOK(c, gin.H{"message": "Product condition deleted"})
}
