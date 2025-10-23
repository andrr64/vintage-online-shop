package dto

import (
	"mime/multipart"
)

// BrandFormBase berisi field umum untuk operasi Brand
type BrandFormBase struct {
	Name       string                `form:"name" binding:"required,min=2"`
	FileHeader *multipart.FileHeader `form:"logo"`
}
