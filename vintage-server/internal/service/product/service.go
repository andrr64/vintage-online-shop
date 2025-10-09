package product

import (
	"context"
	"database/sql"
	"errors"
	"log"
	product "vintage-server/internal/domain"
	"vintage-server/internal/model"
	"vintage-server/internal/repository"
	"vintage-server/pkg/apperror"
	"vintage-server/pkg/auth"
	"vintage-server/pkg/uploader"

	"github.com/google/uuid"
)

// service is a struct that will implement the Service interface from domain.go
type service struct {
	store    repository.ProductStore // Store sudah embed ProductRepository interface
	jwt      auth.JWTService
	uploader uploader.Uploader
}

func NewService(store repository.ProductStore, jwt auth.JWTService, uploader uploader.Uploader) product.ProductService {
	return &service{
		store:    store,
		jwt:      jwt,
		uploader: uploader,
	}
}

// -- CATEGORY MANAGEMENT --
func (s *service) CreateCategory(ctx context.Context, req product.ProductCategory) error {
	data := model.ProductCategory{
		Name: req.Name,
	}
	// Ganti s.repo menjadi s.store
	_, err := s.store.CreateCategory(ctx, data)
	if err != nil {
		return err
	}
	return nil
}

// FindAllCategories adalah operasi baca, tidak memerlukan transaksi.
func (s *service) FindAllCategories(ctx context.Context) ([]product.ProductCategory, error) {
	// Ganti s.repo menjadi s.store
	categoriesFromRepo, err := s.store.FindAllCategories(ctx)
	if err != nil {
		log.Printf("Error finding all categories from repo: %v", err)
		return nil, apperror.New(apperror.ErrCodeInternal, "an internal error occurred")
	}

	var result []product.ProductCategory
	for _, category := range categoriesFromRepo {
		result = append(result, product.ProductCategory{
			ID:   &category.ID,
			Name: category.Name,
		})
	}
	return result, nil
}

// FindById adalah operasi baca, tidak memerlukan transaksi.
func (s *service) FindById(ctx context.Context, id int) (product.ProductCategory, error) {
	// Ganti s.repo menjadi s.store
	res, err := s.store.FindCategoryById(ctx, id)
	if err != nil {
		// Di sini kita bisa menambahkan mapping error not found yang lebih baik
		if errors.Is(err, sql.ErrNoRows) {
			return product.ProductCategory{}, apperror.New(apperror.ErrCodeNotFound, "category not found")
		}
		return product.ProductCategory{}, err
	}
	return product.ProductCategory{
		ID:   &res.ID,
		Name: res.Name,
	}, nil
}

func (s *service) UpdateCategory(ctx context.Context, req product.ProductCategory) error {
	err := s.store.ExecTx(ctx, func(repoInTx product.ProductRepository) error {
		id := *req.ID
		// 1. Baca data LAMA di dalam transaksi
		old, errTx := repoInTx.FindCategoryById(ctx, id)
		if errTx != nil {
			if errors.Is(errTx, sql.ErrNoRows) {
				return apperror.New(apperror.ErrCodeNotFound, "category not found")
			}
			return errTx
		}
		// 2. Modifikasi data
		old.Name = req.Name

		// 3. Tulis data BARU di dalam transaksi yang sama
		errTx = repoInTx.UpdateCategory(ctx, old)
		if errTx != nil {
			return errTx
		}

		return nil // Sukses, transaksi akan di-commit
	})
	return err
}

func (s *service) DeleteCategory(ctx context.Context, categoryID int) error {
	err := s.store.ExecTx(ctx, func(repoInTx product.ProductRepository) error {
		_, errTx := repoInTx.FindCategoryById(ctx, categoryID)
		if errTx != nil {
			if errors.Is(errTx, sql.ErrNoRows) {
				return apperror.New(apperror.ErrCodeNotFound, "category not found")
			}
			return errTx
		}
		// 2. Cek apakah kategori masih digunakan oleh produk (di dalam transaksi)
		count, errTx := repoInTx.CountProductsByCategory(ctx, categoryID)
		if errTx != nil {
			return apperror.New(apperror.ErrCodeInternal, "failed to check category usage")
		}

		if count > 0 {
			// Jika masih dipakai, kembalikan error. Transaksi akan di-rollback.
			return apperror.New(apperror.ErrCodeConflict, "cannot delete category that is still in use by products")
		}

		// 3. Jika aman, hapus kategori (di dalam transaksi yang sama)
		errTx = repoInTx.DeleteCategory(ctx, categoryID)
		if errTx != nil {
			return errTx
		}

		return nil // Sukses, transaksi akan di-commit
	})
	return err
}

