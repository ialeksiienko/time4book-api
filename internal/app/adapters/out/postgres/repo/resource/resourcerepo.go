package resourcerepo

import (
	"time4book/internal/app/adapters/out/postgres"

	"github.com/jackc/pgx/v4/pgxpool"
)

type ResourceRepo struct {
	db *pgxpool.Pool
}

func New(
	db *postgres.Datastore,
) *ResourceRepo {
	return &ResourceRepo{
		db: db.Pool(),
	}
}
