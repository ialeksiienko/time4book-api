package companyresourcetypeport

import (
	"context"

	"github.com/google/uuid"
)

type DeleteRequest struct{ InitiatorID, CompanyID, ID uuid.UUID }

type DeleteResponse struct{}

type Delete interface {
	Execute(context.Context, *DeleteRequest) (*DeleteResponse, error)
}
