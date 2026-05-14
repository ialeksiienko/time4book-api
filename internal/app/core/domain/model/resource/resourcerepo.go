package resource

import (
	"context"

	"github.com/google/uuid"
)

type ListFilter struct {
	CompanyID *uuid.UUID
	Search    *string
	Type      *ResourceType
	Status    *ResourceStatus
	Page      int
	Limit     int
}

type ResourceRepo interface {
	Create(ctx context.Context, r *Resource) error
	ByID(ctx context.Context, id uuid.UUID) (*Resource, error)
	List(ctx context.Context, f ListFilter) ([]*Resource, int64, error)
	Update(ctx context.Context, r *Resource) error
	Delete(ctx context.Context, id uuid.UUID) error
}
