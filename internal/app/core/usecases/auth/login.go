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

type Login struct {
	userRepo userout.Repository
	authRepo authout.Repository

	tokenManager outports.TokenManager

	validator *validator.Facade
	log       *slog.Logger
}

func newLogin(
	userRepo userout.Repository,
	authRepo authout.Repository,

	tokenManager outports.TokenManager,

	validator *validator.Facade,
	log *slog.Logger,
) *Login {
	return &Login{
		userRepo:     userRepo,
		authRepo:     authRepo,
		tokenManager: tokenManager,
		validator:    validator,
		log:          log,
	}
}

func (l *Login) Execute(ctx context.Context, req *port.LoginRequest) (*port.LoginResponse, error) {
	err := l.validator.Struct(req)
	if err != nil {
		return nil, fmt.Errorf("validate error: %w", err)
	}

	creds, err := l.authRepo.RetrieveCredentialsByEmail(ctx, req.Email)
	if err != nil {
		l.log.Error("get user credentials by email", slog.String("error", err.Error()))
		return nil, errors.New("invalid email or password")
	}

	if !creds.VerifyPassword(req.Password) {
		l.log.Error("password is not correct")
		return nil, errors.New("invalid email or password")
	}

	user, err := l.userRepo.Retrieve(ctx, creds.UserID())
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	tokens, err := l.tokenManager.GenerateTokens(user.ID(), user.Role().String())
	if err != nil {
		return nil, fmt.Errorf("generate tokens: %w", err)
	}

	if err := l.authRepo.DeleteSessionsByUserID(ctx, user.ID()); err != nil {
		return nil, fmt.Errorf("delete session by user id %s: %w", user.ID(), err)
	}

	session := auth.NewSession(user.ID(), tokens.RefreshToken.Value, tokens.RefreshToken.ExpiresAt)

	if err := l.authRepo.InsertUserSession(ctx, session); err != nil {
		return nil, fmt.Errorf("insert user session: %w", err)
	}

	return &port.LoginResponse{
		UserID:       user.ID(),
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}
