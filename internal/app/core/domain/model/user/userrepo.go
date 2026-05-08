package user

import (
	"context"

	"github.com/google/uuid"
)

type ListFilter struct {
	CompanyID *uuid.UUID
	Search    *string
	Role      *RoleKey
	Status    *UserStatus
	Page      int
	Limit     int
}

type UserRepo interface {
	Create(ctx context.Context, u *User) error
	ByID(ctx context.Context, id uuid.UUID) (*User, error)
	List(ctx context.Context, f ListFilter) ([]*User, int64, error)
	Update(ctx context.Context, u *User) error
}