// -- BRAND MANAGEMENt--
func (s *service) CreateBrand(ctx context.Context, req product.CreateBrandRequest) (model.Brand, error) {
	logoURL, err := s.uploader.UploadBrandLogo(ctx, req.File)
	if err != nil {
		return model.Brand{}, apperror.New(apperror.ErrCodeInternal, "failed to upload logo")
	}
	// 2. Siapkan data untuk disimpan ke database.
	brand := model.Brand{
		Name:    req.Name,
		LogoURL: &logoURL,
	}
	// 3. Simpan ke database.
	// Menggunakan s.store, bukan s.repo.
	createdBrand, err := s.store.CreateBrand(ctx, brand)
	if err != nil {
		// Jika penyimpanan DB GAGAL, kita harus membatalkan upload (rollback manual).
		log.Printf("Database insertion failed for brand '%s'. Deleting uploaded logo: %s", req.Name, logoURL)
		if delErr := s.uploader.DeleteByURL(ctx, logoURL); delErr != nil {
			log.Printf("CRITICAL: failed to delete uploaded logo during rollback: %v", delErr)
		}
		return model.Brand{}, err // Kembalikan error asli dari database.
	}
	return createdBrand, nil
}

func (s *service) FindAllBrands(ctx context.Context) ([]model.Brand, error) {
	return s.store.FindAllBrands(ctx)
}

func (s *service) FindBrandByID(ctx context.Context, id int) (model.Brand, error) {
	brand, err := s.store.FindBrandByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Brand{}, apperror.New(apperror.ErrCodeNotFound, "brand not found")
		}
		return model.Brand{}, err
	}
	return brand, nil
}

func (s *service) UpdateBrand(ctx context.Context, id int, req product.UpdateBrandRequest) error {
	var oldLogoURL *string
	var newLogoURL string

	// 1. Jika ada file baru, upload terlebih dahulu.
	if req.File != nil {
		var errUpload error
		newLogoURL, errUpload = s.uploader.UploadBrandLogo(ctx, req.File)
		if errUpload != nil {
			return apperror.New(apperror.ErrCodeInternal, "failed to upload new logo")
		}
	}

	// 2. Lakukan operasi database di dalam transaksi yang aman.
	err := s.store.ExecTx(ctx, func(repoInTx product.ProductRepository) error {
		// 2a. Ambil data brand yang ada di dalam transaksi.
		existingBrand, errTx := repoInTx.FindBrandByID(ctx, id)
		if errTx != nil {
			if errors.Is(errTx, sql.ErrNoRows) {
				return apperror.New(apperror.ErrCodeNotFound, "brand not found")
			}
			return errTx
		}

		// Simpan URL logo lama untuk dihapus nanti HANYA JIKA transaksi berhasil.
		oldLogoURL = existingBrand.LogoURL

		// 2b. Update field-fieldnya.
		existingBrand.Name = req.Name
		if newLogoURL != "" { // Jika ada logo baru yang di-upload
			existingBrand.LogoURL = &newLogoURL
		}

		// 2c. Simpan perubahan ke DB di dalam transaksi.
		return repoInTx.UpdateBrand(ctx, existingBrand)
	})

	if err != nil {
		// Jika transaksi GAGAL, kita harus membatalkan upload logo baru (jika ada).
		if newLogoURL != "" {
			log.Printf("Database transaction failed for updating brand %d. Deleting newly uploaded logo: %s", id, newLogoURL)
			s.uploader.DeleteByURL(context.Background(), newLogoURL)
		}
		return err
	}

	// 4. Jika transaksi SUKSES, hapus logo lama dari cloud (jika ada).
	if oldLogoURL != nil && newLogoURL != "" {
		log.Printf("Database transaction successful. Deleting old logo: %s", *oldLogoURL)
		s.uploader.DeleteByURL(context.Background(), *oldLogoURL)
	}
	return nil
}

