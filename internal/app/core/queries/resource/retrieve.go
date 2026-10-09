package resourcequeries

import (
	"context"
	"fmt"
	"log/slog"
	port "time4book/internal/app/core/ports/in/resource"
	resourceout "time4book/internal/app/core/ports/out/resource"
)

type Retrieve struct {
	reader resourceout.Reader
	log    *slog.Logger
}

func newRetrieve(
	reader resourceout.Reader,
	l *slog.Logger,
) *Retrieve {
	return &Retrieve{
		reader: reader,
		log:    l,
	}
}

func (c *Retrieve) Execute(ctx context.Context, req *port.RetrieveRequest) (*port.RetrieveResponse, error) {
	res, err := c.reader.Retrieve(ctx, req.ResourceID, req.CompanyID)
	if err != nil {
		return nil, fmt.Errorf("get resource: %w", err)
	}

	return &port.RetrieveResponse{Resource: res}, nil
}
