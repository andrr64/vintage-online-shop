package handler

import (
	"strconv"
	"vintage-server/internal/product/domain"
	"vintage-server/internal/product/dto"
	"vintage-server/pkg/helper"
	"vintage-server/pkg/response"

	"github.com/gin-gonic/gin"
)

// CreateCategory godoc
func (h *productHandler) CreateCategory(c *gin.Context) {
	var req dto.ProductCategoryBase
	if !helper.BindJSON(c, &req) {
		return
	}
	data := domain.ProductCategory{
		Name: req.Name,
	}
	category, err := h.svc.ProductCategory.CreateCategory(c.Request.Context(), data)
	if err != nil {
		helper.HandleError(c, err)
		return
	}

	response.SuccessCreated(c, category)
}

// ReadCategories godoc
func (h *productHandler) ReadCategories(c *gin.Context) {
	var id *int
	if idParam := c.Query("id"); idParam != "" {
		val, err := strconv.Atoi(idParam)
		if err != nil {
			helper.HandleError(c, err)
			return
		}
		id = &val
	}
	categories, err := h.svc.ProductCategory.FindCategories(c.Request.Context(), id)
	if err != nil {
		helper.HandleError(c, err)
		return
	}

	response.SuccessOK(c, categories)
}

// UpdateCategory godoc
func (h *productHandler) UpdateCategory(c *gin.Context) {
	var req dto.ProductCategoryBase
	if !helper.BindJSON(c, &req) {
		return
	}
	if req.ID == 0 {
		response.ErrorBadRequest(c)
		return
	}

	data := domain.ProductCategory{
		ID:   req.ID,
		Name: req.Name,
	}

	updated, err := h.svc.ProductCategory.UpdateCategory(c.Request.Context(), data)
	if err != nil {
		helper.HandleError(c, err)
		return
	}

	response.SuccessOK(c, updated)
}

// DeleteCategory godoc
func (h *productHandler) DeleteCategory(c *gin.Context) {
	idParam := c.Param("id")
	categoryID, err := strconv.Atoi(idParam)
	if err != nil {
		helper.HandleError(c, err)
		return
	}

	if err := h.svc.ProductCategory.DeleteCategory(c.Request.Context(), categoryID); err != nil {
		helper.HandleError(c, err)
		return
	}

	response.SuccessOK(c, gin.H{"message": "product category deleted successfully"})
}
