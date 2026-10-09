package resourceport

import (
	"context"
	"time4book/internal/app/core/queries"

	"github.com/google/uuid"
)

type GetRequest struct {
	CompanyID uuid.UUID
	Page      int
	Limit     int
	Type      *string
	Status    *string
	Search    *string
}

type GetResponse struct {
	Resources []queries.Resource
	Total     int64
}

type Get interface {
	Execute(context.Context, *GetRequest) (*GetResponse, error)
}
