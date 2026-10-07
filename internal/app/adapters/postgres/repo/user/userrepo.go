package userrepo

import (
	"time4book/internal/app/adapters/postgres"

	"github.com/jackc/pgx/v4/pgxpool"
)

type UserRepo struct {
	db *pgxpool.Pool
}

func New(
	db *postgres.Datastore,
) *UserRepo {
	return &UserRepo{
		db: db.Pool(),
	}
}
