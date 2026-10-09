package userrepo

import (
	"context"
	"fmt"
	"time4book/internal/app/adapters/out/postgres"
	"time4book/internal/app/core/domain/model/user"
)

func (r *UserRepo) Create(ctx context.Context, u *user.User) error {
	q := `INSERT INTO users (id, firstname, lastname, email, role_id, company_id, status, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := postgres.ExtractQuerier(ctx, r.db).Exec(ctx, q, u.ID(), u.Firstname(), u.Lastname(), u.Email(), u.Role().Key(), u.CompanyID(), u.Status().String(), u.CreatedAt())
	if err != nil {
		return fmt.Errorf("exec insert user: %w", err)
	}
	return nil
}
