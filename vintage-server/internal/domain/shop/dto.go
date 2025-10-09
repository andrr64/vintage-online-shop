package shop

import "vintage-server/internal/model"

type ShopRequest struct {
	Name        string `json:"name" binding:"required"`
	Summary     string `json:"summary" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type ShopDetail struct {
	Name        string `json:"name"`
	Summary     string `json:"summary"`
	Description string `json:"description"`
}

func ConvertShopModelToShopDetail(data model.Shop) ShopDetail {
	return ShopDetail{
		Name:        data.Name,
		Summary:     *data.Summary,
		Description: *data.Description,
	}
}
