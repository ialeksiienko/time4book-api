package companyresourcetype

import (
	"errors"
	"net/http"
	"time4book/internal/app/adapters/http/handlers"
	"time4book/internal/app/core/domain/model/companyresourcetype"
	"time4book/internal/app/core/domain/model/user"
	companyresourcetypecommands "time4book/internal/app/core/usecases/companyresourcetype"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	commands *companyresourcetypecommands.Facade
}

func NewHandler(commands *companyresourcetypecommands.Facade) *Handler {
	return &Handler{commands: commands}
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
	case errors.Is(err, companyresourcetype.ErrNotFound):
		status = http.StatusNotFound
	case errors.Is(err, companyresourcetype.ErrInUse), errors.Is(err, companyresourcetype.ErrNameTaken):
		status = http.StatusConflict
	case errors.Is(err, companyresourcetype.ErrInvalidName), errors.Is(err, companyresourcetype.ErrInvalidIconKey):
		status = http.StatusBadRequest
	}
	c.JSON(status, handlers.ErrorResponse{Status: false, Error: err.Error()})
}
