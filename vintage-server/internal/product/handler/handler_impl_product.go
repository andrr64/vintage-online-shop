package handler

import (
	"vintage-server/internal/product/dto"
	"vintage-server/pkg/helper"

	"github.com/gin-gonic/gin"
)

func (h *productHandler) CreateProduct(c *gin.Context) {
	var req dto.ProductBase
	if !helper.CheckBody(c, &req){
		return
	}
	
}

func (h *productHandler) UpdateProduct(c *gin.Context) {

}

func (h *productHandler) GetProuctByID(c *gin.Context) {

}

func (h *productHandler) SellerGetProducts(c *gin.Context) {

}
