package service

import (
	"context"
	"fmt"
	"vintage-server/internal/shop/domain"
	"vintage-server/internal/shop/repository"
	"vintage-server/pkg/uploader"
)

type ShopService interface {
	CreateShop(ctx context.Context, shop domain.Shop) (domain.Shop, error)
}

type shopSvc struct {
	store    repository.ShopStore
	uploader uploader.Uploader
}

// CreateShop implements ShopService.
func (s *shopSvc) CreateShop(ctx context.Context, shop domain.Shop) (domain.Shop, error) {
	var created domain.Shop

	// Jalankan semua query di dalam transaksi
	err := s.store.ExecTx(ctx, func(txStore repository.ShopStore) error {
		shopRepo := txStore.GetShopRepo()

		// 0️⃣ Cek apakah shop sudah ada untuk akun ini
		exists, err := shopRepo.IsShopExists(ctx, shop.AccountID)
		if err != nil {
			return fmt.Errorf("failed to check shop existence: %w", err)
		}
		if exists {
			return fmt.Errorf("user already has a shop")
		}

		// 1️⃣ Buat shop baru
		newShop, err := shopRepo.CreateShop(ctx, shop)
		if err != nil {
			return fmt.Errorf("failed to create shop: %w", err)
		}

		// 2️⃣ Ambil role ID untuk 'seller'
		roleID, err := shopRepo.GetRoleIDByName(ctx, "seller")
		if err != nil {
			return fmt.Errorf("failed to get seller role ID: %w", err)
		}

		// 3️⃣ Tambahkan role seller ke akun
		if err := shopRepo.AddSellerRole(ctx, shop.AccountID, roleID); err != nil {
			return fmt.Errorf("failed to assign seller role: %w", err)
		}

		// 4️⃣ Simpan hasil shop yang sudah dibuat
		created = newShop
		return nil
	})

	if err != nil {
		return domain.Shop{}, fmt.Errorf("service_error: %w", err)
	}

	return created, nil
}

func NewShopSvc(
	store repository.ShopStore,
	uploader uploader.Uploader,
) ShopService {
	return &shopSvc{
		store:    store,
		uploader: uploader,
	}
}
