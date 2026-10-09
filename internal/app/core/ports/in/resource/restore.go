package resourceport

import (
	"context"

	"github.com/google/uuid"
)

type RestoreRequest struct {
	InitiatorID uuid.UUID
	CompanyID   uuid.UUID
	ResourceID  uuid.UUID
}

type RestoreResponse struct{}

type Restore interface {
	Execute(context.Context, *RestoreRequest) (*RestoreResponse, error)
}
