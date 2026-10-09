package resourceport

import (
	"context"
	"time4book/internal/app/core/queries"

	"github.com/google/uuid"
)

type RetrieveRequest struct {
	ResourceID uuid.UUID
	CompanyID  uuid.UUID
}

type RetrieveResponse struct {
	Resource *queries.Resource
}

type Retrieve interface {
	Execute(context.Context, *RetrieveRequest) (*RetrieveResponse, error)
}
