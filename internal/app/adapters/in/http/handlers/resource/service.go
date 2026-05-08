package resource

import (
	"errors"
	"net/http"
	"time"
	"time4book/internal/app/adapters/in/http/handlers"
	"time4book/internal/app/core/usecases/resource"
	"time4book/pkg/validator"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ServiceRequest struct {
	Reason string     `json:"reason" binding:"required"`
	From   time.Time  `json:"from" binding:"required"`
	To     *time.Time `json:"to"`
}

// Service godoc
// @Summary      Mark resource in service
// @Description  Mark a resource as out-of-order/in-service for a period
// @Tags         resources
// @Accept       json
// @Produce      json
// @Param        id path string true "Resource ID"
// @Param        request body ServiceRequest true "Service parameters"
// @Success      200  {object}  handlers.SuccessResponse
// @Failure      400  {object}  handlers.ErrorResponse
// @Failure      500  {object}  handlers.ErrorResponse
// @Security     BearerAuth
// @Router       /resources/{id}/service [post]
func (h *Handler) Service(c *gin.Context) {
	var body ServiceRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, handlers.ErrorResponse{
			Status: false,
			Error:  err.Error(),
		})
		return
	}

	initiatorIDStr := c.GetString("userID")
	initiatorID, _ := uuid.Parse(initiatorIDStr)

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, handlers.ErrorResponse{
			Status: false,
			Error:  "invalid resource id",
		})
		return
	}

	req := &resourcecommands.ServiceRequest{
		InitiatorID: initiatorID,
		ResourceID:  id,
		Reason:      body.Reason,
		From:        body.From,
		To:          body.To,
	}

	_, err = h.commands.Service.Execute(c.Request.Context(), req)
	if err != nil {
		var validationErr validator.ValidationErrors
		if errors.As(err, &validationErr) {
			c.JSON(http.StatusBadRequest, handlers.ErrorResponse{
				Status: false,
				Error:  validationErr.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, handlers.ErrorResponse{
			Status: false,
			Error:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, handlers.SuccessResponse{
		Status:  true,
		Message: "resource marked in service",
	})
}

