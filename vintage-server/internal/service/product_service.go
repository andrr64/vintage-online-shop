package service

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"
	product "vintage-server/internal/domain"
	"vintage-server/internal/model"
	"vintage-server/internal/repository"
	db_error "vintage-server/pkg"
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

// FindProductByID implements product.ProductService.
func (s *service) FindProductByID(ctx context.Context, productID uuid.UUID) (product.ProductDTO, error) {
	data, err := s.store.FindProductByID(ctx, productID)
	if err != nil {
		return product.ProductDTO{}, db_error.HandlePgError(err)
	}
	return product.ToProductDTO(data), nil
}

// CreateProductSize implements product.ProductService.
func (s *service) CreateProductSize(ctx context.Context, request product.ProductConditionRequest) (product.ProductSizeDTO, error) {
	data := model.ProductSize{
		Name: request.Name,
	}

	saved, err := s.store.CreateProductSize(ctx, data)

	if err != nil {
		return product.ProductSizeDTO{}, db_error.HandlePgError(err)
	}
	return product.ProductSizeDTO{
		Name: saved.Name,
		ID:   saved.ID,
	}, nil
}

// -- CATEGORY MANAGEMENT --
func (s *service) CreateCategory(ctx context.Context, req product.ProductCategoryDTO) error {
	data := model.ProductCategory{
		Name: req.Name,
	}
	// Ganti s.repo menjadi s.store
	_, err := s.store.CreateCategory(ctx, data)
	if err != nil {
		return db_error.HandlePgError(err)
	}
	return nil
}

// FindAllCategories adalah operasi baca, tidak memerlukan transaksi.
func (s *service) FindAllCategories(ctx context.Context) ([]product.ProductCategoryDTO, error) {
	// Ganti s.repo menjadi s.store
	categoriesFromRepo, err := s.store.FindAllCategories(ctx)
	if err != nil {
		log.Printf("Error finding all categories from repo: %v", err)
		return nil, apperror.New(apperror.ErrCodeInternal, "an internal error occurred")
	}

	var result []product.ProductCategoryDTO
	for _, category := range categoriesFromRepo {
		result = append(result, product.ProductCategoryDTO{
			ID:   &category.ID,
			Name: category.Name,
		})
	}
	return result, nil
}

// FindById adalah operasi baca, tidak memerlukan transaksi.
func (s *service) FindById(ctx context.Context, id int) (product.ProductCategoryDTO, error) {
	// Ganti s.repo menjadi s.store
	res, err := s.store.FindCategoryById(ctx, id)
	if err != nil {
		// Di sini kita bisa menambahkan mapping error not found yang lebih baik
		if errors.Is(err, sql.ErrNoRows) {
			return product.ProductCategoryDTO{}, apperror.New(apperror.ErrCodeNotFound, "category not found")
		}
		return product.ProductCategoryDTO{}, err
	}
	return product.ProductCategoryDTO{
		ID:   &res.ID,
		Name: res.Name,
	}, nil
}

