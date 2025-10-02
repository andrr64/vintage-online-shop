package product

import "mime/multipart"

type ProductCategory struct {
	ID   *int  `json:"id"`
	Name string `json:"name" binding:"required"`
}


type BrandRequest struct {
	Name    string  `json:"name" binding:"required"`
	LogoURL *string `json:"logo_url"`
}

type CreateBrandRequest struct {
	Name       string
	File       multipart.File
	FileHeader *multipart.FileHeader
}

type UpdateBrandRequest struct {
	Name       string
	File       multipart.File
	FileHeader *multipart.FileHeader
}
