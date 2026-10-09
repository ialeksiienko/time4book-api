package resourcequeries

import (
	"context"
	"fmt"
	"log/slog"
	port "time4book/internal/app/core/ports/in/resource"
	resourceout "time4book/internal/app/core/ports/out/resource"
)

type Get struct {
	reader resourceout.Reader
	log    *slog.Logger
}

func newGet(
	reader resourceout.Reader,
	l *slog.Logger,
) *Get {
	return &Get{
		reader: reader,
		log:    l,
	}
}

func (c *Get) Execute(ctx context.Context, req *port.GetRequest) (*port.GetResponse, error) {
	filter := resourceout.Filter{
		CompanyID: req.CompanyID,
		Page:      req.Page,
		Limit:     req.Limit,
		Search:    req.Search,
	}

	filter.Type = req.Type
	filter.Status = req.Status
	res, total, err := c.reader.Get(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("list resources: %w", err)
	}

	return &port.GetResponse{
		Resources: res,
		Total:     total,
	}, nil
}
