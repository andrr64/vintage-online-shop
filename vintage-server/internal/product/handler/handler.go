package handler

import (
	"vintage-server/internal/product/service"

	"github.com/gin-gonic/gin"
)

type ProductHandler interface {
	// Manajemen Kategori
	CreateCategory(c *gin.Context)
	ReadCategories(c *gin.Context)
	UpdateCategory(c *gin.Context)
	DeleteCategory(c *gin.Context)

	// -- brand management --
	CreateBrand(c *gin.Context)
	ReadBrand(c *gin.Context)
	UpdateBrand(c *gin.Context)
	DeleteBrand(c *gin.Context)

	// -- PRODUCT CONDITION MANAGEMENT --
	CreateCondition(c *gin.Context)
	ReadConditions(c *gin.Context)
	UpdateCondition(c *gin.Context)
	DeleteCondition(c *gin.Context)

	// -- PRODUCT MANAGEMENT --
	CreateProduct(c *gin.Context)
	UpdateProduct(c *gin.Context)
	GetProuctByID(c *gin.Context)
	SellerGetProducts(c *gin.Context)

	// -- PRODUCT SIZE MANAGEMENT --
	CreateProductSize(c *gin.Context)
}

type productHandler struct {
	svc service.ProductServices
}

func NewProductHandler(svc service.ProductServices) ProductHandler {
	return &productHandler{svc: svc}
}
