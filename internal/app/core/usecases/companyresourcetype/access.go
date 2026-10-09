package companyresourcetypecommands

import (
	"context"
	"time4book/internal/app/core/domain/model/user"
	userout "time4book/internal/app/core/ports/out/user"

	"github.com/google/uuid"
)

func authorize(ctx context.Context, repo userout.Repository, initiatorID, companyID uuid.UUID, write bool) error {
	initiator, err := repo.Retrieve(ctx, initiatorID)
	if err != nil {
		return err
	}
	if !initiator.IsActive() {
		return user.ErrUnauthorized
	}
	if initiator.Role().IsDeveloper() {
		return nil
	}
	if initiator.CompanyID() != companyID {
		return user.ErrUnauthorized
	}
	if initiator.Role().IsOwner() || initiator.Role().IsAdmin() {
		return nil
	}
	if !write && initiator.Role().IsEmployee() {
		return nil
	}
	return user.ErrUnauthorized
}
