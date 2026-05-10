package authcommands

import (
	"context"
	"fmt"
	"log/slog"
	"time4book/internal/app/core/domain/model/auth"
	"time4book/internal/app/core/domain/model/company"
	"time4book/internal/app/core/domain/model/user"
	"time4book/internal/app/core/ports"
	"time4book/pkg/validator"

	"github.com/google/uuid"
)

type RegisterRequest struct {
	Firstname       string `validate:"required"`
	Lastname        string `validate:"required"`
	Email           string `validate:"required,email"`
	Password        string `validate:"required"`
	CompanyName     string `validate:"required"`
	CompanyNIP      *string
	CompanyAddress  *string
	CompanyIndustry *string
}

type RegisterResponse struct {
	UserID       uuid.UUID
	CompanyID    uuid.UUID
	AccessToken  *ports.Token
	RefreshToken *ports.Token
}

type Register struct {
	userRepo    user.UserRepo
	authRepo    auth.AuthRepo
	companyRepo company.CompanyRepo

	tx    ports.TxManager
	token ports.TokenManager

	validator *validator.Facade
	log       *slog.Logger
}

func newRegister(
	urepo user.UserRepo,
	arepo auth.AuthRepo,
	crepo company.CompanyRepo,
	tx ports.TxManager,
	token ports.TokenManager,
	v *validator.Facade,
	l *slog.Logger,
) *Register {
	return &Register{
		userRepo:    urepo,
		authRepo:    arepo,
		companyRepo: crepo,
		tx:          tx,
		token:       token,
		validator:   v,
		log:         l,
	}
}

func (r *Register) Execute(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error) {
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

	accessToken, genErr0 := r.token.GenerateToken(usr.ID(), usr.Role().String(), ports.AccessToken)
	if genErr0 != nil {
		return nil, fmt.Errorf("generate access token: %w", genErr0)
	}

	refreshToken, genErr := r.token.GenerateToken(usr.ID(), usr.Role().String(), ports.RefreshToken)
	if genErr != nil {
		return nil, fmt.Errorf("generate access token: %w", genErr)
	}

	session := auth.NewSession(usr.ID(), refreshToken.Value, refreshToken.ExpiresAt)

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
	return &RegisterResponse{
		UserID:       usr.ID(),
		CompanyID:    comp.ID(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
