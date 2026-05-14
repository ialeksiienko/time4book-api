package companyresourcetype

import (
	"context"

	"github.com/google/uuid"
)

type Repo interface {
	Create(ctx context.Context, t *CompanyResourceType) error
	ByIDAndCompany(ctx context.Context, id, companyID uuid.UUID) (*CompanyResourceType, error)
	ListByCompany(ctx context.Context, companyID uuid.UUID) ([]*CompanyResourceType, error)
	Update(ctx context.Context, t *CompanyResourceType) error
	Delete(ctx context.Context, id, companyID uuid.UUID) error
	CountResourcesUsing(ctx context.Context, id, companyID uuid.UUID) (int64, error)
}
