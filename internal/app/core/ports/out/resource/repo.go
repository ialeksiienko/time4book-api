package resourceport

import (
	"context"
	"errors"
	domain "time4book/internal/app/core/domain/model/resource"

	"github.com/google/uuid"
)

var ErrNotFound = errors.New("resource not found")

type Repository interface {
	Create(ctx context.Context, r *domain.Resource) error
	Retrieve(ctx context.Context, id uuid.UUID) (*domain.Resource, error)
	Update(ctx context.Context, r *domain.Resource) error
	Delete(ctx context.Context, id uuid.UUID) error
}
