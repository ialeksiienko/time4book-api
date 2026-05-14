package resourcerepo

import (
	"context"
	"database/sql"
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
	argID := 1

	baseFrom := ` FROM resources r
LEFT JOIN company_resource_types crt ON crt.id = r.company_resource_type_id
WHERE 1=1`

	queryBuilder.WriteString(`SELECT r.id, r.company_id, r.name, r.type, r.description, r.location, r.max_reservation_minutes, r.available_from, r.available_to, r.status, r.unavailable_from, r.unavailable_to, r.unavailable_reason, r.created_at, r.updated_at,
r.company_resource_type_id, crt.name, crt.icon_key`)
	queryBuilder.WriteString(baseFrom)

	countBuilder.WriteString(`SELECT count(r.id)`)
	countBuilder.WriteString(baseFrom)

	if f.CompanyID != nil {
		whereClause := fmt.Sprintf(` AND r.company_id = $%d`, argID)
		queryBuilder.WriteString(whereClause)
		countBuilder.WriteString(whereClause)
		args = append(args, *f.CompanyID)
		argID++
	}

	if f.Type != nil {
		whereClause := fmt.Sprintf(` AND r.type = $%d`, argID)
		queryBuilder.WriteString(whereClause)
		countBuilder.WriteString(whereClause)
		args = append(args, f.Type.String())
		argID++
	}

	if f.Status != nil {
		whereClause := fmt.Sprintf(` AND r.status = $%d`, argID)
		queryBuilder.WriteString(whereClause)
		countBuilder.WriteString(whereClause)
		args = append(args, f.Status.String())
		argID++
	}

	if f.Search != nil && *f.Search != "" {
		searchTerm := "%" + *f.Search + "%"
		whereClause := fmt.Sprintf(` AND (r.name ILIKE $%d OR r.description ILIKE $%d OR r.location ILIKE $%d)`, argID, argID, argID)
		queryBuilder.WriteString(whereClause)
		countBuilder.WriteString(whereClause)
		args = append(args, searchTerm)
		argID++
	}

	queryBuilder.WriteString(fmt.Sprintf(` ORDER BY r.created_at DESC LIMIT $%d OFFSET $%d`, argID, argID+1))
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
			ID                      uuid.UUID
			CompanyID               uuid.UUID
			Name                    string
			Type                    string
			Description             string
			Location                string
			MaxReservationMinutes   *int
			AvailableFrom           *string
			AvailableTo             *string
			Status                  string
			UnavailableFrom         *time.Time
			UnavailableTo           *time.Time
			UnavailableReason       *string
			CreatedAt               time.Time
			UpdatedAt               time.Time
			CompanyResourceTypeIDNu uuid.NullUUID
			CrtName                 sql.NullString
			CrtIconKey              sql.NullString
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
			&row.CompanyResourceTypeIDNu,
			&row.CrtName,
			&row.CrtIconKey,
		); err != nil {
			return nil, 0, fmt.Errorf("scan resource: %w", err)
		}

		var companyCRTID *uuid.UUID
		if row.CompanyResourceTypeIDNu.Valid {
			u := row.CompanyResourceTypeIDNu.UUID
			companyCRTID = &u
		}
		var customName *string
		if row.CrtName.Valid {
			n := row.CrtName.String
			customName = &n
		}
		var customIcon *string
		if row.CrtIconKey.Valid {
			i := row.CrtIconKey.String
			customIcon = &i
		}

		res = append(res, resource.Reconstitute(&resource.Props{
			ID:                    row.ID,
			CompanyID:             row.CompanyID,
			Name:                  row.Name,
			ResourceType:          row.Type,
			CompanyResourceTypeID: companyCRTID,
			CustomTypeName:        customName,
			CustomTypeIconKey:     customIcon,
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
