package authport

import (
	"context"
	"time4book/internal/app/core/domain/model/company"
	"time4book/internal/app/core/domain/model/user"

	"github.com/google/uuid"
)

type RetrieveProfileRequest struct {
	UserID uuid.UUID
}

type RetrieveProfileResponse struct {
	User    *user.User
	Company *company.Company
}

type RetrieveProfile interface {
	Execute(context.Context, *RetrieveProfileRequest) (*RetrieveProfileResponse, error)
}
