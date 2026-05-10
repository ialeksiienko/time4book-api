package authcommands

import (
	"context"
	"fmt"
	"log/slog"
	"time4book/internal/app/core/domain/model/company"
	"time4book/internal/app/core/domain/model/user"

	"github.com/google/uuid"
)

type MeRequest struct {
	UserID uuid.UUID
}

type MeResponse struct {
	User    *user.User
	Company *company.Company
}

type Me struct {
	userRepo    user.UserRepo
	companyRepo company.CompanyRepo
	log         *slog.Logger
}

func newMe(
	urepo user.UserRepo,
	crepo company.CompanyRepo,
	l *slog.Logger,
) *Me {
	return &Me{
		userRepo:    urepo,
		companyRepo: crepo,
		log:         l,
	}
}

func (c *Me) Execute(ctx context.Context, req *MeRequest) (*MeResponse, error) {
	usr, err := c.userRepo.ByID(ctx, req.UserID)
	if err != nil {
		c.log.Error("get user by id", slog.String("error", err.Error()))
		return nil, fmt.Errorf("get user: %w", err)
	}

	var comp *company.Company
	if usr.CompanyID() != uuid.Nil {
		comp, err = c.companyRepo.ByID(ctx, usr.CompanyID())
		if err != nil {
			c.log.Error("get company by id", slog.String("error", err.Error()))
			return nil, fmt.Errorf("get company: %w", err)
		}
	}

	return &MeResponse{
		User:    usr,
		Company: comp,
	}, nil
}
