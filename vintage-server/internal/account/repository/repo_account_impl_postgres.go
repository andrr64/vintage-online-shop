package repository

import (
	"context"
	"vintage-server/internal/account/domain"
	"vintage-server/internal/account/model"
	"vintage-server/internal/shared/db"
	db_error "vintage-server/pkg"
)

// WithTx implements AccountRepository.
func (a *accountRepoPostgres) WithTx(tx db.DBTX) AccountRepository {
	return &accountRepoPostgres{db: tx}
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
