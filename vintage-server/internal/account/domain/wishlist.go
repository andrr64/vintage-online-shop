package domain

import (
	"slices"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Wishlist struct {
	ID         int64
	AccountID  uuid.UUID
	ProductIDs []uuid.UUID
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (w *Wishlist) HasProduct(productID uuid.UUID) bool {
	return slices.Contains(w.ProductIDs, productID)
}

func (w *Wishlist) Clear() {
	w.ProductIDs = []uuid.UUID{}
	w.UpdatedAt = time.Now()
}

func (w *Wishlist) Count() int {
	return len(w.ProductIDs)
}

func (w *Wishlist) IsEmpty() bool {
	return len(w.ProductIDs) == 0
}

func (w *Wishlist) CanAddMore(maxItems int) bool {
	return len(w.ProductIDs) < maxItems
}

// Logika bisnis: tambahkan produk ke wishlist
func (w *Wishlist) AddProduct(productID uuid.UUID) error {
	for _, pid := range w.ProductIDs {
		if pid == productID {
			return fmt.Errorf("produk sudah ada di wishlist")
		}
	}
	w.ProductIDs = append(w.ProductIDs, productID)
	w.UpdatedAt = time.Now()
	return nil
}

// Hapus produk
func (w *Wishlist) RemoveProduct(productID uuid.UUID) error {
	for i, pid := range w.ProductIDs {
		if pid == productID {
			w.ProductIDs = append(w.ProductIDs[:i], w.ProductIDs[i+1:]...)
			w.UpdatedAt = time.Now()
			return nil
		}
	}
	return fmt.Errorf("produk tidak ditemukan di wishlist")
}

func NewWishlist(accountID uuid.UUID, productIDs []uuid.UUID) *Wishlist {
	return &Wishlist{
		AccountID:  accountID,
		ProductIDs: productIDs,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}
