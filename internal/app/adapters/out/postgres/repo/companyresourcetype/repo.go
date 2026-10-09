package companyresourcetyperepo

import (
	"context"
	"fmt"
	"time"
	"time4book/internal/app/adapters/out/postgres"
	"time4book/internal/app/core/domain/model/companyresourcetype"
	companyresourcetypeout "time4book/internal/app/core/ports/out/companyresourcetype"

	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Repo struct {
	db *pgxpool.Pool
}

func New(datastore *postgres.Datastore) *Repo {
	return &Repo{db: datastore.Pool()}
}

func (r *Repo) Create(ctx context.Context, t *companyresourcetype.CompanyResourceType) error {
	q := `INSERT INTO company_resource_types (id, company_id, name, icon_key, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := postgres.ExtractQuerier(ctx, r.db).Exec(ctx, q,
		t.ID(), t.CompanyID(), t.Name(), t.IconKey(), t.CreatedAt(), t.UpdatedAt())
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return companyresourcetypeout.ErrNameTaken
		}
		return fmt.Errorf("insert company_resource_type: %w", err)
	}
	return nil
}

func (r *Repo) Retrieve(ctx context.Context, id uuid.UUID) (*companyresourcetype.CompanyResourceType, error) {
	q := `SELECT id, company_id, name, icon_key, created_at, updated_at
          FROM company_resource_types WHERE id = $1`
	var row struct {
		ID        uuid.UUID
		CompanyID uuid.UUID
		Name      string
		IconKey   string
		CreatedAt time.Time
		UpdatedAt time.Time
	}
	err := postgres.ExtractQuerier(ctx, r.db).QueryRow(ctx, q, id).Scan(
		&row.ID, &row.CompanyID, &row.Name, &row.IconKey, &row.CreatedAt, &row.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, companyresourcetypeout.ErrNotFound
		}
		return nil, fmt.Errorf("scan company resource type: %w", err)
	}
	return companyresourcetype.Reconstitute(row.ID, row.CompanyID, row.Name, row.IconKey, row.CreatedAt, row.UpdatedAt), nil
}

func (r *Repo) Get(ctx context.Context, companyID uuid.UUID) ([]*companyresourcetype.CompanyResourceType, error) {
	q := `SELECT id, company_id, name, icon_key, created_at, updated_at
          FROM company_resource_types WHERE company_id = $1 ORDER BY name ASC`
	rows, err := postgres.ExtractQuerier(ctx, r.db).Query(ctx, q, companyID)
	if err != nil {
		return nil, fmt.Errorf("query company resource types: %w", err)
	}
	defer rows.Close()

	var out []*companyresourcetype.CompanyResourceType
	for rows.Next() {
		var row struct {
			ID        uuid.UUID
			CompanyID uuid.UUID
			Name      string
			IconKey   string
			CreatedAt time.Time
			UpdatedAt time.Time
		}
		if err := rows.Scan(&row.ID, &row.CompanyID, &row.Name, &row.IconKey, &row.CreatedAt, &row.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan company resource type: %w", err)
		}
		out = append(out, companyresourcetype.Reconstitute(row.ID, row.CompanyID, row.Name, row.IconKey, row.CreatedAt, row.UpdatedAt))
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate company resource types: %w", err)
	}
	return out, nil
}

func (r *Repo) Update(ctx context.Context, t *companyresourcetype.CompanyResourceType) error {
	q := `UPDATE company_resource_types
          SET name = $1, icon_key = $2, updated_at = $3
          WHERE id = $4 AND company_id = $5`
	_, err := postgres.ExtractQuerier(ctx, r.db).Exec(ctx, q,
		t.Name(), t.IconKey(), t.UpdatedAt(), t.ID(), t.CompanyID())
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return companyresourcetypeout.ErrNameTaken
		}
		return fmt.Errorf("update company_resource_type: %w", err)
	}
	return nil
}

func (r *Repo) Delete(ctx context.Context, id, companyID uuid.UUID) error {
	q := `DELETE FROM company_resource_types WHERE id = $1 AND company_id = $2`
	_, err := postgres.ExtractQuerier(ctx, r.db).Exec(ctx, q, id, companyID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return companyresourcetypeout.ErrInUse
		}
		return fmt.Errorf("delete company_resource_type: %w", err)
	}
	return nil
}
