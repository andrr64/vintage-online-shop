package product

import (
	"net/http"
	"strconv"
	"vintage-server/internal/domain/product"
	"vintage-server/pkg/helper"
	"vintage-server/pkg/response"

	"github.com/gin-gonic/gin"
)

// -- PRODUCT CONDITION MANAGEMENET --
func (h *Handler) CreateCondition(c *gin.Context) {
	_, err := helper.CheckAuthAndRole(c, "admin")
	if err != nil {
		response.ErrorUnauthorized(c)
		return
	}

	var req product.ProductConditionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.HandleErrorBadRequest(c)
		return
	}

	newCondition, err := h.svc.CreateCondition(c.Request.Context(), req)
	if err != nil {
		helper.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, newCondition)
}

func (h *Handler) ReadConditions(c *gin.Context) {
	idStr := c.Query("id")

	// Jika ada query parameter 'id', cari satu data
	if idStr != "" {
		id, err := strconv.ParseInt(idStr, 10, 16)
		if err != nil {
			response.Error(c, http.StatusBadRequest, "Invalid ID format")
			return
		}

		condition, err := h.svc.FindConditionByID(c.Request.Context(), int16(id))
		if err != nil {
			helper.HandleError(c, err)
			return
		}
		response.Success(c, http.StatusOK, condition)
		return
	}

	// Jika tidak ada query 'id', ambil semua data
	conditions, err := h.svc.FindAllConditions(c.Request.Context())
	if err != nil {
		helper.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, conditions)
}

func (h *Handler) UpdateCondition(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 16)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ID format")
		return
	}

	var req product.ProductConditionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.HandleErrorBadRequest(c)
		return
	}

	updatedCondition, err := h.svc.UpdateCondition(c.Request.Context(), int16(id), req)
	if err != nil {
		helper.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, updatedCondition)
}

func (h *Handler) DeleteCondition(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 16)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ID format")
		return
	}

	err = h.svc.DeleteCondition(c.Request.Context(), int16(id))
	if err != nil {
		helper.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "Product condition deleted successfully"})
}
