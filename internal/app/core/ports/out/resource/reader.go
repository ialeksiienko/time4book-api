package resourceport

import (
	"context"
	"time4book/internal/app/core/queries"

	"github.com/google/uuid"
)

type Filter struct {
	CompanyID uuid.UUID
	Search    *string
	Type      *string
	Status    *string
	Page      int
	Limit     int
}

type Reader interface {
	Get(context.Context, Filter) ([]queries.Resource, int64, error)
	Retrieve(context.Context, uuid.UUID, uuid.UUID) (*queries.Resource, error)
}
