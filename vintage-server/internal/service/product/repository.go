package product

import (
	"context"
	"database/sql"
	"vintage-server/internal/model"
	"vintage-server/pkg/apperror"

	"github.com/jmoiron/sqlx"
)

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) CreateCategory(ctx context.Context, data ProductCategory) error {
	// 1. Cek apakah kategori dengan nama yang sama (case-insensitive) sudah ada.
	// Kita hanya perlu mengambil satu kolom (misal 'id') untuk verifikasi.
	queryCheck := `SELECT id FROM product_categories WHERE LOWER(name) = LOWER($1)`

	var existingID string // Variabel untuk menampung ID jika ditemukan

	err := r.db.GetContext(ctx, &existingID, queryCheck, data.Name)

	// 2. Analisis hasil pengecekan
	if err == nil {
		// Jika err == nil, artinya query berhasil menemukan satu baris. Kategori sudah ada.
		return apperror.New(apperror.ErrCodeConflict, "category with this name already exists")
	} else if err != sql.ErrNoRows {
		// Jika error BUKAN karena tidak ada baris (ErrNoRows), berarti ada masalah lain dengan database.
		// Kita kembalikan error database tersebut.
		return err
	}

	// 3. Jika err == sql.ErrNoRows, berarti kategori belum ada. Lanjutkan proses INSERT.
	// Pastikan untuk INSERT nama aslinya (data.Name), bukan versi lowercase-nya.
	queryInsert := `INSERT INTO product_categories (name) VALUES ($1)`

	_, err = r.db.ExecContext(ctx, queryInsert, data.Name)
	if err != nil {
		// Tangani jika ada error saat proses INSERT
		return err
	}

	return nil
}

func (r *repository) FindAllCategories(ctx context.Context) ([]model.ProductCategory, error) {
	query := "SELECT * FROM product_categories"
	var result []model.ProductCategory

	// Gunakan SelectContext untuk mengambil banyak baris
	err := r.db.SelectContext(ctx, &result, query)
	if err != nil {
		return []model.ProductCategory{}, apperror.HandleDBError(err, "Failed to read data")
	}
	// err akan nil jika tidak ada error, termasuk jika tidak ada baris yang ditemukan
	return result, err
}

func (r *repository) FindById(ctx context.Context, id int) (model.ProductCategory, error) {
	// 1. Deklarasi variabel untuk menampung hasil. (Typo 'resukt' diperbaiki).
	var result model.ProductCategory

	// 2. Tulis query dengan menyebutkan kolom secara eksplisit.
	// Gunakan placeholder $1 untuk keamanan (mencegah SQL Injection).
	query := `SELECT id, name, created_at, updated_at FROM product_categories WHERE id = $1`

	// 3. Eksekusi query menggunakan GetContext.
	// GetContext ideal untuk mengambil satu baris data.
	err := r.db.GetContext(ctx, &result, query, id)
	if err != nil {
		// 4. Gunakan helper untuk menangani error.
		// Helper ini akan otomatis mengubah sql.ErrNoRows menjadi apperror.NotFound
		// dan error lainnya menjadi apperror.Internal.
		return model.ProductCategory{}, apperror.HandleDBError(err, "failed to find product category by id")
	}

	// 5. Jika tidak ada error, kembalikan hasil.
	return result, nil
}