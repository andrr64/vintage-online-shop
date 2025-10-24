package dto

type ShopBase struct {
	Name        string `json:"name" binding:"required"`
	Summary     string `json:"summary" binding:"required"`
	Description string `json:"description" binding:"required"`
}
