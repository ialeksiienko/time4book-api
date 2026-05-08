package userrepo

import (
	"context"
	"fmt"
	"time"
	"time4book/internal/app/adapters/out/postgres"
	"time4book/internal/app/core/domain/model/user"

	"github.com/google/uuid"
)

func (r *UserRepo) ByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	q := `SELECT id, firstname, lastname, email, role_id, company_id, status, created_at, updated_at 
          FROM users WHERE id = $1`

	var row struct {
		ID        uuid.UUID
		Firstname string
		Lastname  string
		Email     string
		RoleKey   string
		CompanyID *uuid.UUID
		Status    string
		CreatedAt time.Time
		UpdatedAt *time.Time
	}

	err := postgres.ExtractQuerier(ctx, r.db).QueryRow(ctx, q, id).Scan(
		&row.ID,
		&row.Firstname,
		&row.Lastname,
		&row.Email,
		&row.RoleKey,
		&row.CompanyID,
		&row.Status,
		&row.CreatedAt,
		&row.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("scan user: %w", err)
	}

	return user.Reconstitute(&user.Props{
		ID:        row.ID,
		CompanyID: row.CompanyID,
		Firstname: row.Firstname,
		Lastname:  row.Lastname,
		Email:     row.Email,
		RoleKey:   row.RoleKey,
		Status:    row.Status,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}), nil
}
