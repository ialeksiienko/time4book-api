package userport

import (
	"context"
	domain "time4book/internal/app/core/domain/model/user"

	"github.com/google/uuid"
)

type Filter struct {
	CompanyID uuid.UUID
	Search    *string
	Role      *domain.RoleKey
	Status    *domain.UserStatus
	Page      int
	Limit     int
}

type Repository interface {
	Create(ctx context.Context, u *domain.User) error
	Retrieve(ctx context.Context, id uuid.UUID) (*domain.User, error)
	Get(ctx context.Context, f Filter) ([]*domain.User, int64, error)
	Update(ctx context.Context, u *domain.User) error
}
