package companyresourcetypeport

import (
	"context"

	"github.com/google/uuid"
)

type CreateRequest struct {
	InitiatorID, CompanyID uuid.UUID
	Name, IconKey          string
}

type CreateResponse struct {
	ID uuid.UUID
}

type Create interface {
	Execute(context.Context, *CreateRequest) (*CreateResponse, error)
}
