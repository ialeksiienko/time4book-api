package resourceport

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type ServiceRequest struct {
	InitiatorID uuid.UUID
	CompanyID   uuid.UUID
	ResourceID  uuid.UUID
	Reason      string    `validate:"required"`
	From        time.Time `validate:"required"`
	To          *time.Time
}

type ServiceResponse struct{}

type Service interface {
	Execute(context.Context, *ServiceRequest) (*ServiceResponse, error)
}
