package product

import (
	"net/http"
	"strconv"
	product "vintage-server/internal/domain"
	"vintage-server/pkg/helper"
	"vintage-server/pkg/response"
	"vintage-server/pkg/utils"

	"github.com/gin-gonic/gin"
)

// -- BRAND MANAGEMENT --
func (h *handler) CreateBrand(c *gin.Context) {
	_, err := helper.CheckAuthAndRole(c, "admin")
	if err != nil {
		response.ErrorForbidden(c)
		return
	}

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
	req := product.BrandRequest{
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

func (h *handler) UpdateBrand(c *gin.Context) {
	_, err := helper.CheckAuthAndRole(c, "admin")
	if err != nil {
		response.ErrorForbidden(c)
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
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
	req := product.BrandRequest{
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

func (h *handler) ReadBrand(c *gin.Context) {
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

func (h *handler) FindBrandByID(c *gin.Context) {
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

func (h *handler) DeleteBrand(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
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
