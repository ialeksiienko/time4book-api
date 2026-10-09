package companyqueries

import (
	"context"
	"fmt"
	"log/slog"
	"time4book/internal/app/core/domain/model/company"
	port "time4book/internal/app/core/ports/in/company"
	companyout "time4book/internal/app/core/ports/out/company"
)

type Get struct {
	companyRepo companyout.Repository
	log         *slog.Logger
}

func newGet(
	crepo companyout.Repository,
	l *slog.Logger,
) *Get {
	return &Get{
		companyRepo: crepo,
		log:         l,
	}
}

func (c *Get) Execute(ctx context.Context, req *port.GetRequest) (*port.GetResponse, error) {
	companies, total, err := c.companyRepo.Get(ctx, req.Page, req.Limit)
	if err != nil {
		c.log.Error("list companies", slog.String("error", err.Error()))
		return nil, fmt.Errorf("list companies: %w", err)
	}

	res := make([]company.Company, len(companies))
	for i, comp := range companies {
		res[i] = *comp
	}

	return &port.GetResponse{
		Companies: res,
		Total:     total,
	}, nil
}
