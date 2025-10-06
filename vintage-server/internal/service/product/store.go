package product

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// Store mendefinisikan semua fungsi untuk berinteraksi dengan database.
// Ia "mewarisi" semua metode dari Repository dan menambahkan kemampuan transaksi.
type Store interface {
	Repository
	ExecTx(ctx context.Context, fn func(Repository) error) error
}

// sqlStore menyediakan implementasi dari Store untuk database SQL.
type sqlStore struct {
	db *sqlx.DB
	Repository
}

// NewStore adalah "pintu masuk" utama dari aplikasi kita ke lapisan data.
// Ia membuat sebuah Store baru yang siap digunakan.
func NewStore(db *sqlx.DB) Store {
	return &sqlStore{
		db:         db,
		Repository: NewRepository(db), // Menggunakan konstruktor Repository yang sudah kita buat
	}
}

// ExecTx adalah metode andalan kita.
// Ia mengeksekusi sebuah fungsi di dalam sebuah transaksi database.
func (s *sqlStore) ExecTx(ctx context.Context, fn func(Repository) error) error {
	// 1. Memulai transaksi ("safety bubble")
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	// 2. Membuat repository khusus yang berjalan di dalam transaksi ini
	repoTx := s.Repository.WithTx(tx)

	// 3. Menjalankan "resep" dari Service Layer (fungsi 'fn')
	//    dengan memberikan repository transaksional.
	err = fn(repoTx)
	if err != nil {
		// 4a. Jika resep gagal, batalkan semua perubahan (Rollback)
		if rbErr := tx.Rollback(); rbErr != nil {
			// Jika rollback juga gagal, laporkan kedua error
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err // Kembalikan error asli dari resep
	}

	// 4b. Jika resep berhasil, simpan semua perubahan (Commit)
	return tx.Commit()
}
