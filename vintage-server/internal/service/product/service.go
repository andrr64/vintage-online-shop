package product

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"vintage-server/internal/model"
	"vintage-server/pkg/apperror"
	"vintage-server/pkg/auth"
	"vintage-server/pkg/uploader"
)

// service is a struct that will implement the Service interface from domain.go
type service struct {
	repo     Repository
	jwt      auth.JWTService // <-- Gunakan interface
	uploader uploader.Uploader
}

func NewService(repo Repository, jwt auth.JWTService, uploader uploader.Uploader) Service {
	return &service{
		repo:     repo,
		jwt:      jwt, // <-- Terima dari parameter
		uploader: uploader,
	}
}

// -- CATEGORY MANAGEMENT --
func (s *service) CreateCategory(ctx context.Context, req ProductCategory) error {
	data := model.ProductCategory{
		Name: req.Name,
	}
	err := s.repo.CreateCategory(ctx, data)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) FindAllCategories(ctx context.Context) ([]ProductCategory, error) {
	categoriesFromRepo, err := s.repo.FindAllCategories(ctx)
	if err != nil {
		// Cukup tangani error non-ErrNoRows.
		// sqlx.Select sudah mengembalikan slice kosong jika tidak ada hasil, bukan error.
		log.Printf("Error finding all categories from repo: %v", err)
		return nil, apperror.New(apperror.ErrCodeInternal, "an internal error occurred")
	}

	// Proses Mapping (sudah benar)
	var result []ProductCategory
	for _, category := range categoriesFromRepo {
		result = append(result, ProductCategory{
			ID:   &category.ID,
			Name: category.Name,
		})
	}

	return result, nil
}

func (s *service) FindById(ctx context.Context, id int) (ProductCategory, error) {
	res, err := s.repo.FindById(ctx, id)
	if err != nil {
		return ProductCategory{}, err
	}
	return ProductCategory{
		ID:   &res.ID,
		Name: res.Name,
	}, nil
}

func (s *service) UpdateCategory(ctx context.Context, req ProductCategory) error {
	id := *req.ID
	old, err := s.repo.FindById(ctx, id)
	if err != nil {
		return err
	}
	old.Name = req.Name
	err = s.repo.UpdateCategory(ctx, old)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) DeleteCategory(ctx context.Context, categoryID int) error {
	// 1. Cek apakah kategori ini masih digunakan oleh produk lain.
	count, err := s.repo.CountProductsByCategory(ctx, categoryID)
	if err != nil {
		// Jika errornya bukan 'not found', maka ini error internal
		if !errors.Is(err, sql.ErrNoRows) {
			log.Printf("Error counting products by category: %v", err)
			return apperror.New(apperror.ErrCodeInternal, "failed to check category usage")
		}
		// Jika errornya 'not found' dari repo lain (misal FindById), lanjutkan saja
		// karena Delete di repo akan menangani kasus 'not found' juga.
	}

	// 2. Jika count > 0, artinya kategori masih dipakai. Kembalikan error konflik.
	if count > 0 {
		return apperror.New(apperror.ErrCodeConflict, "cannot delete category that is still in use by products")
	}

	// 3. Jika aman (count == 0), panggil repository untuk menghapus kategori.
	err = s.repo.DeleteCategory(ctx, categoryID)
	if err != nil {
		// Kembalikan error dari repository (bisa jadi NotFound atau Internal)
		return err
	}

	return nil
}

// -- BRAND MANAGEMENt--

func (s *service) CreateBrand(ctx context.Context, req CreateBrandRequest) (model.Brand, error) {
	// 1. Upload file ke cloud storage
	logoURL, err := s.uploader.UploadBrandLogo(ctx, req.File)
	if err != nil {
		return model.Brand{}, apperror.New(apperror.ErrCodeInternal, "failed to upload logo")
	}

	// 2. Siapkan data untuk disimpan ke database
	brand := model.Brand{
		Name:    req.Name,
		LogoURL: &logoURL, // Kirim URL sebagai pointer string
	}

	// 3. Panggil repository untuk menyimpan
	return s.repo.CreateBrand(ctx, brand)
}

func (s *service) FindAllBrands(ctx context.Context) ([]model.Brand, error) {
	return s.repo.FindAllBrands(ctx)
}

func (s *service) FindBrandByID(ctx context.Context, id int) (model.Brand, error) {
	return s.repo.FindBrandByID(ctx, id)
}

func (s *service) UpdateBrand(ctx context.Context, id int, req UpdateBrandRequest) error {
	// 1. Ambil data brand yang ada
	existingBrand, err := s.repo.FindBrandByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperror.New(apperror.ErrCodeNotFound, "Brand not found, bro...")
		}
		return err
	}

	// 2. Update nama
	existingBrand.Name = req.Name

	// 3. Jika ada file baru, upload dan ganti URL logo
	print(req.File != nil)
	if req.File != nil {
		if existingBrand.LogoURL != nil {
			s.uploader.DeleteByURL(ctx, *existingBrand.LogoURL) // Asumsi ada method Delete
		}
		newLogoURL, err := s.uploader.UploadBrandLogo(ctx, req.File)
		if err != nil {
			return apperror.New(apperror.ErrCodeInternal, "failed to upload new logo")
		}
		existingBrand.LogoURL = &newLogoURL
	}

	return s.repo.UpdateBrand(ctx, existingBrand)
}

func (s *service) DeleteBrand(ctx context.Context, id int) error {
	// LOGIKA BISNIS PENTING: Cek apakah brand masih digunakan oleh produk
	count, err := s.repo.CountProductsByBrand(ctx, id)
	if err != nil {
		log.Printf("Error counting products by brand: %v", err)
		return apperror.New(apperror.ErrCodeInternal, "failed to check brand usage")
	}
	if count > 0 {
		return apperror.New(apperror.ErrCodeConflict, "cannot delete brand that is still in use by products")
	}

	// Jika aman, lanjutkan hapus
	return s.repo.DeleteBrand(ctx, id)
}
