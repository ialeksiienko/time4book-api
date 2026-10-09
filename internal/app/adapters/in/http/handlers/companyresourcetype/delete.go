package companyresourcetype

import (
	"net/http"
	"time4book/internal/app/adapters/in/http/handlers"
	port "time4book/internal/app/core/ports/in/companyresourcetype"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, handlers.ErrorResponse{Status: false, Error: "invalid company resource type id"})
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
	_, err = h.commands.Delete.Execute(c.Request.Context(), &port.DeleteRequest{InitiatorID: initiatorID, CompanyID: scope, ID: id})
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, handlers.SuccessResponse{Status: true, Message: "company resource type deleted"})
}
