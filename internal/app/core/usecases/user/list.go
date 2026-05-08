package usercommands

import (
	"context"
	"fmt"
	"log/slog"
	"time4book/internal/app/core/domain/model/user"

	"github.com/google/uuid"
)

type ListRequest struct {
	CompanyID *uuid.UUID
	Search    *string
	Role      *string
	Status    *string
	Page      int
	Limit     int
}

type ListResponse struct {
	Users []user.User
	Total int64
}

type List struct {
	userRepo user.UserRepo
	log      *slog.Logger
}

func newList(
	urepo user.UserRepo,
	l *slog.Logger,
) *List {
	return &List{
		userRepo: urepo,
		log:      l,
	}
}

func (c *List) Execute(ctx context.Context, req *ListRequest) (*ListResponse, error) {
	filter := user.ListFilter{
		CompanyID: req.CompanyID,
		Search:    req.Search,
		Page:      req.Page,
		Limit:     req.Limit,
	}

	if req.Role != nil && *req.Role != "" {
		r := user.RoleKeyFromString(*req.Role)
		filter.Role = &r
	}

	if req.Status != nil && *req.Status != "" {
		s := user.UserStatus(*req.Status)
		filter.Status = &s
	}

	users, total, err := c.userRepo.List(ctx, filter)
	if err != nil {
		c.log.Error("list users", slog.String("error", err.Error()))
		return nil, fmt.Errorf("list users: %w", err)
	}

	res := make([]user.User, len(users))
	for i, u := range users {
		res[i] = *u
	}

	return &ListResponse{
		Users: res,
		Total: total,
	}, nil
}
