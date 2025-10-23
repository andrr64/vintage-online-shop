package handler

import (
	"mime/multipart"
	"net/http"
	"vintage-server/internal/product/domain"
	"vintage-server/internal/product/dto"
	"vintage-server/pkg/helper"
	"vintage-server/pkg/response"
	"vintage-server/pkg/utils"

	"github.com/gin-gonic/gin"
)

// CreateBrand membuat brand baru
func (h *productHandler) CreateBrand(c *gin.Context) {
	var form dto.BrandFormBase
	if !helper.CheckBody(c, &form) {
		return
	}

	if form.FileHeader == nil {
		response.Error(c, http.StatusBadRequest, "logo file is required")
		return
	}

	if !utils.SizeIsOk(form.FileHeader, utils.Megabytes(2)) {
		response.Error(c, http.StatusBadRequest, "file size must be less than 2MB")
		return
	}
	file, err := form.FileHeader.Open()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to open uploaded file")
		return
	}
	defer file.Close()

	data := domain.Brand{
		Name: form.Name,
	}

	saved, err := h.svc.Brand.CreateBrand(c, data, file)
	if err != nil {
		helper.HandleError(c, err)
		return
	}
	response.SuccessCreated(c, saved)
}

// ReadBrand baca brand (semua / by ID query param)
func (h *productHandler) ReadBrand(c *gin.Context) {
	var brands []domain.Brand
	var err error
	id := helper.GetQueryInt(c, "id", -1)

	if id == -1 {
		brands, err = h.svc.Brand.ReadBrands(c, nil)
	} else {
		brands, err = h.svc.Brand.ReadBrands(c, &id)
	}

	if err != nil {
		helper.HandleError(c, err)
		return
	}
	response.SuccessOK(c, brands)
}

// UpdateBrand update data brand
func (h *productHandler) UpdateBrand(c *gin.Context) {
	id, err := helper.GetParamInt(c, "id")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid brand ID")
		return
	}

	var form dto.BrandFormBase
	if err := c.ShouldBind(&form); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	var f multipart.File

	// File opsional
	if form.FileHeader != nil {
		if !utils.SizeIsOk(form.FileHeader, utils.Megabytes(2)) {
			response.Error(c, http.StatusBadRequest, "file size must be less than 2MB")
			return
		}

		f, err = form.FileHeader.Open()
		if err != nil {
			response.Error(c, http.StatusInternalServerError, "failed to open uploaded file")
			return
		}
		defer f.Close()
	}
	data := domain.Brand{
		ID: id,
		Name: form.Name,
	}
	updated, err := h.svc.Brand.UpdateBrand(c, data, &f)
	if err != nil {
		helper.HandleError(c, err)
		return
	}
	response.SuccessOK(c, updated)
}

// DeleteBrand hapus brand
func (h *productHandler) DeleteBrand(c *gin.Context) {
	id, err := helper.GetParamInt(c, "id")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid brand ID")
		return
	}

	err = h.svc.Brand.DeleteBrand(c, id)
	if err != nil {
		helper.HandleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "OK")
}
