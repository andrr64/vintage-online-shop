package product

import (
	"vintage-server/internal/model"

	"github.com/google/uuid"
)

type ProductImageDTO struct {
	Index int    `json:"index"`
	URL   string `json:"url"`
}

type ProductConditionDTO struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ProductSizeDTO struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ProductCategoryDTO struct {
	ID   *int   `json:"id"`
	Name string `json:"name" binding:"required"`
}

type ProductBrandDTO struct {
	ID      int     `json:"id"`
	Name    string  `json:"name" binding:"required"`
	LogoURL *string `json:"logo_url"`
}

type ProductDTO struct {
	ID          uuid.UUID           `json:"id"`
	Name        string              `json:"name"`
	Price       int64               `json:"price" binding:"required,gt=0"`
	Stock       int                 `json:"stock" binding:"required,gte=0"`
	Description string              `json:"description"`
	Summary     string              `json:"summary"`
	ShopID      uuid.UUID           `json:"shop_id"`
	Size        ProductSizeDTO      `json:"size"`
	Category    ProductCategoryDTO  `json:"category"`
	Condition   ProductConditionDTO `json:"condition"`
	Brand       ProductBrandDTO     `json:"brand"`
	Images      []ProductImageDTO   `json:"images"`
}

type ProductFilterDTO struct {
	Keyword     string  `json:"keyword,omitempty"`      // cari berdasarkan nama produk
	CategoryID  *int    `json:"category_id,omitempty"`  // filter kategori
	BrandID     *int    `json:"brand_id,omitempty"`     // filter brand
	SizeID      *int    `json:"size_id,omitempty"`      // filter ukuran
	ConditionID *int16  `json:"condition_id,omitempty"` // filter kondisi (baru/bekas)
	MinPrice    *int64  `json:"min_price,omitempty"`    // filter harga minimum
	MaxPrice    *int64  `json:"max_price,omitempty"`    // filter harga maksimum
	ShopID      *string `json:"shop_id,omitempty"`      // filter produk milik toko tertentu
}

type PaginatedProductDTO struct {
	Page       int          `json:"page"`
	Size       int          `json:"size"`
	TotalItems int          `json:"total_items"`
	TotalPages int          `json:"total_pages"`
	Items      []ProductDTO `json:"items"`
}

func ToProductDTO(p model.Product) ProductDTO {
	dto := ProductDTO{
		ID:          p.ID,
		Name:        p.Name,
		Price:       p.Price,
		Stock:       p.Stock,
		Description: derefStr(p.Description),
		Summary:     derefStr(p.Summary),
		ShopID:      p.ShopID,
	}

	if p.Category != nil {
		dto.Category = ProductCategoryDTO{
			ID:   &p.Category.ID,
			Name: p.Category.Name,
		}
	}

	if p.Brand != nil {
		dto.Brand = ProductBrandDTO{
			ID:      p.Brand.ID,
			Name:    p.Brand.Name,
			LogoURL: p.Brand.LogoURL,
		}
	}

	if p.Condition != nil {
		dto.Condition = ProductConditionDTO{
			ID:   int(p.Condition.ID),
			Name: p.Condition.Name,
		}
	}

	if p.Size != nil {
		dto.Size = ProductSizeDTO{
			ID:   p.Size.ID,
			Name: p.Size.Name,
		}
	}

	if len(p.Images) > 0 {
		dto.Images = make([]ProductImageDTO, len(p.Images))
		for i, img := range p.Images {
			dto.Images[i] = ProductImageDTO{
				Index: int(img.ImageIndex),
				URL:   img.URL,
			}
		}
	}

	return dto
}

func derefStr(ptr *string) string {
	if ptr != nil {
		return *ptr
	}
	return ""
}
