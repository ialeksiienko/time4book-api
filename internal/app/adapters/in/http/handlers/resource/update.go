package resource

import (
	"errors"
	"net/http"
	"time4book/internal/app/adapters/in/http/handlers"
	"time4book/internal/app/core/usecases/resource"
	"time4book/pkg/validator"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UpdateRequest struct {
	Name                  string  `json:"name" binding:"required"`
	Type                  string  `json:"type" binding:"required"`
	Description           string  `json:"description"`
	Location              string  `json:"location"`
	MaxReservationMinutes *int    `json:"maxReservationMinutes"`
	AvailableFrom         *string `json:"availableFrom"`
	AvailableTo           *string `json:"availableTo"`
}

// Update godoc
// @Summary      Update resource
// @Description  Update resource details
// @Tags         resources
// @Accept       json
// @Produce      json
// @Param        id path string true "Resource ID"
// @Param        request body UpdateRequest true "Update parameters"
// @Success      200  {object}  handlers.SuccessResponse
// @Failure      400  {object}  handlers.ErrorResponse
// @Failure      500  {object}  handlers.ErrorResponse
// @Security     BearerAuth
// @Router       /resources/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	var body UpdateRequest

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

	req := &resourcecommands.UpdateRequest{
		InitiatorID:           initiatorID,
		ResourceID:            id,
		Name:                  body.Name,
		Type:                  body.Type,
		Description:           body.Description,
		Location:              body.Location,
		MaxReservationMinutes: body.MaxReservationMinutes,
		AvailableFrom:         body.AvailableFrom,
		AvailableTo:           body.AvailableTo,
	}

	_, err = h.commands.Update.Execute(c.Request.Context(), req)
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
		Message: "resource updated",
	})
}

