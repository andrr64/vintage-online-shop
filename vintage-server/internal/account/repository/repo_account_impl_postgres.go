package repository

import (
	"context"
	"vintage-server/internal/account/domain"
	"vintage-server/internal/account/model"
	"vintage-server/internal/shared/db"
	db_error "vintage-server/pkg"

	"github.com/google/uuid"
)

// WithTx implements AccountRepository.
func (a *accountRepoPostgres) WithTx(tx db.DBTX) AccountRepository {
	return &accountRepoPostgres{db: tx}
}

func (r *accountRepoPostgres) AssignRole(ctx context.Context, accountID uuid.UUID, role string) error {
	query := `
        INSERT INTO account_roles (account_id, role_id)
        SELECT $1, r.id
        FROM roles r
        WHERE r.name = $2
    `
	_, err := r.db.ExecContext(ctx, query, accountID, role)
	return db_error.HandlePgError(err)
}

// CreateAccount implements AccountRepository.
func (r *accountRepoPostgres) CreateAccount(ctx context.Context, account domain.Account) (domain.Account, error) {
	query := `
        INSERT INTO accounts (username, firstname, lastname, password, email, avatar_url, active, created_at, updated_at)
        VALUES ($1,$2,$3,$4,$5,$6,$7,NOW(),NOW())
        RETURNING id, created_at, updated_at
    `
	row := r.db.QueryRowContext(ctx, query,
		account.Username,
		account.Firstname,
		account.Lastname,
		account.Password,
		account.Email,
		account.AvatarURL,
		account.Active,
	)
	if err := row.Scan(&account.ID, &account.CreatedAt, &account.UpdatedAt); err != nil {
		return domain.Account{}, db_error.HandlePgError(err)
	}
	return account, nil
}

func (a *accountRepoPostgres) FindAccountByUsernameAndRoleString(
	ctx context.Context,
	username domain.Username,
	role string,
) (domain.Account, error) {
	query := `
		SELECT a.id, a.username, a.firstname, a.lastname, a.password, a.email, a.avatar_url, a.active, a.created_at, a.updated_at
		FROM accounts a
		JOIN account_roles ar ON a.id = ar.account_id
		JOIN roles r ON ar.role_id = r.id
		WHERE a.username = $1 AND r.name = $2
		LIMIT 1
	`

	var model model.Account
	if err := a.db.GetContext(ctx, &model, query, username.String(), role); err != nil {
		return domain.Account{}, db_error.HandlePgError(err)
	}

	acc := domain.Account{
		ID:        model.ID,
		Username:  model.Username,
		Firstname: model.Firstname,
		Lastname:  model.Lastname,
		Password:  model.Password,
		Email:     model.Email,
		AvatarURL: model.AvatarURL,
		Active:    model.Active,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}

	return acc, nil
}

// ============================
// Email version
// ============================
func (a *accountRepoPostgres) FindAccountByEmailAndRoleString(
	ctx context.Context,
	email domain.Email,
	role string,
) (domain.Account, error) {
	query := `
		SELECT a.id, a.username, a.firstname, a.lastname, a.password, a.email, a.avatar_url, a.active, a.created_at, a.updated_at
		FROM accounts a
		JOIN account_roles ar ON a.id = ar.account_id
		JOIN roles r ON ar.role_id = r.id
		WHERE a.email = $1 AND r.name = $2
		LIMIT 1
	`

	var model model.Account
	if err := a.db.GetContext(ctx, &model, query, email.String(), role); err != nil {
		return domain.Account{}, db_error.HandlePgError(err)
	}

	acc := domain.Account{
		ID:        model.ID,
		Username:  model.Username,
		Firstname: model.Firstname,
		Lastname:  model.Lastname,
		Password:  model.Password,
		Email:     model.Email,
		AvatarURL: model.AvatarURL,
		Active:    model.Active,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}

	return acc, nil
}
