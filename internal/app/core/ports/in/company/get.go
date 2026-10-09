package companyport

import (
	"context"
	"time4book/internal/app/core/domain/model/company"
)

type GetRequest struct {
	Page  int
	Limit int
}

type GetResponse struct {
	Companies []company.Company
	Total     int64
}

type Get interface {
	Execute(context.Context, *GetRequest) (*GetResponse, error)
}
