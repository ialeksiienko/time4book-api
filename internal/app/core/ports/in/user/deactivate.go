package userport

import (
	"context"

	"github.com/google/uuid"
)

type DeactivateRequest struct {
	InitiatorID uuid.UUID
	CompanyID   uuid.UUID
	TargetID    uuid.UUID
}

type DeactivateResponse struct{}

type Deactivate interface {
	Execute(context.Context, *DeactivateRequest) (*DeactivateResponse, error)
}
