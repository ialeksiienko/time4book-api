package companyresourcetype

import (
	"errors"
	"net/http"
	"time4book/internal/app/adapters/in/http/handlers"
	domaincrt "time4book/internal/app/core/domain/model/companyresourcetype"
	companyresourcetypecommands "time4book/internal/app/core/usecases/companyresourcetype"
	"time4book/pkg/validator"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateRequestBody struct {
	CompanyID uuid.UUID `json:"companyId" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	IconKey   string    `json:"iconKey" binding:"required"`
}

type CreateResponseBody struct {
	Status bool      `json:"status"`
	ID     uuid.UUID `json:"id"`
}

func (h *Handler) Create(c *gin.Context) {
	var body CreateRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, handlers.ErrorResponse{Status: false, Error: err.Error()})
		return
	}

	initiatorID, err := uuid.Parse(c.GetString("userID"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, handlers.ErrorResponse{Status: false, Error: "invalid user"})
		return
	}

	scopeCompanyCtx, scoped := c.Get("companyID")
	if scoped {
		if cid, ok := scopeCompanyCtx.(uuid.UUID); ok && cid != uuid.Nil && body.CompanyID != cid {
			c.JSON(http.StatusForbidden, handlers.ErrorResponse{Status: false, Error: "company mismatch"})
			return
		}
	}

	res, err := h.commands.Create.Execute(c.Request.Context(), &companyresourcetypecommands.CreateRequest{
		InitiatorID: initiatorID,
		CompanyID:   body.CompanyID,
		Name:        body.Name,
		IconKey:     body.IconKey,
	})
	if err != nil {
		if errors.Is(err, domaincrt.ErrInvalidIconKey) || errors.Is(err, domaincrt.ErrInvalidName) {
			c.JSON(http.StatusBadRequest, handlers.ErrorResponse{Status: false, Error: err.Error()})
			return
		}
		var validationErr validator.ValidationErrors
		if errors.As(err, &validationErr) {
			c.JSON(http.StatusBadRequest, handlers.ErrorResponse{Status: false, Error: validationErr.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, handlers.ErrorResponse{Status: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, CreateResponseBody{Status: true, ID: res.ID})
}