func (s *service) DeleteBrand(ctx context.Context, id int) error {
	var brandToDelete model.Brand

	// Gunakan s.store.ExecTx untuk membuat "safe scope".
	err := s.store.ExecTx(ctx, func(repoInTx product.ProductRepository) error {
		// 1. Ambil data brand untuk mendapatkan URL logo dan memastikan brand ada.
		var errTx error
		brandToDelete, errTx = repoInTx.FindBrandByID(ctx, id)
		if errTx != nil {
			if errors.Is(errTx, sql.ErrNoRows) {
				return apperror.New(apperror.ErrCodeNotFound, "brand not found")
			}
			return errTx
		}

		// 2. Cek apakah brand masih digunakan (di dalam transaksi).
		count, errTx := repoInTx.CountProductsByBrand(ctx, id)
		if errTx != nil {
			return apperror.New(apperror.ErrCodeInternal, "failed to check brand usage")
		}
		if count > 0 {
			return apperror.New(apperror.ErrCodeConflict, "cannot delete brand that is still in use by products")
		}

		// 3. Hapus brand dari database (di dalam transaksi).
		return repoInTx.DeleteBrand(ctx, id)
	})
	// Jika transaksi GAGAL, hentikan proses.
	if err != nil {
		return err
	}
	// Jika transaksi SUKSES, hapus logo dari cloud.
	if brandToDelete.LogoURL != nil {
		log.Printf("Database deletion successful for brand %d. Deleting logo from cloud: %s", id, *brandToDelete.LogoURL)
		s.uploader.DeleteByURL(context.Background(), *brandToDelete.LogoURL)
	}
	return nil
}

// -- PRODUCT CONDITION MANAGEMENT --
func (s *service) CreateCondition(ctx context.Context, req product.ProductConditionRequest) (model.ProductCondition, error) {
	result, err := s.store.CreateCondition(ctx, model.ProductCondition{Name: req.Name})
	if err != nil {
		return model.ProductCondition{}, err
	}
	return result, err
}

func (s *service) FindAllConditions(ctx context.Context) ([]model.ProductCondition, error) {
	return s.store.FindAllConditions(ctx)
}

func (s *service) FindConditionByID(ctx context.Context, id int16) (model.ProductCondition, error) {
	condition, err := s.store.FindConditionByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.ProductCondition{}, apperror.New(apperror.ErrCodeNotFound, "Product condition not found")
		}
		return model.ProductCondition{}, err
	}
	return condition, err
}

func (s *service) UpdateCondition(ctx context.Context, id int16, req product.ProductConditionRequest) (model.ProductCondition, error) {
	var updatedCondition model.ProductCondition

	// Gunakan s.store.ExecTx untuk membuat "safe scope".
	err := s.store.ExecTx(ctx, func(repoInTx product.ProductRepository) error {
		// 1. Baca data LAMA di dalam transaksi.
		conditionToUpdate, errTx := repoInTx.FindConditionByID(ctx, id)
		if errTx != nil {
			if errors.Is(errTx, sql.ErrNoRows) {
				return apperror.New(apperror.ErrCodeNotFound, "product condition not found")
			}
			return errTx
		}

		// 2. Modifikasi data.
		conditionToUpdate.Name = req.Name

		// 3. Tulis data BARU di dalam transaksi yang sama.
		updatedCondition, errTx = repoInTx.UpdateCondition(ctx, conditionToUpdate)
		if errTx != nil {
			return errTx
		}

		return nil // Sukses, transaksi akan di-commit.
	})

	return updatedCondition, err // Kembalikan hasil dan error dari transaksi.
}

func (s *service) DeleteCondition(ctx context.Context, id int16) error {
	// 1. Mulai transaksi
	err := s.store.ExecTx(ctx, func(r product.ProductRepository) error {

		// 2. Cek apakah data exists atau tidak
		_, errTx := r.FindConditionByID(ctx, id)
		if errTx != nil {
			if errors.Is(errTx, sql.ErrNoRows) {
				return apperror.New(apperror.ErrCodeNotFound, "product condition not found")
			}
			return errTx
		}

		// 3. JUMLAH PRODUK YANG TERIKAT DENGAN CONDITION HARUS 0
		count, errTx := r.CountProductsByCondition(ctx, id)
		if errTx != nil {
			return apperror.New(apperror.ErrCodeInternal, "failed to check condition usage")
		}
		if count > 0 {
			// Jika masih dipakai, kembalikan error. Transaksi akan otomatis di-rollback.
			return apperror.New(apperror.ErrCodeConflict, "cannot delete condition that is still in use by products")
		}

		return r.DeleteCondition(ctx, id)
	})
	return err
}

// internal/service/product_service.go

