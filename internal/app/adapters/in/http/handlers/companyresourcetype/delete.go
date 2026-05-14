package companyresourcetype

import (
	"errors"
	"net/http"
	"time4book/internal/app/adapters/in/http/handlers"
	"time4book/internal/app/core/domain/model/user"
	companyresourcetypecommands "time4book/internal/app/core/usecases/companyresourcetype"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) Delete(c *gin.Context) {
	typeID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, handlers.ErrorResponse{Status: false, Error: "invalid company resource type id"})
		return
	}

	companyIDValue, ok := c.Get("companyID")
	if !ok {
		c.JSON(http.StatusBadRequest, handlers.ErrorResponse{Status: false, Error: "missing company scope"})
		return
	}
	companyID, ok := companyIDValue.(uuid.UUID)
	if !ok || companyID == uuid.Nil {
		c.JSON(http.StatusBadRequest, handlers.ErrorResponse{Status: false, Error: "invalid company scope"})
		return
	}

	initiatorID, err := uuid.Parse(c.GetString("userID"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, handlers.ErrorResponse{Status: false, Error: "invalid user"})
		return
	}

	_, err = h.commands.Delete.Execute(c.Request.Context(), &companyresourcetypecommands.DeleteRequest{
		InitiatorID: initiatorID,
		CompanyID:   companyID,
		ID:          typeID,
	})
	if err != nil {
		if errors.Is(err, user.ErrUnauthorized) {
			c.JSON(http.StatusForbidden, handlers.ErrorResponse{Status: false, Error: err.Error()})
			return
		}
		if errors.Is(err, companyresourcetypecommands.ErrTypeInUse) {
			c.JSON(http.StatusConflict, handlers.ErrorResponse{Status: false, Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, handlers.ErrorResponse{Status: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, handlers.SuccessResponse{Status: true, Message: "company resource type deleted"})
}
