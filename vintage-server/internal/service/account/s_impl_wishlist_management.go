package service

import (
	"context"
	"log"
	"math"
	common "vintage-server/internal/common/dto"
	"vintage-server/internal/domain/account"
	"vintage-server/internal/model"
	db_error "vintage-server/pkg"

	"github.com/google/uuid"
)

// AddToWishlist implements account.AccountService.
func (s *accountService) AddToWishlist(ctx context.Context, userID uuid.UUID, productID uuid.UUID) error {
	wishlist := model.Wishlist{
		AccountID: userID,
		ProductID: productID,
	}
	log.Printf("Logging %v", wishlist)
	return db_error.HandlePgError(s.store.SaveWishlistItem(ctx, wishlist))
}

// DeactivateUser implements account.AccountService.
func (s *accountService) DeactivateUser(ctx context.Context, userID int64, reason string) error {
	panic("unimplemented")
}

// GetUserProfile implements account.AccountService.
func (s *accountService) GetUserProfile(ctx context.Context, userID int64) (model.Account, error) {
	panic("unimplemented")
}

// GetWishlistByUserID implements account.AccountService.
func (s *accountService) GetWishlistByUserID(
	ctx context.Context,
	req account.WishlistFilter,
) (common.Pagination[account.WishlistItemDetail], error) {

	// Hitung offset untuk pagination
	offset := (req.Page - 1) * req.Size

	// Ambil data wishlist dari repository (model.Wishlist)
	wishlists, total, err := s.store.FindWishlistByAccountID(
		ctx, req.AccountID, req.Keyword, req.Size, offset,
	)
	if err != nil {
		return common.Pagination[account.WishlistItemDetail]{}, err
	}

	// --- Konversi model.Wishlist ke WishlistItemDetail ---
	var items []account.WishlistItemDetail
	for _, w := range wishlists {
		// Ambil detail produk tambahan dari repository/DB
		product, err := s.store.FindProductByID(ctx, w.ProductID)
		if err != nil {
			return common.Pagination[account.WishlistItemDetail]{}, err
		}

		item := account.WishlistItemDetail{
			ProductID:       w.ProductID,
			ProductName:     product.Name,
			ProductPrice:    float32(product.Price),
			ProductImageURL: "http://example.com/image.jpg", // Placeholder, ganti dengan URL gambar sebenarnya
			AddedAt:         w.CreatedAt,
		}
		items = append(items, item)
	}

	// Hitung total halaman
	totalPages := int(math.Ceil(float64(total) / float64(req.Size)))

	// Bungkus hasil ke dalam struct Pagination
	return common.Pagination[account.WishlistItemDetail]{
		Page:       req.Page,
		Size:       req.Size,
		TotalItems: total,
		TotalPages: totalPages,
		Items:      items,
	}, nil
}

// RemoveFromWishlist implements account.AccountService.
func (s *accountService) RemoveFromWishlist(ctx context.Context, userID int64, productID int64) error {
	panic("unimplemented")
}
