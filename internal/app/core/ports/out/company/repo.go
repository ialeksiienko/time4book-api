package companyport

import (
	"context"
	domain "time4book/internal/app/core/domain/model/company"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, c *domain.Company) error
	Retrieve(ctx context.Context, id uuid.UUID) (*domain.Company, error)
	Get(ctx context.Context, page, limit int) ([]*domain.Company, int64, error)
	Update(ctx context.Context, c *domain.Company) error
}
