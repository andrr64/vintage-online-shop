package product

import (
	"net/http"
	"strconv"
	"vintage-server/pkg/helper"
	"vintage-server/pkg/response"
	"vintage-server/pkg/utils"

	"github.com/gin-gonic/gin"
)

// Handler adalah struct yang memegang dependency ke Service
type Handler struct {
	svc Service
}

// Perhatikan bahwa return type-nya adalah interface, bukan struct-nya langsung.
func NewHandler(svc Service) ProductHandler {
	return &Handler{svc: svc}
}

// -- CATEGORY MANAGEMENT --
func (h *Handler) CreateCategory(c *gin.Context) {
	_, err := helper.CheckAuthAndRole(c, "admin")
	if err != nil {
		response.ErrorUnauthorized(c)
		return
	}

	var category ProductCategory
	if err := c.ShouldBindJSON(&category); err != nil {
		response.ErrorBadRequest(c, "Invalid request")
		return
	}

	err = h.svc.CreateCategory(c.Request.Context(), category)
	if err != nil {
		helper.HandleError(c, err)
		return
	}
	response.Success(c, http.StatusCreated, "OK bro")
}

func (h *Handler) ReadCategories(c *gin.Context) {
	// Deklarasikan variabel untuk hasil dan error di luar if/else
	var result interface{}
	var err error

	categoryIdStr := c.Query("id")
	if categoryIdStr == "" {
		// --- Kasus 1: Ambil semua kategori ---
		result, err = h.svc.FindAllCategories(c.Request.Context())
	} else {
		// --- Kasus 2: Ambil kategori berdasarkan ID ---
		// Konversi string ID ke integer
		var categoryId int
		categoryId, err = strconv.Atoi(categoryIdStr)
		if err != nil {
			// Jika ID bukan angka, kembalikan error bad request
			response.Error(c, http.StatusBadRequest, "Invalid category ID format")
			return
		}
		result, err = h.svc.FindById(c.Request.Context(), categoryId)
	}

	// --- Penanganan Error & Sukses (hanya ditulis sekali) ---
	if err != nil {
		helper.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, result)
}

func (h *Handler) UpdateCategory(c *gin.Context) {
	_, err := helper.CheckAuthAndRole(c, "admin")
	if err != nil {
		response.ErrorUnauthorized(c)
		return
	}
	var category ProductCategory
	if err := c.ShouldBindJSON(&category); err != nil {
		response.ErrorBadRequest(c, "Invalid request")
		return
	}
	if category.ID == nil {
		response.ErrorBadRequest(c, "Invalid request")
		return
	}

	err = h.svc.UpdateCategory(c.Request.Context(), category)
	if err != nil {
		helper.HandleError(c, err)
		return
	}

	response.SuccessWithoutData(c, http.StatusOK, "OK")
}

func (h *Handler) DeleteCategory(c *gin.Context) {
	_, err := helper.CheckAuthAndRole(c, "admin")
	if err != nil {
		response.ErrorUnauthorized(c)
		return
	}
	id := c.Query("id")
	if id == "" {
		response.ErrorBadRequest(c)
		return
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		response.ErrorBadRequest(c)
	}
	err = h.svc.DeleteCategory(c.Request.Context(), idInt)
	if err != nil {
		helper.HandleError(c, err)
		return
	}
	response.SuccessWithoutData(c, http.StatusOK, "OK")
}

// -- BRAND MANAGEMENT --
func (h *Handler) CreateBrand(c *gin.Context) {
	// 1. Baca data teks dari form-data
	name := c.PostForm("name")
	if name == "" {
		response.Error(c, http.StatusBadRequest, "name field is required")
		return
	}

	// 2. Baca file dari form-data
	fileHeader, err := c.FormFile("logo")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "logo file is required")
		return
	}
	// 3. Validasi file (ukuran, tipe, dll)
	if !utils.SizeIsOk(fileHeader, utils.Megabytes(2)) {
		response.Error(c, http.StatusBadRequest, "file size must be less than 2MB")
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to open uploaded file")
		return
	}
	defer file.Close()

	// 4. Buat request DTO untuk service
	req := CreateBrandRequest{
		Name:       name,
		File:       file,
		FileHeader: fileHeader,
	}

	// 5. Panggil service
	brand, err := h.svc.CreateBrand(c.Request.Context(), req)
	if err != nil {
		helper.HandleError(c, err)
		return
	}
	response.Success(c, http.StatusCreated, brand)
}

func (h *Handler) UpdateBrand(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request")
		return
	}

	// Baca data teks
	name := c.PostForm("name")
	if name == "" {
		response.Error(c, http.StatusBadRequest, "name field is required")
		return
	}

	// Siapkan request untuk service
	req := UpdateBrandRequest{
		Name: name,
	}

	// File untuk update bersifat OPSIONAL
	fileHeader, err := c.FormFile("logo")
	if err != nil && err != http.ErrMissingFile {
		// Jika ada error selain file tidak dilampirkan, itu error sebenarnya
		response.Error(c, http.StatusBadRequest, "invalid file upload")
		return
	}

	// Jika ada file baru yang di-upload, proses file tersebut
	if fileHeader != nil {
		if !utils.SizeIsOk(fileHeader, utils.Megabytes(2)) {
			response.Error(c, http.StatusBadRequest, "file size must be less than 2MB")
			return
		}

		file, err := fileHeader.Open()
		if err != nil {
			response.Error(c, http.StatusInternalServerError, "failed to open uploaded file")
			return
		}
		defer file.Close()

		req.File = file
		req.FileHeader = fileHeader
	}

	// Panggil service
	err = h.svc.UpdateBrand(c.Request.Context(), id, req)
	if err != nil {
		helper.HandleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, gin.H{"message": "brand updated successfully"})
}

func (h *Handler) ReadBrand(c *gin.Context) {
	var data interface{}
	var err error

	id := c.Query("id")

	if id != "" {
		// 1. Lakukan konversi string ke integer terlebih dahulu.
		idInt, convErr := strconv.Atoi(id)

		// 2. Jika konversi GAGAL, ini adalah kesalahan input dari user (Bad Request).
		//    Langsung hentikan proses dan kirim error.
		if convErr != nil {
			// Asumsi helper.HandleError bisa menangani ini, atau Anda bisa langsung:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format. ID must be an integer."})
			return
		}

		// 3. Jika konversi BERHASIL, baru panggil service untuk mencari data.
		data, err = h.svc.FindBrandByID(c.Request.Context(), idInt)

	} else {
		// Logika ini sudah benar.
		data, err = h.svc.FindAllBrands(c.Request.Context())
	}

	// Penanganan error dari service (misal: data tidak ditemukan, DB error).
	if err != nil {
		helper.HandleError(c, err)
		return
	}

	// Kirim respons sukses.
	response.Success(c, http.StatusOK, data)
}

func (h *Handler) FindBrandByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid brand ID format")
		return
	}

	brand, err := h.svc.FindBrandByID(c.Request.Context(), id)
	if err != nil {
		helper.HandleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, brand)
}

func (h *Handler) DeleteBrand(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid brand ID format")
		return
	}

	err = h.svc.DeleteBrand(c.Request.Context(), id)
	if err != nil {
		helper.HandleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, gin.H{"message": "brand deleted successfully"})
}
