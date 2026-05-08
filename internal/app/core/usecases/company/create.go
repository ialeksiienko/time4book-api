package companycommands

import (
	"context"
	"fmt"
	"log/slog"
	"time4book/internal/app/core/domain/model/company"
	"time4book/internal/app/core/ports"
	"time4book/pkg/validator"

	"github.com/google/uuid"
)

type CreateRequest struct {
	OwnerID  uuid.UUID `validate:"required"`
	Name     string    `validate:"required"`
	NIP      *string
	Address  *string
	Industry *string
}

type CreateResponse struct {
	CompanyID uuid.UUID
}

type Create struct {
	companyRepo company.CompanyRepo
	tx          ports.TxManager
	validator   *validator.Facade
	log         *slog.Logger
}

func newCreate(
	crepo company.CompanyRepo,
	tx ports.TxManager,
	v *validator.Facade,
	l *slog.Logger,
) *Create {
	return &Create{
		companyRepo: crepo,
		tx:          tx,
		validator:   v,
		log:         l,
	}
}

func (c *Create) Execute(ctx context.Context, req *CreateRequest) (*CreateResponse, error) {
	err := c.validator.Struct(req)
	if err != nil {
		return nil, fmt.Errorf("validate error: %w", err)
	}

	comp, err := company.NewCompany(
		req.OwnerID,
		req.Name,
		req.NIP,
		req.Address,
		req.Industry,
	)
	if err != nil {
		c.log.Error("new company", slog.String("error", err.Error()))
		return nil, fmt.Errorf("new company: %w", err)
	}

	err = c.companyRepo.Create(ctx, comp)
	if err != nil {
		c.log.Error("create company repo", slog.String("error", err.Error()))
		return nil, fmt.Errorf("create company: %w", err)
	}

	return &CreateResponse{
		CompanyID: comp.ID(),
	}, nil
}
