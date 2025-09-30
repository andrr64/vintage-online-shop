package product

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"vintage-server/pkg/apperror"
	"vintage-server/pkg/auth"
	"vintage-server/pkg/uploader"
)

// service is a struct that will implement the Service interface from domain.go
type service struct {
	repo     Repository
	jwt      *auth.JWTService
	uploader uploader.Uploader
}

// NewService is the constructor for the service
func NewService(repo Repository, jwtSecret string, uploader uploader.Uploader) Service {
	return &service{
		repo:     repo,
		jwt:      auth.NewJWTService(jwtSecret),
		uploader: uploader,
	}
}

func (s *service) CreateCategory(ctx context.Context, req ProductCategory) error {
	err := s.repo.CreateCategory(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) FindAllCategories(ctx context.Context) ([]ProductCategory, error) {
    // Panggil repository. Gunakan '=' bukan ':=' karena err sudah dideklarasikan di scope fungsi.
    // Kita langsung gunakan variabel baru untuk hasil dari repo agar tidak bingung.
    categoriesFromRepo, err := s.repo.FindAllCategories(ctx)
    if err != nil {
        // Jika tidak ada baris, ini bukan error. Kembalikan slice kosong.
        if errors.Is(err, sql.ErrNoRows) {
            return []ProductCategory{}, nil // <-- Ini cara return slice kosong (best practice)
        }

        // Untuk error database lainnya, catat dan kembalikan error internal.
        // Ganti log message agar lebih sesuai.
        log.Printf("Error finding all categories from repo: %v", err)
        return nil, apperror.New(apperror.ErrCodeInternal, "an internal error occurred")
    }
    // pastikan kita mengembalikan slice kosong yang non-nil.
    if categoriesFromRepo == nil {
        return []ProductCategory{}, nil
    }

    // Jika tidak ada hasil dari repo, kembalikan slice kosong.
    if len(categoriesFromRepo) == 0 {
        return []ProductCategory{}, nil
    }
    
    // --- Proses Mapping dari model.ProductCategory ke service.ProductCategory ---
    // Buat slice dengan tipe return yang benar.
    var result []ProductCategory 
    
    // Lakukan iterasi dan konversi tipe.
    for _, category := range categoriesFromRepo {
        result = append(result, ProductCategory{
            // Asumsi field-nya sama, sesuaikan jika berbeda
            ID:   category.ID,
            Name: category.Name,
            // ... field lainnya
        })
    }

    // Kembalikan hasil yang sudah dipetakan, BUKAN nil.
    return result, nil
}

func (s *service) FindById(ctx context.Context, id int) (ProductCategory, error ){
	res, err := s.repo.FindById(ctx, id)
	if err != nil {
		return ProductCategory{}, err
	}
	return ProductCategory{
		ID: res.ID,
		Name: res.Name,
	}, nil
}