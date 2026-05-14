package companyresourcetype

import (
	"errors"
	"net/http"
	"time4book/internal/app/adapters/in/http/handlers"
	domaincrt "time4book/internal/app/core/domain/model/companyresourcetype"
	"time4book/internal/app/core/domain/model/user"
	companyresourcetypecommands "time4book/internal/app/core/usecases/companyresourcetype"
	"time4book/pkg/validator"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UpdateRequestBody struct {
	Name    string `json:"name" binding:"required"`
	IconKey string `json:"iconKey" binding:"required"`
}

func (h *Handler) Update(c *gin.Context) {
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

	var body UpdateRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, handlers.ErrorResponse{Status: false, Error: err.Error()})
		return
	}

	_, err = h.commands.Update.Execute(c.Request.Context(), &companyresourcetypecommands.UpdateRequest{
		InitiatorID: initiatorID,
		CompanyID:   companyID,
		ID:          typeID,
		Name:        body.Name,
		IconKey:     body.IconKey,
	})
	if err != nil {
		if errors.Is(err, user.ErrUnauthorized) {
			c.JSON(http.StatusForbidden, handlers.ErrorResponse{Status: false, Error: err.Error()})
			return
		}
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

	c.JSON(http.StatusOK, handlers.SuccessResponse{Status: true, Message: "company resource type updated"})
}
