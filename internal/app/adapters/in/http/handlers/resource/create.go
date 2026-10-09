package resource

import (
	"net/http"
	"time4book/internal/app/adapters/in/http/handlers"
	port "time4book/internal/app/core/ports/in/resource"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateRequest struct {
	CompanyID             uuid.UUID  `json:"companyId"`
	Name                  string     `json:"name" binding:"required"`
	Type                  string     `json:"type" binding:"required"`
	CompanyResourceTypeID *uuid.UUID `json:"companyResourceTypeId"`
	Description           string     `json:"description"`
	Location              string     `json:"location"`
	MaxReservationMinutes *int       `json:"maxReservationMinutes"`
	AvailableFrom         *string    `json:"availableFrom"`
	AvailableTo           *string    `json:"availableTo"`
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
	if body.CompanyID != uuid.Nil {
		if companyID != uuid.Nil && body.CompanyID != companyID {
			c.JSON(http.StatusForbidden, handlers.ErrorResponse{Status: false, Error: "company mismatch"})
			return
		}
		companyID = body.CompanyID
	}

	req := &port.CreateRequest{
		InitiatorID:           initiatorID,
		CompanyID:             companyID,
		Name:                  body.Name,
		Type:                  body.Type,
		CompanyResourceTypeID: body.CompanyResourceTypeID,
		Description:           body.Description,
		Location:              body.Location,
		MaxReservationMinutes: body.MaxReservationMinutes,
		AvailableFrom:         body.AvailableFrom,
		AvailableTo:           body.AvailableTo,
	}

	res, err := h.commands.Create.Execute(c.Request.Context(), req)
	if err != nil {
		writeError(c, err)
		return
	}

	c.JSON(http.StatusCreated, CreateResponse{
		Status:     true,
		ResourceID: res.ResourceID,
	})
}
