package companyport

import (
	"context"

	"github.com/google/uuid"
)

type BlockRequest struct {
	InitiatorID uuid.UUID
	CompanyID   uuid.UUID
}

type BlockResponse struct{}

type Block interface {
	Execute(context.Context, *BlockRequest) (*BlockResponse, error)
}
