package product

import (
	"net/http"
	"strconv"
	product "vintage-server/internal/domain"
	"vintage-server/pkg/helper"
	"vintage-server/pkg/response"

	"github.com/gin-gonic/gin"
)

// -- CATEGORY MANAGEMENT --
func (h *handler) CreateCategory(c *gin.Context) {
	var category product.ProductCategoryDTO
	if err := c.ShouldBindJSON(&category); err != nil {
		helper.HandleErrorBadRequest(c)
		return
	}

	err := h.svc.CreateCategory(c.Request.Context(), category)
	if err != nil {
		helper.HandleError(c, err)
		return
	}
	response.SuccessWD_Created(c)
}

func (h *handler) ReadCategories(c *gin.Context) {
	var result interface{}
	var err error

	categoryIdStr := c.Param("id")
	if categoryIdStr == "" {
		result, err = h.svc.FindAllCategories(c.Request.Context())
	} else {
		var categoryId int
		categoryId, err = strconv.Atoi(categoryIdStr)
		if err != nil {
			response.Error(c, http.StatusBadRequest, "Invalid category ID format")
			return
		}
		result, err = h.svc.FindById(c.Request.Context(), categoryId)
	}
	if err != nil {
		helper.HandleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, result)
}

func (h *handler) UpdateCategory(c *gin.Context) {
	var category product.ProductCategoryDTO
	if err := c.ShouldBindJSON(&category); err != nil {
		helper.HandleErrorBadRequest(c)
		return
	}
	if category.ID == nil {
		response.ErrorBadRequest(c, "Invalid request")
		return
	}

	err := h.svc.UpdateCategory(c.Request.Context(), category)
	if err != nil {
		helper.HandleError(c, err)
		return
	}

	response.SuccessWithoutData(c, http.StatusOK, "OK")
}

func (h *handler) DeleteCategory(c *gin.Context) {
	_, err := helper.CheckAuthAndRole(c, "admin")
	if err != nil {
		response.ErrorUnauthorized(c)
		return
	}
	id := c.Query("id")
	if id == "" {
		response.ErrorBadRequest(c)
		return
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		response.ErrorBadRequest(c)
	}
	err = h.svc.DeleteCategory(c.Request.Context(), idInt)
	if err != nil {
		helper.HandleError(c, err)
		return
	}
	response.SuccessWithoutData(c, http.StatusOK, "OK")
}
