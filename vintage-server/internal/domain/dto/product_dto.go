package dto

type ProductCategory struct {
	ID   *int   `json:"id"`
	Name string `json:"name" binding:"required"`
}

type BrandRequest struct {
	Name    string  `json:"name" binding:"required"`
	LogoURL *string `json:"logo_url"`
}

type BrandDetail struct {
	ID      int     `json:"id" binding:"required"`
	Name    string  `json:"name" binding:"required"`
	LogoURL *string `json:"logo_url"  binding:"required"`
}
