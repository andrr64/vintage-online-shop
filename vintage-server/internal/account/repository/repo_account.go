package repository

import (
	"context"
	"vintage-server/internal/account/domain"
	"vintage-server/internal/shared/db"
)

type AccountRepository interface {
	WithTx(tx db.DBTX) AccountRepository

	FindAccountByUsernameAndRoleString(ctx context.Context, username domain.Username, role string) (domain.Account, error)
	FindAccountByEmailAndRoleString(ctx context.Context, email domain.Email, role string) (domain.Account, error)
}

func NewAccountRepositoryPostgres(db db.DBTX) AccountRepository {
	return &accountRepoPostgres{db: db}
}

type accountRepoPostgres struct {
	db db.DBTX
}
