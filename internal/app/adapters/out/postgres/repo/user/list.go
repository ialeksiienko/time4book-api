package userrepo

import (
	"context"
	"fmt"
	"strings"
	"time"
	"time4book/internal/app/adapters/out/postgres"
	"time4book/internal/app/core/domain/model/user"

	"github.com/google/uuid"
)

func (r *UserRepo) List(ctx context.Context, f user.ListFilter) ([]*user.User, int64, error) {
	queryBuilder := strings.Builder{}
	countBuilder := strings.Builder{}
	args := []any{}
	argId := 1

	queryBuilder.WriteString(`SELECT id, firstname, lastname, email, role_id, company_id, status, created_at, updated_at FROM users WHERE 1=1`)
	countBuilder.WriteString(`SELECT count(id) FROM users WHERE 1=1`)

	if f.CompanyID != uuid.Nil {
		whereClause := fmt.Sprintf(` AND company_id = $%d`, argId)
		queryBuilder.WriteString(whereClause)
		countBuilder.WriteString(whereClause)
		args = append(args, f.CompanyID)
		argId++
	}

	if f.Status != nil {
		whereClause := fmt.Sprintf(` AND status = $%d`, argId)
		queryBuilder.WriteString(whereClause)
		countBuilder.WriteString(whereClause)
		args = append(args, f.Status.String())
		argId++
	}

	if f.Role != nil {
		whereClause := fmt.Sprintf(` AND role_id = $%d`, argId)
		queryBuilder.WriteString(whereClause)
		countBuilder.WriteString(whereClause)
		args = append(args, string(*f.Role))
		argId++
	}

	if f.Search != nil && *f.Search != "" {
		searchTerm := "%" + *f.Search + "%"
		whereClause := fmt.Sprintf(` AND (firstname ILIKE $%d OR lastname ILIKE $%d OR email ILIKE $%d)`, argId, argId, argId)
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
		return nil, 0, fmt.Errorf("count users: %w", err)
	}

	rows, err := postgres.ExtractQuerier(ctx, r.db).Query(ctx, q, qArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("query users: %w", err)
	}
	defer rows.Close()

	var res []*user.User
	for rows.Next() {
		var row struct {
			ID        uuid.UUID
			Firstname string
			Lastname  string
			Email     string
			RoleKey   string
			CompanyID uuid.UUID
			Status    string
			CreatedAt time.Time
			UpdatedAt *time.Time
		}

		if err := rows.Scan(
			&row.ID,
			&row.Firstname,
			&row.Lastname,
			&row.Email,
			&row.RoleKey,
			&row.CompanyID,
			&row.Status,
			&row.CreatedAt,
			&row.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("scan user: %w", err)
		}

		res = append(res, user.Reconstitute(&user.Props{
			ID:        row.ID,
			CompanyID: row.CompanyID,
			Firstname: row.Firstname,
			Lastname:  row.Lastname,
			Email:     row.Email,
			RoleKey:   row.RoleKey,
			Status:    row.Status,
			CreatedAt: row.CreatedAt,
			UpdatedAt: row.UpdatedAt,
		}))
	}

	return res, total, nil
}
