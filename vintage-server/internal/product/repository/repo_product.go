package repository

import "vintage-server/internal/shared/db"


type ProductRepository interface {

}

type productRepo struct {
	db db.DBTX
}

func NewProductRepoPostgres(db db.DBTX) ProductRepository {
	return &productRepo {
		db: db,
	}
}