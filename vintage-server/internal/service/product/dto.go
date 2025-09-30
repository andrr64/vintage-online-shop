package product

type ProductCategory struct {
	ID   int  `json:"id"`
	Name string `json:"name" binding:"required"`
}
