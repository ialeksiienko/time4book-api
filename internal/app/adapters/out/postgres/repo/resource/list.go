package resourcerepo

import (
	"context"
	"fmt"
	"strings"
	"time"
	"time4book/internal/app/adapters/out/postgres"
	"time4book/internal/app/core/domain/model/resource"

	"github.com/google/uuid"
)

func (r *ResourceRepo) List(ctx context.Context, f resource.ListFilter) ([]*resource.Resource, int64, error) {
	queryBuilder := strings.Builder{}
	countBuilder := strings.Builder{}
	args := []interface{}{}
	argId := 1

	queryBuilder.WriteString(`SELECT id, company_id, name, type, description, location, max_reservation_minutes, available_from, available_to, status, unavailable_from, unavailable_to, unavailable_reason, created_at, updated_at FROM resources WHERE company_id = $1`)
	countBuilder.WriteString(`SELECT count(id) FROM resources WHERE company_id = $1`)
	args = append(args, f.CompanyID)
	argId++

	if f.Type != nil {
		whereClause := fmt.Sprintf(` AND type = $%d`, argId)
		queryBuilder.WriteString(whereClause)
		countBuilder.WriteString(whereClause)
		args = append(args, f.Type.String())
		argId++
	}

	if f.Status != nil {
		whereClause := fmt.Sprintf(` AND status = $%d`, argId)
		queryBuilder.WriteString(whereClause)
		countBuilder.WriteString(whereClause)
		args = append(args, f.Status.String())
		argId++
	}

	if f.Search != nil && *f.Search != "" {
		searchTerm := "%" + *f.Search + "%"
		whereClause := fmt.Sprintf(` AND (name ILIKE $%d OR description ILIKE $%d OR location ILIKE $%d)`, argId, argId, argId)
		queryBuilder.WriteString(whereClause)
		countBuilder.WriteString(whereClause)
		args = append(args, searchTerm)
		argId++
	}

	queryBuilder.WriteString(fmt.Sprintf(` ORDER BY created_at DESC LIMIT $%d OFFSET $%d`, argId, argId+1))
	qArgs := append(args, f.Limit, (f.Page-1)*f.Limit)

	q := queryBuilder.String()
	countQ := countBuilder.String()

	var total int64
	err := postgres.ExtractQuerier(ctx, r.db).QueryRow(ctx, countQ, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count resources: %w", err)
	}

	rows, err := postgres.ExtractQuerier(ctx, r.db).Query(ctx, q, qArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("query resources: %w", err)
	}
	defer rows.Close()

	var res []*resource.Resource
	for rows.Next() {
		var row struct {
			ID                    uuid.UUID
			CompanyID             uuid.UUID
			Name                  string
			Type                  string
			Description           string
			Location              string
			MaxReservationMinutes *int
			AvailableFrom         *string
			AvailableTo           *string
			Status                string
			UnavailableFrom       *time.Time
			UnavailableTo         *time.Time
			UnavailableReason     *string
			CreatedAt             time.Time
			UpdatedAt             time.Time
		}

		if err := rows.Scan(
			&row.ID,
			&row.CompanyID,
			&row.Name,
			&row.Type,
			&row.Description,
			&row.Location,
			&row.MaxReservationMinutes,
			&row.AvailableFrom,
			&row.AvailableTo,
			&row.Status,
			&row.UnavailableFrom,
			&row.UnavailableTo,
			&row.UnavailableReason,
			&row.CreatedAt,
			&row.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("scan resource: %w", err)
		}

		res = append(res, resource.Reconstitute(&resource.Props{
			ID:                    row.ID,
			CompanyID:             row.CompanyID,
			Name:                  row.Name,
			ResourceType:          row.Type,
			Description:           row.Description,
			Location:              row.Location,
			MaxReservationMinutes: row.MaxReservationMinutes,
			AvailableFrom:         row.AvailableFrom,
			AvailableTo:           row.AvailableTo,
			Status:                row.Status,
			UnavailableFrom:       row.UnavailableFrom,
			UnavailableTo:         row.UnavailableTo,
			UnavailableReason:     row.UnavailableReason,
			CreatedAt:             row.CreatedAt,
			UpdatedAt:             row.UpdatedAt,
		}))
	}

	return res, total, nil
}
