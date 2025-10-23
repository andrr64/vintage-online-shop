package service

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"vintage-server/internal/product/domain"
	"vintage-server/internal/product/repository"
	"vintage-server/pkg/uploader"
)

type BrandService interface {
	CreateBrand(c context.Context, brand domain.Brand, logo multipart.File) (domain.Brand, error)
	UpdateBrand(c context.Context, brand domain.Brand, logo *multipart.File) (domain.Brand, error)
	DeleteBrand(c context.Context, brandID int) error
	ReadBrands(c context.Context, brandID *int) ([]domain.Brand, error)
}

type brandServiceImpl struct {
	store    repository.ProductStore
	uploader uploader.Uploader
}

// CreateBrand implements BrandService.
func (b *brandServiceImpl) CreateBrand(c context.Context, brand domain.Brand, logo multipart.File) (domain.Brand, error) {
	logourl, err := b.uploader.UploadBrandLogo(c, logo)
	if err != nil {

		return domain.Brand{}, fmt.Errorf("error: gagal mengupload file")
	}
	brand.LogoURL = &logourl
	createdBrand, err := b.store.GetBrandRepository().CreateBrand(c, brand)

	if err != nil {
		if delErr := b.uploader.DeleteByURL(c, logourl); delErr != nil {
			log.Printf("CRITICAL: failed to delete uploaded logo during rollback: %v", delErr)
		}
		return domain.Brand{}, err
	}
	return createdBrand, nil
}

// UpdateBrand implements BrandService.
// UpdateBrand implements BrandService.
func (b *brandServiceImpl) UpdateBrand(c context.Context, brand domain.Brand, logo *multipart.File) (domain.Brand, error) {
	// Ambil data brand lama dari database
	brands, err := b.store.GetBrandRepository().ReadBrands(c, &brand.ID)
	if err != nil {
		return domain.Brand{}, fmt.Errorf("error: gagal mengambil data brand")
	}
	if len(brands) == 0 {
		return domain.Brand{}, fmt.Errorf("error: brand tidak ditemukan")
	}
	oldBrand := brands[0]

	var newLogoURL *string

	// Jika ada logo baru diupload
	if logo != nil {
		uploadedURL, err := b.uploader.UploadBrandLogo(c, *logo)
		if err != nil {
			return domain.Brand{}, fmt.Errorf("error: gagal mengupload logo baru")
		}
		newLogoURL = &uploadedURL
		brand.LogoURL = newLogoURL
	} else {
		// Kalau tidak ada upload baru, tetap pakai logo lama
		brand.LogoURL = oldBrand.LogoURL
	}

	// Update brand di database
	updatedBrand, err := b.store.GetBrandRepository().UpdateBrand(c, brand)
	if err != nil {
		// rollback logo baru kalau gagal update
		if newLogoURL != nil {
			_ = b.uploader.DeleteByURL(c, *newLogoURL)
		}
		return domain.Brand{}, fmt.Errorf("error: gagal memperbarui brand")
	}

	// Jika logo baru berhasil diupload dan update sukses → hapus logo lama
	if newLogoURL != nil && oldBrand.LogoURL != nil {
		if err := b.uploader.DeleteByURL(context.Background(), *oldBrand.LogoURL); err != nil {
			log.Printf("warning: gagal menghapus logo lama: %v", err)
		}
	}

	return updatedBrand, nil
}

// DeleteBrand implements BrandService.
func (b *brandServiceImpl) DeleteBrand(c context.Context, brandID int) error {
	var brandToDelete domain.Brand

	err := b.store.ExecTx(c, func(txRepo repository.ProductStore) error {
		brands, errTx := txRepo.GetBrandRepository().ReadBrands(c, &brandID)
		if errTx != nil {
			return fmt.Errorf("failed to read brand: %w", errTx)
		}

		// ✅ Cek dulu kalau hasilnya kosong
		if len(brands) == 0 {
			return fmt.Errorf("error: brand not found")
		}

		brandToDelete = brands[0]

		count, errTx := txRepo.GetBrandRepository().CountProducts(c, brandToDelete)
		if errTx != nil {
			return fmt.Errorf("failed to check brand usage")
		}
		if count > 0 {
			return fmt.Errorf("cannot delete brand that is still in use by products")
		}

		return txRepo.GetBrandRepository().DeleteBrand(c, brandID)
	})

	if err != nil {
		return err
	}

	// ✅ Aman: brandToDelete hanya digunakan kalau berhasil
	if brandToDelete.LogoURL != nil {
		log.Printf("Database deletion successful for brand %d. Deleting logo from cloud: %s", brandID, *brandToDelete.LogoURL)
		b.uploader.DeleteByURL(context.Background(), *brandToDelete.LogoURL)
	}

	return nil
}

// ReadBrands implements BrandService.
func (b *brandServiceImpl) ReadBrands(c context.Context, brandID *int) ([]domain.Brand, error) {
	var brands []domain.Brand = make([]domain.Brand, 0)
	data, err := b.store.GetBrandRepository().ReadBrands(c, brandID)
	if err != nil {
		return brands, err
	}
	brands = data
	return brands, nil
}

func NewBrandService(store repository.ProductStore, upl uploader.Uploader) BrandService {
	return &brandServiceImpl{
		store:    store,
		uploader: upl,
	}
}
