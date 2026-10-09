package companyqueries

import (
	"context"
	"fmt"
	"log/slog"
	port "time4book/internal/app/core/ports/in/company"
	companyout "time4book/internal/app/core/ports/out/company"
)

type Retrieve struct {
	companyRepo companyout.Repository
	log         *slog.Logger
}

func newRetrieve(
	crepo companyout.Repository,
	l *slog.Logger,
) *Retrieve {
	return &Retrieve{
		companyRepo: crepo,
		log:         l,
	}
}

func (c *Retrieve) Execute(ctx context.Context, req *port.RetrieveRequest) (*port.RetrieveResponse, error) {
	comp, err := c.companyRepo.Retrieve(ctx, req.CompanyID)
	if err != nil {
		c.log.Error("get company by id", slog.String("error", err.Error()))
		return nil, fmt.Errorf("get company: %w", err)
	}

	return &port.RetrieveResponse{
		Company: comp,
	}, nil
}
