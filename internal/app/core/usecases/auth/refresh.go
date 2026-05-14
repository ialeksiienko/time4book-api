package authcommands

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time4book/internal/app/core/domain/model/auth"
	"time4book/internal/app/core/domain/model/company"
	"time4book/internal/app/core/domain/model/user"
	"time4book/internal/app/core/ports"
	"time4book/pkg/validator"
)

type RefreshRequest struct {
	RefreshToken string `validate:"required"`
}

type RefreshResponse struct {
	AccessToken  *ports.Token
	RefreshToken *ports.Token
}

type Refresh struct {
	userRepo user.UserRepo
	authRepo auth.AuthRepo
	companyRepo company.CompanyRepo

	token ports.TokenManager

	validator *validator.Facade
	log       *slog.Logger
}

func newRefresh(
	userRepo user.UserRepo,
	authRepo auth.AuthRepo,
	companyRepo company.CompanyRepo,
	token ports.TokenManager,
	validator *validator.Facade,
	log *slog.Logger,
) *Refresh {
	return &Refresh{
		userRepo:   userRepo,
		authRepo:   authRepo,
		companyRepo: companyRepo,
		token:      token,
		validator:  validator,
		log:        log,
	}
}

func (r *Refresh) Execute(ctx context.Context, req *RefreshRequest) (*RefreshResponse, error) {
	if err := r.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("validate error: %w", err)
	}

	session, err := r.authRepo.GetSessionByRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, errors.New("invalid or expired refresh token")
	}

	if session.IsExpired() {
		_ = r.authRepo.DeleteSessionByRefreshToken(ctx, req.RefreshToken)
		return nil, errors.New("refresh token expired")
	}

	user, err := r.userRepo.ByID(ctx, session.UserID())
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	if err := ensureCompanyAccess(ctx, r.companyRepo, user); err != nil {
		return nil, err
	}

	accessToken, err := r.token.GenerateToken(user.ID(), user.Role().String(), ports.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	refreshToken, err := r.token.GenerateToken(user.ID(), user.Role().String(), ports.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}

	if err := r.authRepo.DeleteSessionByRefreshToken(ctx, req.RefreshToken); err != nil {
		return nil, fmt.Errorf("delete old session: %w", err)
	}

	newSession := auth.NewSession(user.ID(), refreshToken.Value, refreshToken.ExpiresAt)
	if err := r.authRepo.InsertUserSession(ctx, newSession); err != nil {
		return nil, fmt.Errorf("insert new session: %w", err)
	}

	return &RefreshResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
