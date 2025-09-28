// File: pkg/uploader/cloudinary.go
package uploader

import (
	"context"
	"io"
	
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

// Uploader mendefinisikan kontrak untuk layanan upload file.
// Dengan interface, kita bisa dengan mudah menggantinya ke GCS atau S3 nanti.
type Uploader interface {
	Upload(ctx context.Context, file io.Reader, filename string) (url string, err error)
}

type cloudinaryUploader struct {
	cld *cloudinary.Cloudinary
}

// NewCloudinaryUploader adalah constructor.
// cloudURL adalah CLOUDINARY_URL dari file .env-mu.
func NewCloudinaryUploader(cloudURL string) (Uploader, error) {
	cld, err := cloudinary.NewFromURL(cloudURL)
	if err != nil {
		return nil, err
	}
	return &cloudinaryUploader{cld: cld}, nil
}


// Upload mengimplementasikan logika upload ke Cloudinary.
func (u *cloudinaryUploader) Upload(ctx context.Context, file io.Reader, filename string) (string, error) {
	uploadParams := uploader.UploadParams{
		// PublicID bisa di-set agar nama file di Cloudinary lebih rapi
		// Contoh: "products/nama-file-unik"
		PublicID: filename,
		// Opsi lain seperti folder, tag, transformasi bisa ditambahkan di sini
	}

	uploadResult, err := u.cld.Upload.Upload(ctx, file, uploadParams)
	if err != nil {
		return "", err
	}

	// Kembalikan URL yang aman (HTTPS) dari gambar yang sudah di-upload
	return uploadResult.SecureURL, nil
}