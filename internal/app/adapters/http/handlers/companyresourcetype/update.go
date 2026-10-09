package companyresourcetype

import (
	"net/http"
	"time4book/internal/app/adapters/http/handlers"
	companyresourcetypecommands "time4book/internal/app/core/usecases/companyresourcetype"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UpdateRequest struct {
	Name    string `json:"name" binding:"required"`
	IconKey string `json:"iconKey" binding:"required"`
}

func (h *Handler) Update(c *gin.Context) {
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
	var body UpdateRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, handlers.ErrorResponse{Status: false, Error: err.Error()})
		return
	}
	_, err = h.commands.Update.Execute(c.Request.Context(), &companyresourcetypecommands.UpdateRequest{InitiatorID: initiatorID, CompanyID: scope, ID: id, Name: body.Name, IconKey: body.IconKey})
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, handlers.SuccessResponse{Status: true, Message: "company resource type updated"})
}
