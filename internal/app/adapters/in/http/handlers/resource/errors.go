package resource

import (
	"errors"
	"net/http"
	"time4book/internal/app/adapters/in/http/handlers"
	"time4book/internal/app/core/domain/model/resource"
	"time4book/internal/app/core/domain/model/user"
	companytypeout "time4book/internal/app/core/ports/out/companyresourcetype"
	resourceout "time4book/internal/app/core/ports/out/resource"
	"time4book/pkg/validator"

	"github.com/gin-gonic/gin"
)

func writeError(c *gin.Context, err error) {
	status := http.StatusInternalServerError
	var validationErr validator.ValidationErrors
	switch {
	case errors.Is(err, resource.ErrInvalidType), errors.Is(err, resource.ErrInvalidDetails), errors.Is(err, resource.ErrInvalidServiceInterval), errors.Is(err, resource.ErrInvalidServiceReason), errors.As(err, &validationErr):
		status = http.StatusBadRequest
	case errors.Is(err, user.ErrUnauthorized):
		status = http.StatusForbidden
	case errors.Is(err, resourceout.ErrNotFound), errors.Is(err, companytypeout.ErrNotFound):
		status = http.StatusNotFound
	}
	c.JSON(status, handlers.ErrorResponse{Status: false, Error: err.Error()})
}
