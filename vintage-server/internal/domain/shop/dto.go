package shop

type CreateShop struct {
	Name        string `json:"name" binding:"required"`
	Summary     string `json:"summary" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type ShopDetail struct {
	Name        string `json:"name"`
	Summary     string `json:"summary"`
	Description string `json:"description"`
}