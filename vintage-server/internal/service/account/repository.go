package user

import (
	"context"
	"vintage-server/internal/model" // Sesuaikan path

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// repository adalah struct yang mengimplementasikan kontrak Repository dari domain.go
type repository struct {
	db *sqlx.DB
}

// NewRepository adalah constructor untuk implementasi repository
func NewRepository(db *sqlx.DB) Repository {
	return &repository{
		db: db,
	}
}

// --- Account ---
func (r *repository) IsUsernameUsed(ctx context.Context, username string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS (SELECT 1 FROM accounts WHERE username = $1)"
	err := r.db.GetContext(ctx, &exists, query, username)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *repository) FindAccountByID(ctx context.Context, id uuid.UUID) (model.Account, error) {
	var account model.Account
	query := "SELECT * FROM accounts WHERE id = $1"
	err := r.db.GetContext(ctx, &account, query, id)
	return account, err
}

func (r *repository) FindAccountByEmailWithRole(ctx context.Context, email string, roleName string) (model.Account, error) {
	var account model.Account
	query := `
		SELECT a.*
		FROM accounts a
		JOIN account_roles ar ON a.id = ar.account_id
		JOIN roles r ON ar.role_id = r.id
		WHERE a.email = $1 AND r.name = $2
		LIMIT 1
	`
	err := r.db.GetContext(ctx, &account, query, email, roleName)
	return account, err
}

func (r *repository) FindAccountByUsernameWithRole(ctx context.Context, username string, roleName string) (model.Account, error) {
	var account model.Account
	query := `
		SELECT a.*
		FROM accounts a
		JOIN account_roles ar ON a.id = ar.account_id
		JOIN roles r ON ar.role_id = r.id
		WHERE a.username = $1 AND r.name = $2
		LIMIT 1
	`
	err := r.db.GetContext(ctx, &account, query, username, roleName)
	return account, err
}

func (r *repository) FindAccountByUsername(ctx context.Context, username string) (model.Account, error) {
	var account model.Account
	query := `
		SELECT a.* 
		FROM accounts a
		WHERE a.username = $1
		LIMIT 1
	`
	err := r.db.GetContext(ctx, &account, query, username)
	return account, err
}

func (r *repository) FindAccountByEmail(ctx context.Context, email string) (model.Account, error) {
	var account model.Account
	query := `
		SELECT a.* 
		FROM accounts a
		WHERE a.email = $1
		LIMIT 1
	`
	err := r.db.GetContext(ctx, &account, query, email)
	return account, err
}

// File: internal/service/user/repository.go

func (r *repository) SaveAccount(ctx context.Context, account model.Account, roleName string) (savedAccount model.Account, err error) {
	// 1. Persiapan transaksi
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return model.Account{}, err
	}

	// 2. Fungsi defer commit atau rollback
	defer func() {
		if p := recover(); p != nil || err != nil {
			tx.Rollback()
			if p != nil {
				panic(p)
			}
		} else {
			err = tx.Commit()
		}
	}()

	// 3. Insert ke tabel 'accounts'
	queryAcc := `
        INSERT INTO accounts ( firstname, lastname, username, email, password, active, created_at, updated_at)
        VALUES (:firstname, :lastname, :username, :email, :password, :active, :created_at, :updated_at)
        RETURNING *`

	// ================= PERBAIKAN DI SINI =================
	// Siapkan named statement di dalam transaksi
	stmt, err := tx.PrepareNamedContext(ctx, queryAcc)
	if err != nil {
		return model.Account{}, err
	}
	defer stmt.Close() // Jangan lupa tutup statement

	// Eksekusi statement dengan data dari struct 'account'
	rows, err := stmt.QueryxContext(ctx, account)
	if err != nil {
		return model.Account{}, err
	}
	defer rows.Close()
	// =======================================================

	if rows.Next() {
		if err = rows.StructScan(&savedAccount); err != nil {
			return model.Account{}, err
		}
	}

	// ... (sisa kodenya sudah benar) ...
	var roleID int64
	err = tx.GetContext(ctx, &roleID, "SELECT id FROM roles WHERE name = $1", roleName)
	if err != nil {
		return model.Account{}, err
	}

	_, err = tx.ExecContext(ctx, "INSERT INTO account_roles (account_id, role_id) VALUES ($1, $2)", savedAccount.ID, roleID)
	if err != nil {
		return model.Account{}, err
	}

	return savedAccount, nil
}

func (r *repository) UpdateAccount(ctx context.Context, account model.Account) error {
	query := `UPDATE accounts SET 
				username = :username, 
				avatar_url = :avatar_url, 
				firstname = :firstname,
				lastname = :lastname,
				active = :active,
				updated_at = :updated_at
			  WHERE id = :id`
	_, err := r.db.NamedExecContext(ctx, query, account)
	return err
}

// --- Address ---

