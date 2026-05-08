package companycommands

import (
	"context"
	"fmt"
	"log/slog"
	"time4book/internal/app/core/domain/model/company"

	"github.com/google/uuid"
)

type GetByIDRequest struct {
	CompanyID uuid.UUID
}

type GetByIDResponse struct {
	Company *company.Company
}

type GetByID struct {
	companyRepo company.CompanyRepo
	log         *slog.Logger
}

func newGetByID(
	crepo company.CompanyRepo,
	l *slog.Logger,
) *GetByID {
	return &GetByID{
		companyRepo: crepo,
		log:         l,
	}
}

func (c *GetByID) Execute(ctx context.Context, req *GetByIDRequest) (*GetByIDResponse, error) {
	comp, err := c.companyRepo.ByID(ctx, req.CompanyID)
	if err != nil {
		c.log.Error("get company by id", slog.String("error", err.Error()))
		return nil, fmt.Errorf("get company: %w", err)
	}

	return &GetByIDResponse{
		Company: comp,
	}, nil
}
