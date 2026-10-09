package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Config struct {
	User string
	Pass string
	Host string
	Port int
	Name string

	Logger *slog.Logger
}

type Datastore struct {
	pool *pgxpool.Pool
}

func New(p *pgxpool.Pool) *Datastore {
	return &Datastore{
		pool: p,
	}
}

func (d Datastore) Pool() *pgxpool.Pool {
	return d.pool
}

func (db Config) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		db.User, db.Pass, db.Host, db.Port, db.Name)
}

func NewDBPool(dbc *Config) (*pgxpool.Pool, func(), error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	f := func() {}

	pool, err := pgxpool.Connect(ctx, dbc.DSN())
	if err != nil {
		return nil, f, err
	}

	err = validateDBPool(pool, dbc.Logger)
	if err != nil {
		return nil, f, err
	}

	return pool, func() { pool.Close() }, nil
}

func validateDBPool(p *pgxpool.Pool, l *slog.Logger) error {
	err := p.Ping(context.Background())
	if err != nil {
		return err
	}

	var (
		currentDatabase string
		currentUser     string
		dbVersion       string
	)

	sqlStatement := `select current_database(), current_user, version();`
	row := p.QueryRow(context.Background(), sqlStatement)
	err = row.Scan(&currentDatabase, &currentUser, &dbVersion)

	switch {
	case err == sql.ErrNoRows:
		return errors.New("no rows were returned")
	case err != nil:
		return fmt.Errorf("database connection error: %w", err)
	default:
		l.Debug(fmt.Sprintf("database version: %s\n", dbVersion))
		l.Debug(fmt.Sprintf("current database: %s\n", currentDatabase))
		l.Debug(fmt.Sprintf("current database user: %s\n", currentUser))
	}

	return nil
}
