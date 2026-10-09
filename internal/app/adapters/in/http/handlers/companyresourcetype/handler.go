package companyresourcetype

import (
	"errors"
	"net/http"
	"time4book/internal/app/adapters/in/http/handlers"
	"time4book/internal/app/core/domain/model/companyresourcetype"
	"time4book/internal/app/core/domain/model/user"
	port "time4book/internal/app/core/ports/in/companyresourcetype"
	companyresourcetypeout "time4book/internal/app/core/ports/out/companyresourcetype"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	commands port.Commands
	queries  port.Queries
}

func NewHandler(commands port.Commands, queries port.Queries) *Handler {
	return &Handler{commands: commands, queries: queries}
}
func userID(c *gin.Context) (uuid.UUID, error) { return uuid.Parse(c.GetString("userID")) }
func companyID(c *gin.Context) (uuid.UUID, error) {
	value := c.GetString("companyID")
	if value == "" {
		return uuid.Nil, nil
	}
	return uuid.Parse(value)
}
func writeError(c *gin.Context, err error) {
	status := http.StatusInternalServerError
	switch {
	case errors.Is(err, user.ErrUnauthorized):
		status = http.StatusForbidden
	case errors.Is(err, companyresourcetypeout.ErrNotFound):
		status = http.StatusNotFound
	case errors.Is(err, companyresourcetypeout.ErrInUse), errors.Is(err, companyresourcetypeout.ErrNameTaken):
		status = http.StatusConflict
	case errors.Is(err, companyresourcetype.ErrInvalidName), errors.Is(err, companyresourcetype.ErrInvalidIconKey):
		status = http.StatusBadRequest
	}
	c.JSON(status, handlers.ErrorResponse{Status: false, Error: err.Error()})
}
