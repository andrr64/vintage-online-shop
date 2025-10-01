package product

import (
	"context"
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

func (s *service) DeleteCategory(ctx context.Context, id int) error {
	// Hitung jumlah produk yang pakai kategori ini
	count, err := s.repo.CountProductsByCategory(ctx, id)
	if err != nil {
		return apperror.New(apperror.ErrCodeInternal, "gagal mengecek produk terkait kategori")
	}

	if count > 0 {
		return apperror.New(apperror.ErrCodeConflict, "kategori tidak bisa dihapus karena masih digunakan oleh produk")
	}

	// Kalau aman, hapus kategori
	err = s.repo.DeleteCategory(ctx, id)
	if err != nil {
		return apperror.New(apperror.ErrCodeInternal, "gagal menghapus kategori")
	}

	return nil
}
