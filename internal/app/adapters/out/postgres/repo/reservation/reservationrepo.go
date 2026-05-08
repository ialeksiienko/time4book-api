package reservationrepo

import (
	"time4book/internal/app/adapters/out/postgres"

	"github.com/jackc/pgx/v4/pgxpool"
)

type ReservationRepo struct {
	db *pgxpool.Pool
}

func New(
	db *postgres.Datastore,
) *ReservationRepo {
	return &ReservationRepo{
		db: db.Pool(),
	}
}
