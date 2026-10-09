package companyresourcetypeport

import (
	"context"
	"errors"
	domain "time4book/internal/app/core/domain/model/companyresourcetype"

	"github.com/google/uuid"
)

var (
	ErrNotFound  = errors.New("company resource type not found")
	ErrInUse     = errors.New("company resource type is in use")
	ErrNameTaken = errors.New("company resource type name already exists")
)

type Repository interface {
	Create(ctx context.Context, t *domain.CompanyResourceType) error
	Retrieve(ctx context.Context, id uuid.UUID) (*domain.CompanyResourceType, error)
	Get(ctx context.Context, companyID uuid.UUID) ([]*domain.CompanyResourceType, error)
	Update(ctx context.Context, t *domain.CompanyResourceType) error
	Delete(ctx context.Context, id, companyID uuid.UUID) error
}