func (r *repository) SaveAddress(ctx context.Context, address model.Address) (model.Address, error) {
	var savedAddress model.Address
	query := `
		INSERT INTO addresses (account_id, district_id, regency_id, province_id, label, recipient_name, recipient_phone, street, postal_code, is_primary, created_at, updated_at)
		VALUES (:account_id, :district_id, :regency_id, :province_id, :label, :recipient_name, :recipient_phone, :street, :postal_code, :is_primary, :created_at, :updated_at)
		RETURNING *`

	rows, err := r.db.NamedQueryContext(ctx, query, address)
	if err != nil {
		return model.Address{}, err
	}
	defer rows.Close()
	if rows.Next() {
		if err := rows.StructScan(&savedAddress); err != nil {
			return model.Address{}, err
		}
	}
	return savedAddress, nil
}

func (r *repository) FindAddressesByAccountID(ctx context.Context, accountID int64) ([]model.Address, error) {
	var addresses []model.Address
	query := "SELECT * FROM addresses WHERE account_id = $1 ORDER BY is_primary DESC, updated_at DESC"
	err := r.db.SelectContext(ctx, &addresses, query, accountID)
	return addresses, err
}

func (r *repository) FindAddressByIDAndAccountID(ctx context.Context, addressID, accountID int64) (model.Address, error) {
	var address model.Address
	query := "SELECT * FROM addresses WHERE id = $1 AND account_id = $2"
	err := r.db.GetContext(ctx, &address, query, addressID, accountID)
	return address, err
}

func (r *repository) UpdateAddress(ctx context.Context, address model.Address) (model.Address, error) {
	var updatedAddress model.Address
	query := `
		UPDATE addresses SET
			label = :label,
			recipient_name = :recipient_name,
			recipient_phone = :recipient_phone,
			street = :street,
			postal_code = :postal_code,
			updated_at = :updated_at
		WHERE id = :id AND account_id = :account_id
		RETURNING *`

	rows, err := r.db.NamedQueryContext(ctx, query, address)
	if err != nil {
		return model.Address{}, err
	}
	defer rows.Close()
	if rows.Next() {
		if err := rows.StructScan(&updatedAddress); err != nil {
			return model.Address{}, err
		}
	}
	return updatedAddress, nil
}

func (r *repository) DeleteAddress(ctx context.Context, addressID, accountID int64) error {
	query := "DELETE FROM addresses WHERE id = $1 AND account_id = $2"
	_, err := r.db.ExecContext(ctx, query, addressID, accountID)
	return err
}

func (r *repository) TransactionSetPrimaryAddress(ctx context.Context, accountID, addressID int64) error {
	// Memulai transaksi
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	// Defer rollback jika terjadi panic atau error
	defer tx.Rollback()

	// Query 1: Set semua alamat user menjadi BUKAN primary
	queryUnset := "UPDATE addresses SET is_primary = FALSE WHERE account_id = $1"
	if _, err := tx.ExecContext(ctx, queryUnset, accountID); err != nil {
		return err
	}

	// Query 2: Set alamat yang dipilih menjadi primary
	querySet := "UPDATE addresses SET is_primary = TRUE WHERE id = $1 AND account_id = $2"
	if _, err := tx.ExecContext(ctx, querySet, addressID, accountID); err != nil {
		return err
	}

	// Jika semua berhasil, commit transaksi
	return tx.Commit()
}

// --- Wishlist ---

func (r *repository) SaveWishlistItem(ctx context.Context, item model.Wishlist) error {
	query := `INSERT INTO wishlist (account_id, product_id, created_at, updated_at)
			  VALUES (:account_id, :product_id, :created_at, :updated_at)`
	_, err := r.db.NamedExecContext(ctx, query, item)
	return err
}

func (r *repository) FindWishlistByAccountID(ctx context.Context, accountID int64) ([]WishlistItemDetail, error) {
	var wishlistItems []WishlistItemDetail
	// Query ini melakukan JOIN antara tabel wishlist, products, dan product_images
	// untuk mengambil data yang dibutuhkan oleh DTO WishlistItemDetail.
	query := `
		SELECT 
			w.product_id,
			p.name as product_name,
			p.price,
			pi.url as product_image_url,
			w.created_at
		FROM wishlist w
		JOIN products p ON w.product_id = p.id
		LEFT JOIN product_images pi ON p.id = pi.product_id AND pi.image_index = 0
		WHERE w.account_id = $1
		ORDER BY w.created_at DESC`

	err := r.db.SelectContext(ctx, &wishlistItems, query, accountID)
	return wishlistItems, err
}

func (r *repository) DeleteWishlistItem(ctx context.Context, accountID, productID int64) error {
	query := "DELETE FROM wishlist WHERE account_id = $1 AND product_id = $2"
	_, err := r.db.ExecContext(ctx, query, accountID, productID)
	return err
}

func (r *repository) CheckWishlistItemExists(ctx context.Context, accountID, productID int64) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM wishlist WHERE account_id = $1 AND product_id = $2)"
	err := r.db.GetContext(ctx, &exists, query, accountID, productID)
	return exists, err
}
