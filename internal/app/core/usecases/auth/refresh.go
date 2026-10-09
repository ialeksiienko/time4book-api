package authcommands

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time4book/internal/app/core/domain/model/auth"
	port "time4book/internal/app/core/ports/in/auth"
	outports "time4book/internal/app/core/ports/out"
	authout "time4book/internal/app/core/ports/out/auth"
	userout "time4book/internal/app/core/ports/out/user"
	"time4book/pkg/validator"
)

type Refresh struct {
	userRepo userout.Repository
	authRepo authout.Repository

	tokenManager outports.TokenManager

	validator *validator.Facade
	log       *slog.Logger
}

func newRefresh(
	userRepo userout.Repository,
	authRepo authout.Repository,
	tokenManager outports.TokenManager,
	validator *validator.Facade,
	log *slog.Logger,
) *Refresh {
	return &Refresh{
		userRepo:     userRepo,
		authRepo:     authRepo,
		tokenManager: tokenManager,
		validator:    validator,
		log:          log,
	}
}

func (r *Refresh) Execute(ctx context.Context, req *port.RefreshRequest) (*port.RefreshResponse, error) {
	if err := r.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("validate error: %w", err)
	}

	session, err := r.authRepo.RetrieveSessionByRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, errors.New("invalid or expired refresh token")
	}

	if session.IsExpired() {
		_ = r.authRepo.DeleteSessionByRefreshToken(ctx, req.RefreshToken)
		return nil, errors.New("refresh token expired")
	}

	user, err := r.userRepo.Retrieve(ctx, session.UserID())
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	tokens, err := r.tokenManager.GenerateTokens(user.ID(), user.Role().String())
	if err != nil {
		return nil, fmt.Errorf("generate tokens: %w", err)
	}

	if err := r.authRepo.DeleteSessionByRefreshToken(ctx, req.RefreshToken); err != nil {
		return nil, fmt.Errorf("delete old session: %w", err)
	}

	newSession := auth.NewSession(user.ID(), tokens.RefreshToken.Value, tokens.RefreshToken.ExpiresAt)
	if err := r.authRepo.InsertUserSession(ctx, newSession); err != nil {
		return nil, fmt.Errorf("insert new session: %w", err)
	}

	return &port.RefreshResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}
