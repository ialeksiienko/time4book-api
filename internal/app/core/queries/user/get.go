package userqueries

import (
	"context"
	"fmt"
	"log/slog"
	"time4book/internal/app/core/domain/model/user"
	port "time4book/internal/app/core/ports/in/user"
	userout "time4book/internal/app/core/ports/out/user"
)

type Get struct {
	userRepo userout.Repository
	log      *slog.Logger
}

func newGet(
	urepo userout.Repository,
	l *slog.Logger,
) *Get {
	return &Get{
		userRepo: urepo,
		log:      l,
	}
}

func (c *Get) Execute(ctx context.Context, req *port.GetRequest) (*port.GetResponse, error) {
	filter := userout.Filter{
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

	users, total, err := c.userRepo.Get(ctx, filter)
	if err != nil {
		c.log.Error("list users", slog.String("error", err.Error()))
		return nil, fmt.Errorf("list users: %w", err)
	}

	res := make([]user.User, len(users))
	for i, u := range users {
		res[i] = *u
	}

	return &port.GetResponse{
		Users: res,
		Total: total,
	}, nil
}