func (s *service) CreateProduct(ctx context.Context, accountID uuid.UUID, request product.CreateProductRequest) {
	// 1. VALIDASI BISNIS: Pastikan akun yang login memiliki toko yang terdaftar.
	// Kita menggunakan s.store karena ini adalah operasi baca.
	shop, err := s.store.FindShopByAccountID(ctx, accountID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return //TODO
		}
		// Error database lainnya
		return //TODO
	}

	// 2. INTERAKSI LAYANAN EKSTERNAL: Upload semua gambar ke cloud storage.
	var uploadedImageURLs []string

	// Fungsi cleanup ini adalah "rollback manual" untuk cloud storage.
	cleanupCloudImages := func() {
		log.Println("Cleanup triggered: deleting product images from cloud storage...")
		cleanupCtx := context.Background()
		for _, url := range uploadedImageURLs {
			if delErr := s.uploader.DeleteByURL(cleanupCtx, url); delErr != nil {
				log.Printf("CRITICAL: Failed to delete image %s during cleanup: %v", url, delErr)
			}
		}
	}

	// Loop untuk upload file satu per satu.
	for _, fileHeader := range request.Images {
		file, err := fileHeader.Open()
		if err != nil {
			cleanupCloudImages()
			return //TODO
		}

		url, err := s.uploader.UploadProductImage(ctx, file)
		file.Close() // Tutup file sesegera mungkin.

		if err != nil {
			cleanupCloudImages() // Jika upload gagal, bersihkan semua yang sudah berhasil.
			return               //TODO
		}
		uploadedImageURLs = append(uploadedImageURLs, url)
	}

	// 3. OPERASI DATABASE: Simpan produk dan gambar dalam satu transaksi atomik.
	var createdProduct model.Product
	err = s.store.ExecTx(ctx, func(repoInTx product.ProductRepository) error {
		// Semua kode di dalam blok ini dijamin aman dari race condition.

		// 3a. Siapkan data produk utama.
		productData := model.Product{
			ShopID:      shop.ID,
			Name:        request.Name,
			CategoryID:  request.CategoryID,
			ConditionID: request.ConditionID,
			Price:       request.Price,
			Stock:       request.Stock,
			Description: &request.Description,
			Summary:     &request.Summary,
			BrandID:     request.BrandID,
			SizeID:      request.SizeID,
		}

		// 3b. Simpan produk utama ke DB.
		var errTx error
		createdProduct, errTx = repoInTx.CreateProduct(ctx, productData)
		if errTx != nil {
			return errTx // Error di sini akan memicu Rollback oleh Store.
		}

		// 3c. Siapkan data gambar-gambar produk dengan ID produk yang baru dibuat.
		var imageModels []model.ProductImage
		for i, imgURL := range uploadedImageURLs {
			imageModels = append(imageModels, model.ProductImage{
				ProductID:  createdProduct.ID,
				ImageIndex: int16(i + 1), // Index gambar mulai dari 1
				URL:        imgURL,
			})
		}

		// 3d. Simpan semua data gambar ke DB (bulk insert).
		if errTx = repoInTx.CreateProductImages(ctx, imageModels); errTx != nil {
			return errTx // Error di sini juga akan memicu Rollback.
		}

		return nil // Sukses! Store akan otomatis Commit transaksi ini.
	})

	// 4. PENANGANAN HASIL TRANSAKSI
	if err != nil {
		// Jika ExecTx mengembalikan error, berarti transaksi database sudah di-rollback.
		// Tugas kita sekarang adalah melakukan rollback manual untuk file di cloud.
		cleanupCloudImages()
		return //TODO
	}

	// // 5. SIAPKAN RESPONSE SUKSES
	// imageResponses := []product.ProductImageResponse{}
	// for i, url := range uploadedImageURLs {
	// 	imageResponses = append(imageResponses, ProductImageResponse{
	// 		URL:        url,
	// 		ImageIndex: int16(i + 1),
	// 	})
	// }

	// response := product.ProductResponse{
	// 	ID:          createdProduct.ID,
	// 	ShopID:      createdProduct.ShopID,
	// 	Name:        createdProduct.Name,
	// 	Price:       createdProduct.Price,
	// 	Stock:       createdProduct.Stock,
	// 	Summary:     createdProduct.Summary,
	// 	Description: createdProduct.Description,
	// 	Images:      imageResponses,
	// }

	// return response, nil
	return
}
