package resource

import (
	"errors"
	"net/http"
	"time4book/internal/app/adapters/in/http/handlers"
	resourcecommands "time4book/internal/app/core/usecases/resource"
	"time4book/pkg/validator"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateRequest struct {
	Name                  string  `json:"name" binding:"required"`
	Type                  string  `json:"type" binding:"required"`
	Description           string  `json:"description"`
	Location              string  `json:"location"`
	MaxReservationMinutes *int    `json:"maxReservationMinutes"`
	AvailableFrom         *string `json:"availableFrom"`
	AvailableTo           *string `json:"availableTo"`
}

type CreateResponse struct {
	Status     bool      `json:"status"`
	ResourceID uuid.UUID `json:"resourceId"`
}

// Create godoc
// @Summary      Create resource
// @Description  Create a new resource within a company
// @Tags         resources
// @Accept       json
// @Produce      json
// @Param        request body CreateRequest true "Resource details"
// @Success      201  {object}  CreateResponse
// @Failure      400  {object}  handlers.ErrorResponse
// @Failure      500  {object}  handlers.ErrorResponse
// @Security     BearerAuth
// @Router       /resources [post]
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

	req := &resourcecommands.CreateRequest{
		InitiatorID:           initiatorID,
		CompanyID:             companyID,
		Name:                  body.Name,
		Type:                  body.Type,
		Description:           body.Description,
		Location:              body.Location,
		MaxReservationMinutes: body.MaxReservationMinutes,
		AvailableFrom:         body.AvailableFrom,
		AvailableTo:           body.AvailableTo,
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
		Status:     true,
		ResourceID: res.ResourceID,
	})
}
