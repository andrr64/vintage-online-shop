// File: pkg/uploader/cloudinary.go
package uploader

import (
	"context"
	"io"
	"path"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

// Uploader mendefinisikan kontrak untuk layanan upload file.
// Dengan interface, kita bisa dengan mudah menggantinya ke GCS atau S3 nanti.
type Uploader interface {
	Upload(ctx context.Context, file io.Reader, filename string) (url string, err error)
	Delete(ctx context.Context, publicID string) error
}

type cloudinaryUploader struct {
	cld *cloudinary.Cloudinary
}

func extractPublicID(url string) string {
	// Contoh: https://res.cloudinary.com/<cloud>/image/upload/v12345/vintage/avatars/abc-avatar.png
	parts := strings.Split(url, "/")
	if len(parts) == 0 {
		return ""
	}

	// buang bagian versi (v12345) → ambil setelah itu
	for i, p := range parts {
		if strings.HasPrefix(p, "v") && len(p) > 1 {
			// ambil semua setelah versi
			publicID := path.Join(parts[i+1:]...)
			// buang ekstensi (".png", ".jpg", dll)
			dot := strings.LastIndex(publicID, ".")
			if dot != -1 {
				publicID = publicID[:dot]
			}
			return publicID
		}
	}
	return ""
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
	overwrite := true

	uploadParams := uploader.UploadParams{
		// PublicID bisa di-set agar nama file di Cloudinary lebih rapi
		// Contoh: "products/nama-file-unik"
		PublicID:  filename,
		Tags:      api.CldAPIArray{"vintage", "avatar"},
		Folder:    "vintage/avatars",
		Overwrite: &overwrite,
		// Opsi lain seperti folder, tag, transformasi bisa ditambahkan di sini
	}

	uploadResult, err := u.cld.Upload.Upload(ctx, file, uploadParams)
	if err != nil {
		return "", err
	}

	// Kembalikan URL yang aman (HTTPS) dari gambar yang sudah di-upload
	return uploadResult.SecureURL, nil
}

// Implementasi di cloudinaryUploader
func (u *cloudinaryUploader) Delete(ctx context.Context, avatar_url string) error {
	// Panggil API destroy Cloudinary
	_, err := u.cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: extractPublicID(avatar_url),
		// ResourceType biasanya "image", bisa juga "video" dll
		Invalidate: api.Bool(true), // biar cache CDN dihapus juga
	})
	return err
}
