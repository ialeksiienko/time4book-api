package reservation

import (
	"errors"
	"net/http"
	"time"
	"time4book/internal/app/adapters/in/http/handlers"
	reservationcommands "time4book/internal/app/core/usecases/reservation"
	"time4book/pkg/validator"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateRequest struct {
	ResourceID  uuid.UUID `json:"resourceId" binding:"required"`
	StartDate   time.Time `json:"startDate" binding:"required"`
	EndDate     time.Time `json:"endDate" binding:"required"`
	Description *string   `json:"description,omitempty"`
}

type CreateResponse struct {
	Status        bool      `json:"status"`
	ReservationID uuid.UUID `json:"reservationId"`
}

// Create godoc
// @Summary      Create reservation
// @Description  Create a new reservation for a resource
// @Tags         reservations
// @Accept       json
// @Produce      json
// @Param        request body CreateRequest true "Reservation details"
// @Success      201  {object}  CreateResponse
// @Failure      400  {object}  handlers.ErrorResponse
// @Failure      500  {object}  handlers.ErrorResponse
// @Security     BearerAuth
// @Router       /reservations [post]
func (h *Handler) Create(c *gin.Context) {
	var body CreateRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, handlers.ErrorResponse{
			Status: false,
			Error:  err.Error(),
		})
		return
	}

	initiatorIDStr := c.GetString("userID")
	initiatorID, _ := uuid.Parse(initiatorIDStr)

	companyIDStr := c.GetString("companyID")
	companyID, _ := uuid.Parse(companyIDStr)

	req := &reservationcommands.CreateRequest{
		InitiatorID: initiatorID,
		CompanyID:   companyID,
		ResourceID:  body.ResourceID,
		StartDate:   body.StartDate,
		EndDate:     body.EndDate,
		Description: body.Description,
	}

	res, err := h.commands.Create.Execute(c.Request.Context(), req)
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

	c.JSON(http.StatusCreated, CreateResponse{
		Status:        true,
		ReservationID: res.ReservationID,
	})
}
