package reservationrepo

import (
	"context"
	"fmt"
	"strings"
	"time"
	"time4book/internal/app/adapters/out/postgres"
	"time4book/internal/app/core/domain/model/reservation"

	"github.com/google/uuid"
)

func (r *ReservationRepo) List(ctx context.Context, f reservation.ListFilter) ([]*reservation.Reservation, int64, error) {
	queryBuilder := strings.Builder{}
	countBuilder := strings.Builder{}
	args := []any{}
	argId := 1

	queryBuilder.WriteString(`SELECT r.id, r.user_id, r.company_id, r.resource_id, r.start_date, r.end_date, r.description, r.status, r.created_at, r.updated_at FROM reservations r JOIN resources res ON r.resource_id = res.id WHERE 1=1`)
	countBuilder.WriteString(`SELECT count(r.id) FROM reservations r JOIN resources res ON r.resource_id = res.id WHERE 1=1`)

	if f.CompanyID != nil {
		whereClause := fmt.Sprintf(` AND res.company_id = $%d`, argId)
		queryBuilder.WriteString(whereClause)
		countBuilder.WriteString(whereClause)
		args = append(args, *f.CompanyID)
		argId++
	}

	if f.UserID != nil {
		whereClause := fmt.Sprintf(` AND r.user_id = $%d`, argId)
		queryBuilder.WriteString(whereClause)
		countBuilder.WriteString(whereClause)
		args = append(args, *f.UserID)
		argId++
	}

	if f.ResourceID != nil {
		whereClause := fmt.Sprintf(` AND r.resource_id = $%d`, argId)
		queryBuilder.WriteString(whereClause)
		countBuilder.WriteString(whereClause)
		args = append(args, *f.ResourceID)
		argId++
	}

	if f.Status != nil {
		whereClause := fmt.Sprintf(` AND r.status = $%d`, argId)
		queryBuilder.WriteString(whereClause)
		countBuilder.WriteString(whereClause)
		args = append(args, f.Status.String())
		argId++
	}

	if f.From != nil {
		whereClause := fmt.Sprintf(` AND r.start_date >= $%d`, argId)
		queryBuilder.WriteString(whereClause)
		countBuilder.WriteString(whereClause)
		args = append(args, *f.From)
		argId++
	}

	if f.To != nil {
		whereClause := fmt.Sprintf(` AND r.end_date <= $%d`, argId)
		queryBuilder.WriteString(whereClause)
		countBuilder.WriteString(whereClause)
		args = append(args, *f.To)
		argId++
	}

	queryBuilder.WriteString(fmt.Sprintf(` ORDER BY r.start_date ASC LIMIT $%d OFFSET $%d`, argId, argId+1))
	qArgs := append(args, f.Limit, (f.Page-1)*f.Limit)

	q := queryBuilder.String()
	countQ := countBuilder.String()

	var total int64
	err := postgres.ExtractQuerier(ctx, r.db).QueryRow(ctx, countQ, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count reservations: %w", err)
	}

	rows, err := postgres.ExtractQuerier(ctx, r.db).Query(ctx, q, qArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("query reservations: %w", err)
	}
	defer rows.Close()

	var res []*reservation.Reservation
	for rows.Next() {
		var row struct {
			ID          uuid.UUID
			UserID      uuid.UUID
			CompanyID   uuid.UUID
			ResourceID  uuid.UUID
			StartDate   time.Time
			EndDate     time.Time
			Description *string
			Status      string
			CreatedAt   time.Time
			UpdatedAt   time.Time
		}

		if err := rows.Scan(
			&row.ID,
			&row.UserID,
			&row.CompanyID,
			&row.ResourceID,
			&row.StartDate,
			&row.EndDate,
			&row.Description,
			&row.Status,
			&row.CreatedAt,
			&row.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("scan reservation: %w", err)
		}

		res = append(res, reservation.Reconstitute(&reservation.Props{
			ID:          row.ID,
			UserID:      row.UserID,
			CompanyID:   row.CompanyID,
			ResourceID:  row.ResourceID,
			Description: row.Description,
			StartDate:   row.StartDate,
			EndDate:     row.EndDate,
			Status:      row.Status,
			CreatedAt:   row.CreatedAt,
			UpdatedAt:   row.UpdatedAt,
		}))
	}

	return res, total, nil
}
