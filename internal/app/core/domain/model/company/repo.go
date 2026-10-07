package company

import (
	"context"

	"github.com/google/uuid"
)

type CompanyRepo interface {
	Create(ctx context.Context, c *Company) error
	ByID(ctx context.Context, id uuid.UUID) (*Company, error)
	List(ctx context.Context, page, limit int) ([]*Company, int64, error)
	Update(ctx context.Context, c *Company) error
}
