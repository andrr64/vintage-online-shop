package service

import (
	"vintage-server/internal/product/repository"
	"vintage-server/pkg/uploader"
)

type ProductServices struct {
	Brand BrandService
	ProductCondition ProductConditionService
	ProductCategory ProductCategoryService
}

func NewProductServices(store repository.ProductStore, jwtSecret string, upSvc uploader.Uploader) *ProductServices {
	return &ProductServices{
		Brand: NewBrandService(store, upSvc),
		ProductCondition: NewProductConditionService(store),
		ProductCategory: NewProductCategorySvc(store),
	}
}
