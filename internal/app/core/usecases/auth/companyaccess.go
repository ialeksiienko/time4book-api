package authcommands

import (
	"context"
	"fmt"
	"time4book/internal/app/core/domain/model/company"
	"time4book/internal/app/core/domain/model/user"
)

func ensureCompanyAccess(ctx context.Context, companyRepo company.CompanyRepo, usr *user.User) error {
	if usr.Role().IsDeveloper() || usr.CompanyID() == nil {
		return nil
	}

	comp, err := companyRepo.ByID(ctx, *usr.CompanyID())
	if err != nil {
		return fmt.Errorf("get company: %w", err)
	}

	if comp.IsBlocked() {
		return ErrCompanyBlocked
	}

	return nil
}
