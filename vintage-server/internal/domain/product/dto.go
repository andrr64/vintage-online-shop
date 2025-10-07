package product

import "mime/multipart"

type ProductCategory struct {
	ID   *int   `json:"id"`
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

type ProductConditionRequest struct {
	Name string `json:"name" binding:"required"`
}

type CreateProductRequest struct {
	Name        string                  `form:"name" binding:"required"`
	CategoryID  int                     `form:"category_id" binding:"required"`
	ConditionID int16                   `form:"condition_id" binding:"required"`
	Price       int64                   `form:"price" binding:"required,gt=0"`
	Stock       int                     `form:"stock" binding:"required,gte=0"`
	Description string                  `form:"description"`
	Summary     string                  `form:"summary"`
	BrandID     *int                    `form:"brand_id"`
	SizeID      *int                    `form:"size_id"`
	Images      []*multipart.FileHeader `form:"images" binding:"required,min=1,max=5"` // Validasi: min 1, max 5 gambar
}
