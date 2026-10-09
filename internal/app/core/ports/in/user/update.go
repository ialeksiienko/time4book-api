package userport

import (
	"context"

	"github.com/google/uuid"
)

type UpdateRequest struct {
	InitiatorID uuid.UUID
	CompanyID   uuid.UUID
	TargetID    uuid.UUID
	Firstname   *string
	Lastname    *string
	Role        *string
	Status      *string
}

type UpdateResponse struct{}

type Update interface {
	Execute(context.Context, *UpdateRequest) (*UpdateResponse, error)
}
