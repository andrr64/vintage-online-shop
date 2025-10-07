package product

import "vintage-server/internal/domain/product"

type Handler struct {
	svc product.ProductService
}

func NewHandler(svc product.ProductService) product.ProductHandler{
	return &Handler{svc: svc}
}
