package resource

import (
	"net/http"
	"strconv"
	"time4book/internal/app/adapters/in/http/handlers"
	resourcecommands "time4book/internal/app/core/usecases/resource"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ListResponse struct {
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

type PaginatedResourceResponse struct {
	Status bool           `json:"status"`
	Data   []ListResponse `json:"data"`
	Total  int64          `json:"total"`
	Page   int            `json:"page"`
	Limit  int            `json:"limit"`
}

// List godoc
// @Summary      List resources
// @Description  Get a paginated list of resources
// @Tags         resources
// @Produce      json
// @Param        page query int false "Page number"
// @Param        limit query int false "Items per page"
// @Param        companyId query string false "Company ID"
// @Param        search query string false "Search term"
// @Param        type query string false "Type"
// @Param        status query string false "Status"
// @Success      200  {object}  PaginatedResourceResponse
// @Failure      500  {object}  handlers.ErrorResponse
// @Security     BearerAuth
// @Router       /resources [get]
func (h *Handler) List(c *gin.Context) {
	req := &resourcecommands.ListRequest{
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

	if compID := c.Query("companyId"); compID != "" {
		if id, err := uuid.Parse(compID); err == nil {
			req.CompanyID = id
		}
	} else if compContext, exists := c.Get("companyID"); exists {
		if cid, ok := compContext.(uuid.UUID); ok {
			req.CompanyID = cid
		}
	}

	if search := c.Query("search"); search != "" {
		req.Search = &search
	}
	if typeStr := c.Query("type"); typeStr != "" {
		req.Type = &typeStr
	}
	if status := c.Query("status"); status != "" {
		req.Status = &status
	}

	res, err := h.commands.List.Execute(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, handlers.ErrorResponse{
			Status: false,
			Error:  err.Error(),
		})
		return
	}

	items := make([]ListResponse, len(res.Resources))
	for i, r := range res.Resources {
		items[i] = ListResponse{
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
		}
	}

	c.JSON(http.StatusOK, PaginatedResourceResponse{
		Status: true,
		Data:   items,
		Total:  res.Total,
		Page:   req.Page,
		Limit:  req.Limit,
	})
}
