package companyresourcetype

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrNotFound  = errors.New("company resource type not found")
	ErrInUse     = errors.New("company resource type is in use")
	ErrNameTaken = errors.New("company resource type name already exists")
)

type Repo interface {
	Create(ctx context.Context, t *CompanyResourceType) error
	Retrieve(ctx context.Context, id uuid.UUID) (*CompanyResourceType, error)
	Get(ctx context.Context, companyID uuid.UUID) ([]*CompanyResourceType, error)
	Update(ctx context.Context, t *CompanyResourceType) error
	Delete(ctx context.Context, id, companyID uuid.UUID) error
}
