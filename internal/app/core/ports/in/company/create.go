package companyport

import (
	"context"

	"github.com/google/uuid"
)

type CreateRequest struct {
	OwnerID  uuid.UUID `validate:"required"`
	Name     string    `validate:"required"`
	NIP      *string
	Address  *string
	Industry *string
}

type CreateResponse struct {
	CompanyID uuid.UUID
}

type Create interface {
	Execute(context.Context, *CreateRequest) (*CreateResponse, error)
}
