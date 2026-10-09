package authrepo

import (
	"time4book/internal/app/adapters/out/postgres"

	"github.com/jackc/pgx/v4/pgxpool"
)

type AuthRepo struct {
	db *pgxpool.Pool
}

func New(
	db *postgres.Datastore,
) *AuthRepo {
	return &AuthRepo{
		db: db.Pool(),
	}
}
