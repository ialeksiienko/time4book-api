package resource

import (
	"net/http"
	"time4book/internal/app/adapters/in/http/handlers"
	"time4book/internal/app/core/usecases/resource"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GetByIDResponse struct {
	Status                bool      `json:"status"`
	ID                    uuid.UUID `json:"id"`
	CompanyID             uuid.UUID `json:"companyId"`
	Name                  string    `json:"name"`
	Type                  string    `json:"type"`
	Description           string    `json:"description"`
	Location              string    `json:"location"`
	MaxReservationMinutes *int      `json:"maxReservationMinutes,omitempty"`
	AvailableFrom         *string   `json:"availableFrom,omitempty"`
	AvailableTo           *string   `json:"availableTo,omitempty"`
	ResourceStatus        string    `json:"resourceStatus"`
}

// GetByID godoc
// @Summary      Get resource by ID
// @Description  Get resource details by ID
// @Tags         resources
// @Produce      json
// @Param        id path string true "Resource ID"
// @Success      200  {object}  GetByIDResponse
// @Failure      400  {object}  handlers.ErrorResponse
// @Failure      500  {object}  handlers.ErrorResponse
// @Security     BearerAuth
// @Router       /resources/{id} [get]
func (h *Handler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, handlers.ErrorResponse{
			Status: false,
			Error:  "invalid resource id",
		})
		return
	}

	req := &resourcecommands.GetByIDRequest{
		ResourceID: id,
	}

	if companyIDStr := c.GetString("companyID"); companyIDStr != "" {
		if cid, err := uuid.Parse(companyIDStr); err == nil {
			req.CompanyID = cid
		}
	}

	res, err := h.commands.GetByID.Execute(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, handlers.ErrorResponse{
			Status: false,
			Error:  err.Error(),
		})
		return
	}

	r := res.Resource
	c.JSON(http.StatusOK, GetByIDResponse{
		Status:                true,
		ID:                    r.ID(),
		CompanyID:             r.CompanyID(),
		Name:                  r.Name(),
		Type:                  r.Type().String(),
		Description:           r.Description(),
		Location:              r.Location(),
		MaxReservationMinutes: r.MaxReservationMinutes(),
		AvailableFrom:         r.AvailableFrom(),
		AvailableTo:           r.AvailableTo(),
		ResourceStatus:        r.Status().String(),
	})
}

