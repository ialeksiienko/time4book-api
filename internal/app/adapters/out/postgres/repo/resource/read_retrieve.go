package resourcerepo

import (
	"context"
	"errors"
	"fmt"
	"time4book/internal/app/adapters/out/postgres"
	resourceout "time4book/internal/app/core/ports/out/resource"
	"time4book/internal/app/core/queries"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

func (r *Reader) Retrieve(ctx context.Context, id, companyID uuid.UUID) (*queries.Resource, error) {
	query := selectResources + " WHERE r.id = $1"
	args := []interface{}{id}
	if companyID != uuid.Nil {
		query += " AND r.company_id = $2"
		args = append(args, companyID)
	}
	result, err := scanResource(postgres.ExtractQuerier(ctx, r.db).QueryRow(ctx, query, args...))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, resourceout.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("scan resource: %w", err)
	}
	return result, nil
}
