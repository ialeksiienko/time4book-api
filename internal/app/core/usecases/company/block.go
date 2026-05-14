package companycommands

import (
	"context"
	"fmt"
	"log/slog"
	"time4book/internal/app/core/domain/model/company"
	"time4book/internal/app/core/domain/model/user"

	"github.com/google/uuid"
)

type BlockRequest struct {
	InitiatorID uuid.UUID
	CompanyID   uuid.UUID
}

type BlockResponse struct{}

type Block struct {
	userRepo    user.UserRepo
	companyRepo company.CompanyRepo
	log         *slog.Logger
}

func newBlock(
	urepo user.UserRepo,
	crepo company.CompanyRepo,
	l *slog.Logger,
) *Block {
	return &Block{
		userRepo:    urepo,
		companyRepo: crepo,
		log:         l,
	}
}

func (c *Block) Execute(ctx context.Context, req *BlockRequest) (*BlockResponse, error) {
	initiator, err := c.userRepo.ByID(ctx, req.InitiatorID)
	if err != nil {
		return nil, fmt.Errorf("get initiator: %w", err)
	}

	comp, err := c.companyRepo.ByID(ctx, req.CompanyID)
	if err != nil {
		return nil, fmt.Errorf("get company: %w", err)
	}

	if !initiator.Role().IsDeveloper() {
		if initiator.CompanyID() == nil || *initiator.CompanyID() != req.CompanyID {
			return nil, user.ErrUnauthorized
		}
		if !initiator.Role().IsOwner() || initiator.ID() != comp.OwnerID() {
			return nil, user.ErrUnauthorized
		}
	}

	comp.Block()

	if err := c.companyRepo.Update(ctx, comp); err != nil {
		return nil, fmt.Errorf("update company: %w", err)
	}

	return &BlockResponse{}, nil
}
