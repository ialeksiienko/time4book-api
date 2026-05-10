package resourcecommands

import (
	"context"
	"fmt"
	"log/slog"
	"time4book/internal/app/core/domain/model/resource"

	"github.com/google/uuid"
)

type GetByIDRequest struct {
	ResourceID uuid.UUID
	CompanyID  uuid.UUID
}

type GetByIDResponse struct {
	Resource *resource.Resource
}

type GetByID struct {
	resourceRepo resource.ResourceRepo
	log          *slog.Logger
}

func newGetByID(
	resrepo resource.ResourceRepo,
	l *slog.Logger,
) *GetByID {
	return &GetByID{
		resourceRepo: resrepo,
		log:          l,
	}
}

func (c *GetByID) Execute(ctx context.Context, req *GetByIDRequest) (*GetByIDResponse, error) {
	res, err := c.resourceRepo.ByID(ctx, req.ResourceID)
	if err != nil {
		return nil, fmt.Errorf("get resource: %w", err)
	}

	if req.CompanyID != uuid.Nil && req.CompanyID != res.CompanyID() {
		return nil, fmt.Errorf("resource not found in company")
	}

	return &GetByIDResponse{Resource: res}, nil
}
