package companyrepo

import (
	"time4book/internal/app/adapters/out/postgres"

	"github.com/jackc/pgx/v4/pgxpool"
)

type CompanyRepo struct {
	db *pgxpool.Pool
}

func New(
	db *postgres.Datastore,
) *CompanyRepo {
	return &CompanyRepo{
		db: db.Pool(),
	}
}
