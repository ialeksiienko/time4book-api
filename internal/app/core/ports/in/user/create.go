package userport

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type CreateRequest struct {
	InitiatorID uuid.UUID
	CompanyID   uuid.UUID
	Firstname   string `validate:"required"`
	Lastname    string `validate:"required"`
	Email       string `validate:"required,email"`
	Password    string `validate:"required"`
	Role        string `validate:"required"`
}

type CreateResponse struct {
	UserID uuid.UUID
}

type Create interface {
	Execute(context.Context, *CreateRequest) (*CreateResponse, error)
}

var (
	ErrCompanyMismatch = errors.New("security violation: company mismatch")
	ErrRoleNotAllowed  = errors.New("security violation: role not allowed")
)
