package companycommands

import (
	"context"
	"fmt"
	"log/slog"
	"time4book/internal/app/core/domain/model/company"
	port "time4book/internal/app/core/ports/in/company"
	outports "time4book/internal/app/core/ports/out"
	companyout "time4book/internal/app/core/ports/out/company"
	"time4book/pkg/validator"
)

type Create struct {
	companyRepo companyout.Repository
	tx          outports.TxManager
	validator   *validator.Facade
	log         *slog.Logger
}

func newCreate(
	crepo companyout.Repository,
	tx outports.TxManager,
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

func (c *Create) Execute(ctx context.Context, req *port.CreateRequest) (*port.CreateResponse, error) {
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

	return &port.CreateResponse{
		CompanyID: comp.ID(),
	}, nil
}
