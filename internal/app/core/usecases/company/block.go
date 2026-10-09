package companycommands

import (
	"context"
	"fmt"
	"log/slog"
	"time4book/internal/app/core/domain/model/user"
	port "time4book/internal/app/core/ports/in/company"
	companyout "time4book/internal/app/core/ports/out/company"
	userout "time4book/internal/app/core/ports/out/user"
)

type Block struct {
	userRepo    userout.Repository
	companyRepo companyout.Repository
	log         *slog.Logger
}

func newBlock(
	urepo userout.Repository,
	crepo companyout.Repository,
	l *slog.Logger,
) *Block {
	return &Block{
		userRepo:    urepo,
		companyRepo: crepo,
		log:         l,
	}
}

func (c *Block) Execute(ctx context.Context, req *port.BlockRequest) (*port.BlockResponse, error) {
	initiator, err := c.userRepo.Retrieve(ctx, req.InitiatorID)
	if err != nil {
		return nil, fmt.Errorf("get initiator: %w", err)
	}

	if !initiator.Role().IsDeveloper() {
		return nil, user.ErrUnauthorized
	}

	comp, err := c.companyRepo.Retrieve(ctx, req.CompanyID)
	if err != nil {
		return nil, fmt.Errorf("get company: %w", err)
	}

	comp.Block()

	if err := c.companyRepo.Update(ctx, comp); err != nil {
		return nil, fmt.Errorf("update company: %w", err)
	}

	return &port.BlockResponse{}, nil
}
