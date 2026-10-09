package companyresourcetype

import (
	"net/http"
	"time4book/internal/app/adapters/in/http/handlers"
	port "time4book/internal/app/core/ports/in/companyresourcetype"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateRequest struct {
	CompanyID uuid.UUID `json:"companyId" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	IconKey   string    `json:"iconKey" binding:"required"`
}
type CreateResponse struct {
	Status bool      `json:"status"`
	ID     uuid.UUID `json:"id"`
}

func (h *Handler) Create(c *gin.Context) {
	var body CreateRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, handlers.ErrorResponse{Status: false, Error: err.Error()})
		return
	}
	initiatorID, err := userID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, handlers.ErrorResponse{Status: false, Error: "invalid user"})
		return
	}
	scope, err := companyID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, handlers.ErrorResponse{Status: false, Error: "invalid company scope"})
		return
	}
	if scope != uuid.Nil && scope != body.CompanyID {
		c.JSON(http.StatusForbidden, handlers.ErrorResponse{Status: false, Error: "company mismatch"})
		return
	}
	result, err := h.commands.Create.Execute(c.Request.Context(), &port.CreateRequest{InitiatorID: initiatorID, CompanyID: body.CompanyID, Name: body.Name, IconKey: body.IconKey})
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusCreated, CreateResponse{Status: true, ID: result.ID})
}
