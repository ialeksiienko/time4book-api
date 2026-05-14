package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"time4book/internal/app/core/domain/model/auth"
	"time4book/internal/app/core/domain/model/company"
	"time4book/internal/app/core/domain/model/user"
	"time4book/internal/app/core/ports"
)

const (
	DeveloperBootstrapEmail     = "vladbevl@gmail.com"
	DeveloperBootstrapPassword  = "123856vlad"
	devFirstname                = "Vlad"
	devLastname                 = "Bevl"
	DeveloperSandboxCompanyName = "Sandbox (developer)"
)

// EnsureDeveloperUser creates the standard developer credentials when missing (never in prod env).
func EnsureDeveloperUser(
	ctx context.Context,
	env string,
	passwordPlain string,
	tx ports.TxManager,
	users userrepoIface,
	authR authIface,
	companies companyIface,
	logger *slog.Logger,
) error {
	if env == "prod" {
		return nil
	}
	if passwordPlain == "" {
		logger.Warn("developer bootstrap skipped: empty password")
		return nil
	}

	u, errUser := users.ByEmail(ctx, DeveloperBootstrapEmail)
	if errUser == nil {
		creds, err := auth.NewCredentials(u.ID(), DeveloperBootstrapEmail, passwordPlain)
		if err != nil {
			return fmt.Errorf("bootstrap dev hash: %w", err)
		}
		txErr := tx.ReadCommitted(ctx, func(txCtx context.Context) error {
			return authR.UpsertUserCredentials(txCtx, creds)
		})
		if txErr != nil {
			return fmt.Errorf("bootstrap upsert credentials for existing developer: %w", txErr)
		}
		logger.Info("ensured credentials for developer user", slog.String("email", DeveloperBootstrapEmail))
		return nil
	}
	if !errors.Is(errUser, user.ErrNotFound) {
		return fmt.Errorf("developer bootstrap lookup user: %w", errUser)
	}

	role, err := user.NewRole(user.RoleDeveloperKey, user.RoleDeveloperKey.FriendlyName())
	if err != nil {
		return fmt.Errorf("bootstrap dev role: %w", err)
	}

	usr, err := user.NewUser(devFirstname, devLastname, DeveloperBootstrapEmail, role, nil)
	if err != nil {
		return fmt.Errorf("bootstrap dev user: %w", err)
	}

	comp, err := company.NewCompany(usr.ID(), DeveloperSandboxCompanyName, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("bootstrap dev company: %w", err)
	}
	cid := comp.ID()
	usr.SetCompanyID(&cid)

	creds, err := auth.NewCredentials(usr.ID(), DeveloperBootstrapEmail, passwordPlain)
	if err != nil {
		return fmt.Errorf("bootstrap dev hash: %w", err)
	}

	txErr := tx.ReadCommitted(ctx, func(txCtx context.Context) error {
		if err := companies.Create(txCtx, comp); err != nil {
			return err
		}
		if err := users.Create(txCtx, usr); err != nil {
			return err
		}
		if err := authR.UpsertUserCredentials(txCtx, creds); err != nil {
			return err
		}
		return nil
	})
	if txErr != nil {
		return fmt.Errorf("bootstrap developer transaction: %w", txErr)
	}
	logger.Info("created developer user and sandbox company", slog.String("email", DeveloperBootstrapEmail))
	return nil
}

// Narrow interfaces avoid import cycles between bootstrap and repos.
type userrepoIface interface {
	ByEmail(ctx context.Context, email string) (*user.User, error)
	Create(ctx context.Context, u *user.User) error
}

type authIface interface {
	UpsertUserCredentials(ctx context.Context, c *auth.Credentials) error
}

type companyIface interface {
	Create(ctx context.Context, c *company.Company) error
}
