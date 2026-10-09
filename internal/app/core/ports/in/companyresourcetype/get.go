package companyresourcetypeport

import (
	"context"
	"time4book/internal/app/core/domain/model/companyresourcetype"

	"github.com/google/uuid"
)

type GetRequest struct{ InitiatorID, CompanyID uuid.UUID }

type GetResponse struct {
	Items []*companyresourcetype.CompanyResourceType
}

type Get interface {
	Execute(context.Context, *GetRequest) (*GetResponse, error)
}
