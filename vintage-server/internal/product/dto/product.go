package dto

import "mime/multipart"

type ProductBase struct {
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