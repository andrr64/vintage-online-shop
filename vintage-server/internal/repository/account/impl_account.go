// InsertAccount implements account.AccountRepository.
package repository

import (
	"context"
	"vintage-server/internal/model"

	"github.com/google/uuid"
)

func (r *sqlAccountRepository) InsertAccount(ctx context.Context, account model.Account) (model.Account, error) {
	query := `
		INSERT INTO accounts (
			firstname, lastname, username, email, password, active, created_at, updated_at
		) VALUES (
			:firstname, :lastname, :username, :email, :password, :active, :created_at, :updated_at
		)
		RETURNING id, firstname, lastname, username, email, active, avatar_url, created_at, updated_at
	`

	var savedAccount model.Account
	nstmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return model.Account{}, err
	}
	defer nstmt.Close()

	if err := nstmt.GetContext(ctx, &savedAccount, account); err != nil {
		return model.Account{}, err
	}

	return savedAccount, nil
}

// GetRoleIDByName implements account.AccountRepository.
func (r *sqlAccountRepository) GetRoleIDByName(ctx context.Context, roleName string) (int64, error) {
	var roleID int64
	err := r.db.GetContext(ctx, &roleID, `SELECT id FROM roles WHERE name = $1 LIMIT 1`, roleName)
	if err != nil {
		return 0, err
	}

	return roleID, nil
}

// InsertAccountRole implements account.AccountRepository.
func (r *sqlAccountRepository) InsertAccountRole(ctx context.Context, accountID uuid.UUID, roleID int64) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO account_roles (account_id, role_id) VALUES ($1, $2)`,
		accountID, roleID)
	return err
}

// FindAccountByEmail implements account.AccountRepository.
func (r *sqlAccountRepository) FindAccountByEmail(ctx context.Context, email string) (model.Account, error) {
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

// FindAccountByEmailWithRole implements account.AccountRepository.
func (r *sqlAccountRepository) FindAccountByEmailWithRole(ctx context.Context, email string, roleName string) (model.Account, error) {
	var account model.Account
	query := `
		SELECT a.*
		FROM accounts AS a
		JOIN account_roles AS ar ON a.id = ar.account_id
		JOIN roles AS r ON r.id = ar.role_id
		WHERE a.email = $1 AND r.name = $2
		LIMIT 1
	`
	err := r.db.GetContext(ctx, &account, query, email, roleName)
	return account, err
}

// FindAccountByID implements account.AccountRepository.
func (r *sqlAccountRepository) FindAccountByID(ctx context.Context, id uuid.UUID) (model.Account, error) {
	var account model.Account
	query := "SELECT * FROM accounts WHERE id = $1"
	err := r.db.GetContext(ctx, &account, query, id)
	return account, err
}

// FindAccountByUsername implements account.AccountRepository.
func (r *sqlAccountRepository) FindAccountByUsername(ctx context.Context, username string) (model.Account, error) {
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

// FindAccountByUsernameWithRole implements account.AccountRepository.
func (r *sqlAccountRepository) FindAccountByUsernameWithRole(ctx context.Context, username string, roleName string) (model.Account, error) {
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

// IsUsernameUsed implements account.AccountRepository.
func (r *sqlAccountRepository) IsUsernameUsed(ctx context.Context, username string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS (SELECT 1 FROM accounts WHERE username = $1)"
	err := r.db.GetContext(ctx, &exists, query, username)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// UpdateAccount implements account.AccountRepository.
func (r *sqlAccountRepository) UpdateAccount(ctx context.Context, account model.Account) (model.Account, error) {
	query := `
		UPDATE accounts SET 
			username = :username, 
			avatar_url = :avatar_url, 
			firstname = :firstname,
			lastname = :lastname,
			active = :active,
			updated_at = :updated_at
		WHERE id = :id
		RETURNING *;
	`

	var updatedAccount model.Account
	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return model.Account{}, err
	}
	defer stmt.Close()

	if err := stmt.GetContext(ctx, &updatedAccount, account); err != nil {
		return model.Account{}, err
	}

	return updatedAccount, nil
}

// UpdateAvatarTx implements account.AccountRepository.
func (r *sqlAccountRepository) UpdateAvatar(ctx context.Context, avatarUrl string, id uuid.UUID) (model.Account, error) {
	query := `
		UPDATE accounts 
		SET avatar_url = $1, updated_at = NOW()
		WHERE id = $2
		RETURNING id, username, firstname, lastname, email, avatar_url;
	`

	var savedAccount model.Account
	if err := r.db.GetContext(ctx, &savedAccount, query, avatarUrl, id); err != nil {
		return model.Account{}, err
	}

	return savedAccount, nil
}


func (r *sqlAccountRepository) FindProductByID(ctx context.Context, productID uuid.UUID) (model.Product, error) {
	// 🔹 Ambil data utama produk
	query := `
		SELECT 
			p.id, p.shop_id, p.condition_id, p.category_id, p.brand_id, p.size_id,
			p.name, p.summary, p.description, p.price, p.stock,
			p.created_at, p.updated_at
		FROM products p
		WHERE p.id = $1
	`
	var product model.Product
	if err := r.db.GetContext(ctx, &product, query, productID); err != nil {
		return model.Product{}, err
	}

	// 🔹 Ambil relasi: Brand
	if product.BrandID != nil {
		var brand model.Brand
		err := r.db.GetContext(ctx, &brand, `SELECT id, name, logo_url, created_at, updated_at FROM brands WHERE id = $1`, *product.BrandID)
		if err == nil {
			product.Brand = &brand
		}
	}

	// 🔹 Ambil relasi: Category
	var category model.ProductCategory
	err := r.db.GetContext(ctx, &category, `SELECT id, name, created_at, updated_at FROM product_categories WHERE id = $1`, product.CategoryID)
	if err == nil {
		product.Category = &category
	}

	// 🔹 Ambil relasi: Condition
	var condition model.ProductCondition
	err = r.db.GetContext(ctx, &condition, `SELECT id, name, created_at, updated_at FROM product_conditions WHERE id = $1`, product.ConditionID)
	if err == nil {
		product.Condition = &condition
	}

	// 🔹 Ambil relasi: Size
	if product.SizeID != nil {
		var size model.ProductSize
		err := r.db.GetContext(ctx, &size, `SELECT id, name FROM product_size WHERE id = $1`, *product.SizeID)
		if err == nil {
			product.Size = &size
		}
	}

	// 🔹 Ambil relasi: Shop
	var shop model.Shop
	err = r.db.GetContext(ctx, &shop, `
		SELECT id, account_id, name, summary, description, active, created_at, updated_at
		FROM shop WHERE id = $1
	`, product.ShopID)
	if err == nil {
		product.Shop = &shop
	}

	// 🔹 Ambil relasi: Images
	var images []model.ProductImage
	err = r.db.SelectContext(ctx, &images, `
		SELECT id, product_id, image_index, url, created_at
		FROM product_images
		WHERE product_id = $1
		ORDER BY image_index ASC
	`, product.ID)
	if err == nil {
		product.Images = images
	}

	return product, nil
}