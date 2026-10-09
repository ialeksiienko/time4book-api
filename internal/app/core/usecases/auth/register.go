package authcommands

import (
	"context"
	"fmt"
	"log/slog"
	"time4book/internal/app/core/domain/model/auth"
	"time4book/internal/app/core/domain/model/company"
	"time4book/internal/app/core/domain/model/user"
	port "time4book/internal/app/core/ports/in/auth"
	outports "time4book/internal/app/core/ports/out"
	authout "time4book/internal/app/core/ports/out/auth"
	companyout "time4book/internal/app/core/ports/out/company"
	userout "time4book/internal/app/core/ports/out/user"
	"time4book/pkg/validator"

	"github.com/google/uuid"
)

type Register struct {
	userRepo    userout.Repository
	authRepo    authout.Repository
	companyRepo companyout.Repository

	tx           outports.TxManager
	tokenManager outports.TokenManager

	validator *validator.Facade
	log       *slog.Logger
}

func newRegister(
	urepo userout.Repository,
	arepo authout.Repository,
	crepo companyout.Repository,
	tx outports.TxManager,
	tokenManager outports.TokenManager,
	v *validator.Facade,
	l *slog.Logger,
) *Register {
	return &Register{
		userRepo:     urepo,
		authRepo:     arepo,
		companyRepo:  crepo,
		tx:           tx,
		tokenManager: tokenManager,
		validator:    v,
		log:          l,
	}
}

func (r *Register) Execute(ctx context.Context, req *port.RegisterRequest) (*port.RegisterResponse, error) {
	err := r.validator.Struct(req)
	if err != nil {
		return nil, fmt.Errorf("validate error: %w", err)
	}

	ownerRole := user.RoleOwnerKey

	role, err0 := user.NewRole(ownerRole, ownerRole.FriendlyName())
	if err0 != nil {
		r.log.Error("new role", slog.String("error", err0.Error()))
		return nil, fmt.Errorf("new role: %w", err0)
	}

	usr, err := user.NewUser(
		req.Firstname,
		req.Lastname,
		req.Email,
		role,
		uuid.Nil, // company id will be set after company creation
	)
	if err != nil {
		r.log.Error("new user", slog.String("error", err.Error()))
		return nil, fmt.Errorf("new user: %w", err)
	}

	creds, cErr := auth.NewCredentials(
		usr.ID(),
		usr.Email(),
		req.Password,
	)
	if cErr != nil {
		r.log.Error("new auth credentials", slog.String("error", cErr.Error()))
		return nil, fmt.Errorf("new auth creds: %w", cErr)
	}

	tokens, err := r.tokenManager.GenerateTokens(usr.ID(), usr.Role().String())
	if err != nil {
		return nil, fmt.Errorf("generate tokens: %w", err)
	}

	session := auth.NewSession(usr.ID(), tokens.RefreshToken.Value, tokens.RefreshToken.ExpiresAt)

	comp, cErr2 := company.NewCompany(
		usr.ID(),
		req.CompanyName,
		req.CompanyNIP,
		req.CompanyAddress,
		req.CompanyIndustry,
	)
	if cErr2 != nil {
		r.log.Error("new company", slog.String("error", cErr2.Error()))
		return nil, fmt.Errorf("new company: %w", cErr2)
	}

	compID := comp.ID()
	usr.SetCompanyID(compID)

	txErr := r.tx.ReadCommitted(ctx, func(txCtx context.Context) error {
		compErr := r.companyRepo.Create(txCtx, comp)
		if compErr != nil {
			r.log.Error("create company", slog.String("error", compErr.Error()))
			return fmt.Errorf("create company: %w", compErr)
		}

		userErr := r.userRepo.Create(txCtx, usr)
		if userErr != nil {
			r.log.Error("create new user", slog.String("error", userErr.Error()))
			return fmt.Errorf("create user: %w", userErr)
		}

		credsErr := r.authRepo.InsertUserCredentials(txCtx, creds)
		if credsErr != nil {
			r.log.Error("create new auth user credentials", slog.String("error", credsErr.Error()))
			return fmt.Errorf("create user credentials: %w", credsErr)
		}

		sessErr := r.authRepo.InsertUserSession(txCtx, session)
		if sessErr != nil {
			r.log.Error("create new auth user session", slog.String("error", sessErr.Error()))
			return fmt.Errorf("create user session: %w", sessErr)
		}

		return nil
	})
	if txErr != nil {
		r.log.Error("transaction", slog.String("error", txErr.Error()))
		return nil, fmt.Errorf("transaction failed: %w", txErr)
	}
	return &port.RegisterResponse{
		UserID:       usr.ID(),
		CompanyID:    comp.ID(),
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}
