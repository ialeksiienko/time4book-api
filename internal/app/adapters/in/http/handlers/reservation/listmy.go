package reservation

import (
	"net/http"
	"strconv"
	"time4book/internal/app/adapters/in/http/handlers"
	"time4book/internal/app/core/usecases/reservation"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ListMyResponse struct {
	Status bool           `json:"status"`
	Data   []ListResponse `json:"data"`
	Total  int64          `json:"total"`
	Page   int            `json:"page"`
	Limit  int            `json:"limit"`
}

// ListMy godoc
// @Summary      List my reservations
// @Description  Get a paginated list of my reservations
// @Tags         reservations
// @Produce      json
// @Param        page query int false "Page number"
// @Param        limit query int false "Items per page"
// @Success      200  {object}  ListMyResponse
// @Failure      500  {object}  handlers.ErrorResponse
// @Security     BearerAuth
// @Router       /reservations/my [get]
func (h *Handler) ListMy(c *gin.Context) {
	req := &reservationcommands.ListMyRequest{
		Page:  1,
		Limit: 20,
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

	userIDStr := c.GetString("userID")
	userID, _ := uuid.Parse(userIDStr)
	req.UserID = userID

	companyIDStr := c.GetString("companyID")
	companyID, _ := uuid.Parse(companyIDStr)
	req.CompanyID = companyID

	res, err := h.commands.ListMy.Execute(c.Request.Context(), req)
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

	c.JSON(http.StatusOK, ListMyResponse{
		Status: true,
		Data:   items,
		Total:  res.Total,
		Page:   req.Page,
		Limit:  req.Limit,
	})
}

