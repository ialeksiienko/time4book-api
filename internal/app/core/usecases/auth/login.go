package authcommands

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time4book/internal/app/core/domain/model/auth"
	"time4book/internal/app/core/domain/model/user"
	"time4book/internal/app/core/ports"
	"time4book/pkg/validator"

	"github.com/google/uuid"
)

type LoginRequest struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
}

type LoginResponse struct {
	UserID       uuid.UUID
	AccessToken  *ports.Token
	RefreshToken *ports.Token
}

type Login struct {
	userRepo user.UserRepo
	authRepo auth.AuthRepo

	tokenManager ports.TokenManager

	validator *validator.Facade
	log       *slog.Logger
}

func newLogin(
	userRepo user.UserRepo,
	authRepo auth.AuthRepo,

	tokenManager ports.TokenManager,

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

func (l *Login) Execute(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	err := l.validator.Struct(req)
	if err != nil {
		return nil, fmt.Errorf("validate error: %w", err)
	}

	creds, err := l.authRepo.GetCredentialsByEmail(ctx, req.Email)
	if err != nil {
		l.log.Error("get user credentials by email", slog.String("error", err.Error()))
		return nil, errors.New("invalid email or password")
	}

	if !creds.VerifyPassword(req.Password) {
		l.log.Error("password is not correct")
		return nil, errors.New("invalid email or password")
	}

	user, err := l.userRepo.ByID(ctx, creds.UserID())
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

	return &LoginResponse{
		UserID:       user.ID(),
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}
