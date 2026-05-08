package companycommands

import (
	"context"
	"fmt"
	"log/slog"
	"time4book/internal/app/core/domain/model/company"
)

type ListRequest struct {
	Page  int
	Limit int
}

type ListResponse struct {
	Companies []company.Company
	Total     int64
}

type List struct {
	companyRepo company.CompanyRepo
	log         *slog.Logger
}

func newList(
	crepo company.CompanyRepo,
	l *slog.Logger,
) *List {
	return &List{
		companyRepo: crepo,
		log:         l,
	}
}

func (c *List) Execute(ctx context.Context, req *ListRequest) (*ListResponse, error) {
	companies, total, err := c.companyRepo.List(ctx, req.Page, req.Limit)
	if err != nil {
		c.log.Error("list companies", slog.String("error", err.Error()))
		return nil, fmt.Errorf("list companies: %w", err)
	}

	res := make([]company.Company, len(companies))
	for i, comp := range companies {
		res[i] = *comp
	}

	return &ListResponse{
		Companies: res,
		Total:     total,
	}, nil
}
