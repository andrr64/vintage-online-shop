package product

import "mime/multipart"

type BrandRequest struct {
	Name       string
	File       multipart.File
	FileHeader *multipart.FileHeader
}

type CreateProductRequest struct {
	Name        string                  `form:"name" binding:"required"`
	CategoryID  int                     `form:"category_id" binding:"required"`
	ConditionID int16                   `form:"condition_id" binding:"required"`
	Price       int64                   `form:"price" binding:"required,gt=0"`
	Stock       int                     `form:"stock" binding:"required,gte=0"`
	Description string                  `form:"description"`
	Summary     string                  `form:"summary"`
	BrandID     int                     `form:"brand_id"`
	SizeID      int                     `form:"size_id"`
	Thumbnail   *multipart.FileHeader   `form:"thumbnail" binding:"required"`      // wajib 1
	Images      []*multipart.FileHeader `form:"images" binding:"omitempty,max=10"` // opsional max 10
}

// ✅ DTO untuk update produk
type UpdateProductDTO struct {
	ID          string  `form:"id" binding:"required,uuid"`
	Name        *string `form:"name"`
	CategoryID  *int    `form:"category_id"`
	ConditionID *int16  `form:"condition_id"`
	Price       *int64  `form:"price"`
	Stock       *int    `form:"stock"`
	Description *string `form:"description"`
	Summary     *string `form:"summary"`
	BrandID     *int    `form:"brand_id"`
	SizeID      *int    `form:"size_id"`
}

type UpdateBrandRequest BrandRequest
type ProductConditionRequest ProductConditionDTO
