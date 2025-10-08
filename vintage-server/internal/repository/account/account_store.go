package repository

import (
	"context"
	"fmt"
	"vintage-server/internal/domain/account"

	"github.com/jmoiron/sqlx"
)

type AccountStore interface {
	account.AccountRepository
	ExecTx(ctx context.Context, fn func(account.AccountRepository) error) error
}

type accountSqlStore struct {
	db *sqlx.DB
	account.AccountRepository
}

func NewAccountStore(db *sqlx.DB) AccountStore {
	return &accountSqlStore{
		db:                db,
		AccountRepository: NewAccountRepository(db),
	}
}

func (s *accountSqlStore) ExecTx(ctx context.Context, fn func(account.AccountRepository) error) error {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	repoTx := s.AccountRepository.WithTx(tx)
	if err := fn(repoTx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rollback err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}
