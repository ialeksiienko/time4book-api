package userport

import (
	"context"
	"time4book/internal/app/core/domain/model/user"

	"github.com/google/uuid"
)

type GetRequest struct {
	CompanyID uuid.UUID
	Search    *string
	Role      *string
	Status    *string
	Page      int
	Limit     int
}

type GetResponse struct {
	Users []user.User
	Total int64
}

type Get interface {
	Execute(context.Context, *GetRequest) (*GetResponse, error)
}
