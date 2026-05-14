package reservation

import (
	"net/http"
	"strconv"
	"time"
	"time4book/internal/app/adapters/in/http/handlers"
	"time4book/internal/app/core/domain/model/user"
	"time4book/internal/app/core/usecases/reservation"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ListResponse struct {
	ID          uuid.UUID  `json:"id"`
	UserID      uuid.UUID  `json:"userId"`
	ResourceID  uuid.UUID  `json:"resourceId"`
	Status      string     `json:"reservationStatus"`
	StartDate   time.Time  `json:"startDate"`
	EndDate     time.Time  `json:"endDate"`
	Description *string    `json:"description,omitempty"`
	CancelledAt *time.Time `json:"cancelledAt,omitempty"`
}

type PaginatedReservationResponse struct {
	Status bool           `json:"status"`
	Data   []ListResponse `json:"data"`
	Total  int64          `json:"total"`
	Page   int            `json:"page"`
	Limit  int            `json:"limit"`
}

// List godoc
// @Summary      List reservations
// @Description  Get a paginated list of reservations globally
// @Tags         reservations
// @Produce      json
// @Param        page query int false "Page number"
// @Param        limit query int false "Items per page"
// @Param        companyId query string false "Company ID"
// @Param        userId query string false "User ID"
// @Param        resourceId query string false "Resource ID"
// @Param        status query string false "Status"
// @Param        from query string false "From date (RFC3339)"
// @Param        to query string false "To date (RFC3339)"
// @Success      200  {object}  PaginatedReservationResponse
// @Failure      500  {object}  handlers.ErrorResponse
// @Security     BearerAuth
// @Router       /reservations [get]
func (h *Handler) List(c *gin.Context) {
	req := &reservationcommands.ListRequest{
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

	role := c.GetString("role")
	var companyScoped *uuid.UUID
	if role == string(user.RoleDeveloperKey) {
		if compID := c.Query("companyId"); compID != "" {
			if id, err := uuid.Parse(compID); err == nil {
				companyScoped = &id
			}
		} else if compContext, exists := c.Get("companyID"); exists {
			if cid, ok := compContext.(uuid.UUID); ok {
				companyScoped = &cid
			}
		}
	} else if compContext, exists := c.Get("companyID"); exists {
		if cid, ok := compContext.(uuid.UUID); ok {
			companyScoped = &cid
		}
	}
	req.CompanyID = companyScoped

	if userID := c.Query("userId"); userID != "" {
		if id, err := uuid.Parse(userID); err == nil {
			req.UserID = &id
		}
	}

	if resourceID := c.Query("resourceId"); resourceID != "" {
		if id, err := uuid.Parse(resourceID); err == nil {
			req.ResourceID = &id
		}
	}

	if status := c.Query("status"); status != "" {
		req.Status = &status
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

	res, err := h.commands.List.Execute(c.Request.Context(), req)
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

	c.JSON(http.StatusOK, PaginatedReservationResponse{
		Status: true,
		Data:   items,
		Total:  res.Total,
		Page:   req.Page,
		Limit:  req.Limit,
	})
}
