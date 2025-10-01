package product

import (
	"net/http"
	"strconv"
	"vintage-server/pkg/helper"
	"vintage-server/pkg/response"

	"github.com/gin-gonic/gin"
)

// Handler adalah struct yang memegang dependency ke Service
type Handler struct {
	svc Service
}

// Perhatikan bahwa return type-nya adalah interface, bukan struct-nya langsung.
func NewHandler(svc Service) ProductHandler {
	return &Handler{svc: svc}
}

func (h *Handler) CreateCategory(c *gin.Context) {
	_, err := helper.CheckAuthAndRole(c, "admin")
	if err != nil {
		response.ErrorUnauthorized(c)
		return
	}

	var category ProductCategory
	if err := c.ShouldBindJSON(&category); err != nil {
		response.ErrorBadRequest(c, "Invalid request")
		return
	}

	err = h.svc.CreateCategory(c.Request.Context(), category)
	if err != nil {
		helper.HandleError(c, err)
		return
	}
	response.Success(c, http.StatusCreated, "OK bro")
}

func (h *Handler) ReadCategories(c *gin.Context) {
	// Deklarasikan variabel untuk hasil dan error di luar if/else
	var result interface{}
	var err error

	categoryIdStr := c.Query("id")
	if categoryIdStr == "" {
		// --- Kasus 1: Ambil semua kategori ---
		result, err = h.svc.FindAllCategories(c.Request.Context())
	} else {
		// --- Kasus 2: Ambil kategori berdasarkan ID ---
		// Konversi string ID ke integer
		var categoryId int
		categoryId, err = strconv.Atoi(categoryIdStr)
		if err != nil {
			// Jika ID bukan angka, kembalikan error bad request
			response.Error(c, http.StatusBadRequest, "Invalid category ID format")
			return
		}
		result, err = h.svc.FindById(c.Request.Context(), categoryId)
	}

	// --- Penanganan Error & Sukses (hanya ditulis sekali) ---
	if err != nil {
		helper.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, result)
}

func (h *Handler) UpdateCategory(c *gin.Context) {
	_, err := helper.CheckAuthAndRole(c, "admin")
	if err != nil {
		response.ErrorUnauthorized(c)
		return
	}
	var category ProductCategory
	if err := c.ShouldBindJSON(&category); err != nil {
		response.ErrorBadRequest(c, "Invalid request")
		return
	}
	if category.ID == nil {
		response.ErrorBadRequest(c, "Invalid request")
		return
	}

	err = h.svc.UpdateCategory(c.Request.Context(), category)
	if err != nil {
		helper.HandleError(c, err)
		return
	}

	response.SuccessWithoutData(c, http.StatusOK, "OK")
}

func (h *Handler) DeleteCategory(c *gin.Context) {
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
