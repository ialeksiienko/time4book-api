package companyport

import (
	"context"

	"github.com/google/uuid"
)

type UpdateRequest struct {
	InitiatorID uuid.UUID
	CompanyID   uuid.UUID
	Name        string `validate:"required"`
	NIP         *string
	Address     *string
	Industry    *string
}

type UpdateResponse struct{}

type Update interface {
	Execute(context.Context, *UpdateRequest) (*UpdateResponse, error)
}
