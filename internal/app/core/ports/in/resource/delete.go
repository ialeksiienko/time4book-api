package resourceport

import (
	"context"

	"github.com/google/uuid"
)

type DeleteRequest struct {
	InitiatorID uuid.UUID
	CompanyID   uuid.UUID
	ResourceID  uuid.UUID
}

type DeleteResponse struct{}

type Delete interface {
	Execute(context.Context, *DeleteRequest) (*DeleteResponse, error)
}
