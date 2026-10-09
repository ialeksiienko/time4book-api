package userrepo

import (
	"context"
	"fmt"
	"time4book/internal/app/adapters/out/postgres"
	"time4book/internal/app/core/domain/model/user"
)

func (r *UserRepo) Update(ctx context.Context, u *user.User) error {
	q := `UPDATE users SET firstname = $1, lastname = $2, role_id = $3, status = $4, updated_at = $5 WHERE id = $6 AND company_id = $7`

	_, err := postgres.ExtractQuerier(ctx, r.db).Exec(ctx, q, u.Firstname(), u.Lastname(), u.Role().Key(), u.Status().String(), u.UpdatedAt(), u.ID(), u.CompanyID())
	if err != nil {
		return fmt.Errorf("exec update user: %w", err)
	}

	return nil
}
