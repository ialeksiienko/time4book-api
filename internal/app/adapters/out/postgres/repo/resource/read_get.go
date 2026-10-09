package resourcerepo

import (
	"context"
	"fmt"
	"strings"
	"time4book/internal/app/adapters/out/postgres"
	resourceout "time4book/internal/app/core/ports/out/resource"
	"time4book/internal/app/core/queries"

	"github.com/google/uuid"
)

func (r *Reader) Get(ctx context.Context, filter resourceout.Filter) ([]queries.Resource, int64, error) {
	if filter.CompanyID == uuid.Nil {
		return nil, 0, fmt.Errorf("company scope required")
	}

	where := strings.Builder{}
	where.WriteString(" WHERE r.company_id = $1")
	args := []interface{}{filter.CompanyID}
	argID := 2

	if filter.Type != nil && *filter.Type != "" {
		where.WriteString(fmt.Sprintf(" AND r.type = $%d", argID))
		args = append(args, *filter.Type)
		argID++
	}
	if filter.Status != nil && *filter.Status != "" {
		where.WriteString(fmt.Sprintf(" AND r.status = $%d", argID))
		args = append(args, *filter.Status)
		argID++
	}
	if filter.Search != nil && *filter.Search != "" {
		where.WriteString(fmt.Sprintf(" AND (r.name ILIKE $%d OR r.description ILIKE $%d OR r.location ILIKE $%d)", argID, argID, argID))
		args = append(args, "%"+*filter.Search+"%")
		argID++
	}

	querier := postgres.ExtractQuerier(ctx, r.db)
	var total int64
	countQuery := "SELECT count(r.id) FROM resources r" + where.String()
	if err := querier.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count resources: %w", err)
	}

	query := selectResources + where.String() + fmt.Sprintf(" ORDER BY r.created_at DESC LIMIT $%d OFFSET $%d", argID, argID+1)
	queryArgs := append(args, filter.Limit, (filter.Page-1)*filter.Limit)
	rows, err := querier.Query(ctx, query, queryArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("query resources: %w", err)
	}
	defer rows.Close()

	resources := make([]queries.Resource, 0)
	for rows.Next() {
		item, err := scanResource(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("scan resource: %w", err)
		}
		resources = append(resources, *item)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate resources: %w", err)
	}
	return resources, total, nil
}
