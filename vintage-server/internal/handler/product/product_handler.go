package product

import (
	"net/http"
	product "vintage-server/internal/domain"
	"vintage-server/pkg/helper"
	"vintage-server/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type handler struct {
	svc product.ProductService
}

func NewHandler(svc product.ProductService) product.ProductHandler {
	return &handler{svc: svc}
}

func (h *handler) CreateProduct(c *gin.Context) {
	// Mengekstrak ID akun dari konteks (misalnya, dari token JWT)
	accountID, err := helper.ExtractAccountID(c)
	if err != nil {
		response.ErrorUnauthorized(c, "Invalid Account")
		return
	}

	var request product.CreateProductRequest
	// FIX 1: Perbaiki typo 'shouldbind' menjadi 'ShouldBind'
	// 'ShouldBind' akan otomatis menangani multipart/form-data
	if err := c.ShouldBind(&request); err != nil {
		// FIX 2: Ganti status 502 menjadi 400 untuk error validasi dari klien
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// FIX 3: Tangkap dan tangani error yang mungkin terjadi dari service layer
	// (misalnya, error saat menyimpan ke database)
	newProduk, err := h.svc.CreateProduct(c.Request.Context(), accountID, request)
	// if err != nil {
	//     // Jika service gagal, kembalikan error 500 (Internal Server Error)
	//     response.Error(c, http.StatusInternalServerError, "Failed to create product")
	//     return
	// }
	if err != nil {
		helper.HandleError(c, err)
	}

	// FIX 4: Kembalikan data produk yang baru dibuat dengan status 201 Created
	response.Success(c, http.StatusAccepted, newProduk)
}

func (h *handler) CreateProductSize(c *gin.Context) {
	// Mengekstrak ID akun dari konteks (misalnya, dari token JWT)
	_, err := helper.ExtractAccountID(c)
	if err != nil {
		response.ErrorUnauthorized(c, "Invalid Account")
		return
	}

	var req product.ProductConditionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.HandleError(c, err)
		return
	}
	result, err := h.svc.CreateProductSize(c, req)
	if err != nil {
		helper.HandleError(c, err)
		return
	}
	response.SuccessCreated(c, result)
}

// UpdateProduct implements product.ProductHandler.
func (h *handler) UpdateProduct(c *gin.Context) {
	accountID, err := helper.ExtractAccountID(c)
	if err != nil {
		response.ErrorUnauthorized(c, "Invalid Account")
		return
	}

	var req product.UpdateProductDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.HandleError(c, err)
		return
	}
	result, err := h.svc.UpdateProduct(c, accountID, req)
	if err != nil {
		helper.HandleError(c, err)
		return
	}
	response.SuccessCreated(c, result)
}

func (h *handler) GetProuctByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, "Product ID is required")
		return
	}
	productID, err := uuid.Parse(id)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid Product ID format")
		return
	}
	result, err := h.svc.FindProductByID(c.Request.Context(), productID)
	if err != nil {
		helper.HandleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, result)
}