func (s *service) UpdateCategory(ctx context.Context, req product.ProductCategoryDTO) error {
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

func (s *service) FindAllBrands(ctx context.Context) ([]product.ProductBrandDTO, error) {
	// Inisialisasi slice kosong biar gak nil
	brands := make([]product.ProductBrandDTO, 0)

	// Ambil data mentah dari store
	data, err := s.store.FindAllBrands(ctx)
	if err != nil {
		// Tetap return slice kosong, bukan nil
		return brands, db_error.HandlePgError(err)
	}

	// Mapping dari data store ke DTO
	for _, d := range data {
		brands = append(brands, product.ProductBrandDTO{
			ID:      d.ID,
			Name:    d.Name,
			LogoURL: d.LogoURL,
		})
	}

	return brands, nil
}

func (s *service) FindBrandByID(ctx context.Context, id int) (model.Brand, error) {
	panic("Unimplementeed")
}

func (s *service) CreateBrand(ctx context.Context, req product.BrandRequest) (model.Brand, error) {
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

func (s *service) UpdateBrand(ctx context.Context, id int, req product.BrandRequest) error {
	var oldLogoURL *string
	var newLogoURL *string

	// 1️⃣ Cek dulu apakah brand-nya ada, sebelum upload file baru
	existingBrand, err := s.store.FindBrandByID(ctx, id)
	if err != nil {
		return db_error.HandlePgError(err)
	}

	oldLogoURL = existingBrand.LogoURL

	// 2️⃣ Jika ada file baru → upload dulu (tapi jangan commit apa pun dulu)
	if req.File != nil {
		logo, errUpload := s.uploader.UploadBrandLogo(ctx, req.File)
		if errUpload != nil {
			return apperror.New(apperror.ErrCodeInternal, "gagal mengunggah logo baru")
		}
		newLogoURL = &logo
	}

	// 3️⃣ Update data di DB pakai transaksi
	err = s.store.ExecTx(ctx, func(repoInTx product.ProductRepository) error {
		// Update field nama
		existingBrand.Name = req.Name

		// Update logo jika ada file baru
		if newLogoURL != nil {
			existingBrand.LogoURL = newLogoURL
		}

		// Jalankan update
		return repoInTx.UpdateBrand(ctx, existingBrand)
	})

	// 4️⃣ Kalau transaksi gagal → rollback upload logo baru
	if err != nil {
		if newLogoURL != nil {
			log.Printf("[ROLLBACK] Update brand %d gagal, hapus logo baru: %s", id, *newLogoURL)
			if delErr := s.uploader.DeleteByURL(context.Background(), *newLogoURL); delErr != nil {
				log.Printf("[CRITICAL] Gagal hapus logo baru saat rollback: %v", delErr)
			}
		}
		return db_error.HandlePgError(err)
	}

	// 5️⃣ Jika transaksi berhasil → hapus logo lama (jika ada logo baru)
	if newLogoURL != nil && oldLogoURL != nil && *newLogoURL != *oldLogoURL {
		log.Printf("[CLEANUP] Hapus logo lama: %s", *oldLogoURL)
		if delErr := s.uploader.DeleteByURL(context.Background(), *oldLogoURL); delErr != nil {
			log.Printf("[WARN] Gagal hapus logo lama: %v", delErr)
		}
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
		return model.ProductCondition{}, db_error.HandlePgError(err)
	}
	return result, nil
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
	return condition, db_error.HandlePgError(err)
}

func (s *service) UpdateCondition(ctx context.Context, id int16, req product.ProductConditionRequest) (model.ProductCondition, error) {
	var updatedCondition model.ProductCondition
	err := s.store.ExecTx(ctx, func(repoInTx product.ProductRepository) error {
		conditionToUpdate, errTx := repoInTx.FindConditionByID(ctx, id)
		if errTx != nil {
			return errTx
		}
		conditionToUpdate.Name = req.Name
		errTx = repoInTx.UpdateCondition(ctx, conditionToUpdate)
		if errTx != nil {
			return errTx
		}
		return nil
	})
	return updatedCondition, db_error.HandlePgError(err)
}

func (s *service) DeleteCondition(ctx context.Context, id int16) error {
	// 1. Mulai transaksi
	err := s.store.ExecTx(ctx, func(r product.ProductRepository) error {

		// 2. Cek apakah data exists atau tidak
		_, errTx := r.FindConditionByID(ctx, id)
		if errTx != nil {
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

func (s *service) CreateProduct(ctx context.Context, accountID uuid.UUID, request product.CreateProductRequest) (product.ProductDTO, error) {
	// 1. Validasi akun memiliki toko
	log.Printf("[SERVICE-INFO] Mencari shop untuk accountID: %v", accountID)
	shop, err := s.store.FindShopByAccountID(ctx, accountID)
	if err != nil {
		log.Printf("[DB-INFO] Akun tidak ditemukan: %v | %s", accountID, err.Error())
		return product.ProductDTO{}, apperror.New(apperror.ErrCodeNotFound, "Shop not found")
	}
	log.Printf("[SERVICE-INFO] Shop ditemukan: %v", shop.ID)

	// 2. Upload thumbnail (WAJIB)
	log.Println("[SERVICE-INFO] Mengupload thumbnail product")
	var uploadedImageURLs []string
	cleanupCloudImages := func() {
		log.Println("[SERVICE-INFO] Trigger cleanup cloud images akibat error")
		cleanupCtx := context.Background()
		for _, url := range uploadedImageURLs {
			if delErr := s.uploader.DeleteByURL(cleanupCtx, url); delErr != nil {
				log.Printf("[SERVICE-ERROR] Gagal hapus image %s: %v", url, delErr)
			} else {
				log.Printf("[SERVICE-INFO] Berhasil hapus image %s", url)
			}
		}
	}

	// Upload thumbnail
	thumbnailFile, err := request.Thumbnail.Open()
	if err != nil {
		log.Printf("[SERVICE-ERROR] Gagal membuka file thumbnail: %v", err)
		return product.ProductDTO{}, apperror.New(apperror.ErrCodeInternal, "failed to open thumbnail file")
	}
	thumbnailURL, err := s.uploader.UploadProductImage(ctx, thumbnailFile)
	thumbnailFile.Close()
	if err != nil {
		log.Printf("[SERVICE-ERROR] Gagal upload thumbnail: %v", err)
		return product.ProductDTO{}, apperror.New(apperror.ErrCodeInternal, "failed to upload thumbnail")
	}
	uploadedImageURLs = append(uploadedImageURLs, thumbnailURL)
	log.Printf("[SERVICE-INFO] Berhasil upload thumbnail: %s", thumbnailURL)

	// 3. Upload gambar lain (jika ada)
	for _, fileHeader := range request.Images {
		file, err := fileHeader.Open()
		if err != nil {
			log.Printf("[SERVICE-ERROR] Gagal membuka file image: %v", err)
			cleanupCloudImages()
			return product.ProductDTO{}, apperror.New(apperror.ErrCodeInternal, "failed to open image file")
		}

		url, err := s.uploader.UploadProductImage(ctx, file)
		file.Close()
		if err != nil {
			log.Printf("[SERVICE-ERROR] Gagal upload image: %v", err)
			cleanupCloudImages()
			return product.ProductDTO{}, apperror.New(apperror.ErrCodeInternal, "failed to upload image")
		}

		log.Printf("[SERVICE-INFO] Berhasil upload image tambahan: %s", url)
		uploadedImageURLs = append(uploadedImageURLs, url)
	}

	// 4. Transaksi database: simpan produk + gambar
	log.Printf("[SERVICE-INFO] Memulai transaksi untuk menyimpan product dan images")
	var createdProduct model.Product
	err = s.store.ExecTx(ctx, func(repoInTx product.ProductRepository) error {
		productData := model.Product{
			ShopID:      shop.ID,
			Name:        request.Name,
			CategoryID:  request.CategoryID,
			ConditionID: request.ConditionID,
			Price:       request.Price,
			Stock:       request.Stock,
			Description: &request.Description,
			Summary:     &request.Summary,
			BrandID:     &request.BrandID,
			SizeID:      &request.SizeID,
		}

		createdProduct, err = repoInTx.CreateProduct(ctx, productData)
		if err != nil {
			log.Printf("[DB-ERROR] Gagal create product: %v", err)
			return err
		}
		log.Printf("[DB-INFO] Product berhasil dibuat: %v", createdProduct.ID)

		// Data gambar lain (jika ada)
		if len(request.Images) > 0 {
			var images []model.ProductImage
			for i, url := range uploadedImageURLs[1:] { // [1:] skip thumbnail
				images = append(images, model.ProductImage{
					ProductID:  createdProduct.ID,
					ImageIndex: int16(i + 1),
					URL:        url,
				})
			}

			if err = repoInTx.CreateProductImages(ctx, images); err != nil {
				log.Printf("[DB-ERROR] Gagal create product images: %v", err)
				return err
			}
			log.Printf("[DB-INFO] Product images berhasil dibuat untuk productID: %v", createdProduct.ID)
		}

		return nil
	})

	if err != nil {
		log.Printf("[SERVICE-ERROR] Gagal menyimpan product atau images: %v", err)
		cleanupCloudImages()
		return product.ProductDTO{}, db_error.HandlePgError(err)
	}

	// 5. Siapkan response
	log.Printf("[SERVICE-INFO] Menyiapkan response product")
	var imageResponses []product.ProductImageDTO
	for i, url := range uploadedImageURLs { // skip thumbnail
		imageResponses = append(imageResponses, product.ProductImageDTO{
			URL:   url,
			Index: i + 1,
		})
	}

	response := product.ProductDTO{
		ID:          createdProduct.ID,
		ShopID:      createdProduct.ShopID,
		Name:        createdProduct.Name,
		Price:       createdProduct.Price,
		Stock:       createdProduct.Stock,
		Summary:     *createdProduct.Summary,
		Description: *createdProduct.Description,
		Images:      imageResponses,
	}

	log.Printf("[SERVICE-INFO] Product creation selesai: %v", response.ID)
	return response, nil
}

func (s *service) UpdateProduct(ctx context.Context, accountID uuid.UUID, request product.UpdateProductDTO) (product.ProductDTO, error) {

	// 1. ektraksi dan validasi product ID
	productID, err := uuid.Parse(request.ID)
	if err != nil {
		return product.ProductDTO{}, apperror.New(apperror.ErrCodeBadRequest, "invalid product ID format")
	}

	// 2. Validasi akun memiliki toko dan produk ada di toko tersebut
	shop, err := s.store.FindShopByAccountID(ctx, accountID)
	if err != nil {
		return product.ProductDTO{}, apperror.New(apperror.ErrCodeNotFound, "shop not found")
	}

	// 3. mengambil data produk lama
	oldProduct, err := s.store.FindProductByIDAndShop(ctx, productID, shop.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return product.ProductDTO{}, apperror.New(apperror.ErrCodeNotFound, "product not found in your shop")
		}
		return product.ProductDTO{}, db_error.HandlePgError(err)
	}

	// 3️⃣ Update field yang dikirim (non-nil)
	if request.Name != nil {
		oldProduct.Name = *request.Name
	}
	if request.CategoryID != nil {
		oldProduct.CategoryID = *request.CategoryID
	}
	if request.ConditionID != nil {
		oldProduct.ConditionID = *request.ConditionID
	}
	if request.Price != nil {
		oldProduct.Price = *request.Price
	}
	if request.Stock != nil {
		oldProduct.Stock = *request.Stock
	}
	if request.Description != nil {
		oldProduct.Description = request.Description
	}
	if request.Summary != nil {
		oldProduct.Summary = request.Summary
	}
	if request.BrandID != nil {
		oldProduct.BrandID = request.BrandID
	}
	if request.SizeID != nil {
		oldProduct.SizeID = request.SizeID
	}

	oldProduct.UpdatedAt = time.Now()

	// 4️⃣ Simpan perubahan
	updatedProduct, err := s.store.UpdateProduct(ctx, oldProduct)
	if err != nil {
		return product.ProductDTO{}, db_error.HandlePgError(err)
	}

	// 5. ambil data produk lengkap dengan categori, size, brand dll, lengkapp
	updatedProduct, err = s.store.FindProductByID(ctx, updatedProduct.ID)
	if err != nil {
		return product.ProductDTO{}, db_error.HandlePgError(err)
	}

	resp := product.ToProductDTO(updatedProduct)

	log.Printf("[SERVICE-INFO] UpdateProduct berhasil untuk ID: %v", resp.ID)
	return resp, nil
}

// FindProductsBySeller implements product.ProductService.
func (s *service) FindProductsBySeller(
	ctx context.Context,
	accountID uuid.UUID,
	filter product.ProductFilterDTO,
	page int,
	size int,
) (product.PaginatedProductDTO, error) {

	// 🔹 Ambil data dari repository
	products, total, err := s.store.FindProductsBySeller(ctx, accountID, filter, page, size)
	if err != nil {
		return product.PaginatedProductDTO{}, err
	}

	// 🔹 Mapping ke DTO
	items := make([]product.ProductDTO, len(products))
	for i, p := range products {
		items[i] = product.ToProductDTO(p)
	}

	// 🔹 Hitung pagination
	totalPages := 0
	if size > 0 {
		totalPages = (total + size - 1) / size
	}

	paginated := product.PaginatedProductDTO{
		Page:       page,
		Size:       size,
		TotalItems: total,
		TotalPages: totalPages,
		Items:      items,
	}

	return paginated, nil
}

// helper kecil biar gak nulis manual
func ptrInt(v int) *int {
	return &v
}

func NewProductService(store repository.ProductStore, jwt auth.JWTService, uploader uploader.Uploader) product.ProductService {
	return &service{
		store:    store,
		jwt:      jwt,
		uploader: uploader,
	}
}
