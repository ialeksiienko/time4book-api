package companyresourcetypeport

import (
	"context"

	"github.com/google/uuid"
)

type UpdateRequest struct {
	InitiatorID, CompanyID, ID uuid.UUID
	Name, IconKey              string
}

type UpdateResponse struct{}

type Update interface {
	Execute(context.Context, *UpdateRequest) (*UpdateResponse, error)
}
