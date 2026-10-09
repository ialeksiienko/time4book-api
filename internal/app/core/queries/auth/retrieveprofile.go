package authqueries

import (
	"context"
	"fmt"
	"log/slog"
	"time4book/internal/app/core/domain/model/company"
	port "time4book/internal/app/core/ports/in/auth"
	companyout "time4book/internal/app/core/ports/out/company"
	userout "time4book/internal/app/core/ports/out/user"

	"github.com/google/uuid"
)

type RetrieveProfile struct {
	userRepo    userout.Repository
	companyRepo companyout.Repository
	log         *slog.Logger
}

func newRetrieveProfile(
	urepo userout.Repository,
	crepo companyout.Repository,
	l *slog.Logger,
) *RetrieveProfile {
	return &RetrieveProfile{
		userRepo:    urepo,
		companyRepo: crepo,
		log:         l,
	}
}

func (c *RetrieveProfile) Execute(ctx context.Context, req *port.RetrieveProfileRequest) (*port.RetrieveProfileResponse, error) {
	usr, err := c.userRepo.Retrieve(ctx, req.UserID)
	if err != nil {
		c.log.Error("get user by id", slog.String("error", err.Error()))
		return nil, fmt.Errorf("get user: %w", err)
	}

	var comp *company.Company
	if usr.CompanyID() != uuid.Nil {
		comp, err = c.companyRepo.Retrieve(ctx, usr.CompanyID())
		if err != nil {
			c.log.Error("get company by id", slog.String("error", err.Error()))
			return nil, fmt.Errorf("get company: %w", err)
		}
	}

	return &port.RetrieveProfileResponse{
		User:    usr,
		Company: comp,
	}, nil
}
