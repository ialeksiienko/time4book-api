package reservation

import (
	"net/http"
	"strconv"
	"time"
	"time4book/internal/app/adapters/in/http/handlers"
	"time4book/internal/app/core/usecases/reservation"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ListByResourceResponse struct {
	Status bool           `json:"status"`
	Data   []ListResponse `json:"data"`
	Total  int64          `json:"total"`
	Page   int            `json:"page"`
	Limit  int            `json:"limit"`
}

// ListByResource godoc
// @Summary      List reservations by resource
// @Description  Get a paginated list of reservations for a specific resource
// @Tags         reservations
// @Produce      json
// @Param        id path string true "Resource ID"
// @Param        page query int false "Page number"
// @Param        limit query int false "Items per page"
// @Param        from query string false "From date (RFC3339)"
// @Param        to query string false "To date (RFC3339)"
// @Success      200  {object}  ListByResourceResponse
// @Failure      400  {object}  handlers.ErrorResponse
// @Failure      500  {object}  handlers.ErrorResponse
// @Security     BearerAuth
// @Router       /reservations/resource/{id} [get]
func (h *Handler) ListByResource(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, handlers.ErrorResponse{
			Status: false,
			Error:  "invalid resource id",
		})
		return
	}

	req := &reservationcommands.ListByResourceRequest{
		ResourceID: id,
		Page:       1,
		Limit:      20,
	}

	if page := c.Query("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			req.Page = p
		}
	}
	if limit := c.Query("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil && l > 0 && l <= 100 {
			req.Limit = l
		}
	}

	if from := c.Query("from"); from != "" {
		if t, err := time.Parse(time.RFC3339, from); err == nil {
			req.From = &t
		}
	}

	if to := c.Query("to"); to != "" {
		if t, err := time.Parse(time.RFC3339, to); err == nil {
			req.To = &t
		}
	}

	res, err := h.commands.ListByResource.Execute(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, handlers.ErrorResponse{
			Status: false,
			Error:  err.Error(),
		})
		return
	}

	items := make([]ListResponse, len(res.Reservations))
	for i, r := range res.Reservations {
		items[i] = ListResponse{
			ID:          r.ID(),
			UserID:      r.UserID(),
			ResourceID:  r.ResourceID(),
			Status:      r.Status().String(),
			StartDate:   r.StartDate(),
			EndDate:     r.EndDate(),
			Description: r.Description(),
			CancelledAt: r.CancelledAt(),
		}
	}

	c.JSON(http.StatusOK, ListByResourceResponse{
		Status: true,
		Data:   items,
		Total:  res.Total,
		Page:   req.Page,
		Limit:  req.Limit,
	})
}
