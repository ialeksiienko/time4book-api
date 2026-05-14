package usercommands

import (
	"context"
	"fmt"
	"log/slog"
	"time4book/internal/app/core/domain/model/auth"
	"time4book/internal/app/core/domain/model/user"
	"time4book/internal/app/core/ports"
	"time4book/pkg/validator"

	"github.com/google/uuid"
)

type CreateRequest struct {
	InitiatorID uuid.UUID
	CompanyID   uuid.UUID
	Firstname   string `validate:"required"`
	Lastname    string `validate:"required"`
	Email       string `validate:"required,email"`
	Password    string `validate:"required"`
	Role        string `validate:"required"`
}

type CreateResponse struct {
	UserID uuid.UUID
}

type Create struct {
	userRepo  user.UserRepo
	authRepo  auth.AuthRepo
	tx        ports.TxManager
	validator *validator.Facade
	log       *slog.Logger
}

func newCreate(
	urepo user.UserRepo,
	arepo auth.AuthRepo,
	tx ports.TxManager,
	v *validator.Facade,
	l *slog.Logger,
) *Create {
	return &Create{
		userRepo:  urepo,
		authRepo:  arepo,
		tx:        tx,
		validator: v,
		log:       l,
	}
}

func (c *Create) Execute(ctx context.Context, req *CreateRequest) (*CreateResponse, error) {
	if err := c.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("validate error: %w", err)
	}

	initiator, err := c.userRepo.ByID(ctx, req.InitiatorID)
	if err != nil {
		return nil, fmt.Errorf("get initiator: %w", err)
	}

	roleKey := user.RoleKeyFromString(req.Role)
	if !initiator.Role().IsDeveloper() && (initiator.CompanyID() == nil || *initiator.CompanyID() != req.CompanyID) {
		return nil, user.ErrUnauthorized
	}

	targetRole, err := user.NewRole(roleKey, roleKey.FriendlyName())
	if err != nil {
		return nil, fmt.Errorf("invalid role: %w", err)
	}

	if !initiator.CanCreateUserWithRole(targetRole) {
		return nil, user.ErrUnauthorized
	}

	newUser, err := user.NewUser(
		req.Firstname,
		req.Lastname,
		req.Email,
		targetRole,
		&req.CompanyID,
	)
	if err != nil {
		return nil, fmt.Errorf("new user: %w", err)
	}

	creds, err := auth.NewCredentials(newUser.ID(), newUser.Email(), req.Password)
	if err != nil {
		return nil, fmt.Errorf("new auth creds: %w", err)
	}

	txErr := c.tx.ReadCommitted(ctx, func(txCtx context.Context) error {
		if err := c.userRepo.Create(txCtx, newUser); err != nil {
			return fmt.Errorf("create user: %w", err)
		}
		if err := c.authRepo.InsertUserCredentials(txCtx, creds); err != nil {
			return fmt.Errorf("create user credentials: %w", err)
		}
		return nil
	})

	if txErr != nil {
		c.log.Error("transaction failed", slog.String("error", txErr.Error()))
		return nil, txErr
	}

	return &CreateResponse{UserID: newUser.ID()}, nil
}
