package companyport

import (
	"context"
	"time4book/internal/app/core/domain/model/company"

	"github.com/google/uuid"
)

type RetrieveRequest struct {
	CompanyID uuid.UUID
}

type RetrieveResponse struct {
	Company *company.Company
}

type Retrieve interface {
	Execute(context.Context, *RetrieveRequest) (*RetrieveResponse, error)
}
