package company

import (
	"net/http"
	"strconv"
	"time4book/internal/app/adapters/in/http/handlers"
	companycommands "time4book/internal/app/core/usecases/company"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ListResponse struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	NIP           *string   `json:"nip,omitempty"`
	Address       *string   `json:"address,omitempty"`
	Industry      *string   `json:"industry,omitempty"`
	CompanyStatus string    `json:"companyStatus"`
}

type PaginatedCompanyResponse struct {
	Status bool           `json:"status"`
	Data   []ListResponse `json:"data"`
	Total  int64          `json:"total"`
	Page   int            `json:"page"`
	Limit  int            `json:"limit"`
}

// List godoc
// @Summary      List companies
// @Description  Get a paginated list of companies
// @Tags         companies
// @Produce      json
// @Param        page query int false "Page number"
// @Param        limit query int false "Items per page"
// @Param        search query string false "Search term"
// @Param        status query string false "Status"
// @Success      200  {object}  PaginatedCompanyResponse
// @Failure      500  {object}  handlers.ErrorResponse
// @Security     BearerAuth
// @Router       /companies [get]
func (h *Handler) List(c *gin.Context) {
	req := &companycommands.ListRequest{
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

	res, err := h.commands.List.Execute(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, handlers.ErrorResponse{
			Status: false,
			Error:  err.Error(),
		})
		return
	}

	items := make([]ListResponse, len(res.Companies))
	for i, comp := range res.Companies {
		items[i] = ListResponse{
			ID:            comp.ID(),
			Name:          comp.Name(),
			NIP:           comp.NIP(),
			Address:       comp.Address(),
			Industry:      comp.Industry(),
			CompanyStatus: comp.Status().String(),
		}
	}

	c.JSON(http.StatusOK, PaginatedCompanyResponse{
		Status: true,
		Data:   items,
		Total:  res.Total,
		Page:   req.Page,
		Limit:  req.Limit,
	})
}
