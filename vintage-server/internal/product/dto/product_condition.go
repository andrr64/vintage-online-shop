package dto

type ProductConditionBase struct {
	ID   int    `json:"id"`
	Name string `json:"name" binding:"required"`
}