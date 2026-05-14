package postgres

import (
	"context"
	"database/sql"
	"embed"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed sql/*.sql
var migrationsFS embed.FS

const migrationsDir = "sql"

func RunMigrations(ctx context.Context, dbc *Config) error {
	db, err := sql.Open("pgx", dbc.DSN())
	if err != nil {
		return fmt.Errorf("open migration database connection: %w", err)
	}
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("ping migration database connection: %w", err)
	}

	goose.SetBaseFS(migrationsFS)
	defer goose.SetBaseFS(nil)

	if err := goose.SetDialect("pgx"); err != nil {
		return fmt.Errorf("set goose dialect: %w", err)
	}

	if err := goose.UpContext(ctx, db, migrationsDir); err != nil {
		return fmt.Errorf("run database migrations: %w", err)
	}

	return nil
}
