package resourcecommands

import (
	"context"
	"fmt"
	"log/slog"
	"time4book/internal/app/core/domain/model/resource"

	"github.com/google/uuid"
)

type ListRequest struct {
	CompanyID uuid.UUID
	Page      int
	Limit     int
	Type      *string
	Status    *string
	Search    *string
}

type ListResponse struct {
	Resources []resource.Resource
	Total     int64
}

type List struct {
	resourceRepo resource.ResourceRepo
	log          *slog.Logger
}

func newList(
	resrepo resource.ResourceRepo,
	l *slog.Logger,
) *List {
	return &List{
		resourceRepo: resrepo,
		log:          l,
	}
}

func (c *List) Execute(ctx context.Context, req *ListRequest) (*ListResponse, error) {
	filter := resource.ListFilter{
		CompanyID: req.CompanyID,
		Page:      req.Page,
		Limit:     req.Limit,
		Search:    req.Search,
	}

	if req.Type != nil && *req.Type != "" {
		t := resource.ResourceType(*req.Type)
		filter.Type = &t
	}

	if req.Status != nil && *req.Status != "" {
		s := resource.ResourceStatus(*req.Status)
		filter.Status = &s
	}

	res, total, err := c.resourceRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("list resources: %w", err)
	}

	resources := make([]resource.Resource, len(res))
	for i, r := range res {
		resources[i] = *r
	}

	return &ListResponse{
		Resources: resources,
		Total:     total,
	}, nil
}
