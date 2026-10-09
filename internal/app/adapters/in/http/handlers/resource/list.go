package resource

import (
	"net/http"
	"strconv"
	"time4book/internal/app/adapters/in/http/handlers"
	port "time4book/internal/app/core/ports/in/resource"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PaginatedResourceResponse struct {
	Status bool           `json:"status"`
	Data   []ResourceBody `json:"data"`
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
	req := &port.GetRequest{
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

	if compContext := c.GetString("companyID"); compContext != "" {
		req.CompanyID, _ = uuid.Parse(compContext)
	}
	if req.CompanyID == uuid.Nil {
		c.JSON(http.StatusBadRequest, handlers.ErrorResponse{Status: false, Error: "company scope required"})
		return
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

	res, err := h.queries.Get.Execute(c.Request.Context(), req)
	if err != nil {
		writeError(c, err)
		return
	}

	items := make([]ResourceBody, len(res.Resources))
	for i := range res.Resources {
		items[i] = toResourceBody(&res.Resources[i])
	}

	c.JSON(http.StatusOK, PaginatedResourceResponse{
		Status: true,
		Data:   items,
		Total:  res.Total,
		Page:   req.Page,
		Limit:  req.Limit,
	})
}
