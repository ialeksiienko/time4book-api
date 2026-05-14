package user

import (
	"net/http"
	"strconv"
	"time"
	"time4book/internal/app/adapters/in/http/handlers"
	usermodel "time4book/internal/app/core/domain/model/user"
	usercommands "time4book/internal/app/core/usecases/user"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ListResponse struct {
	ID        uuid.UUID `json:"id"`
	Firstname string    `json:"firstName"`
	Lastname  string    `json:"lastName"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}

type PaginatedUserResponse struct {
	Status bool           `json:"status"`
	Data   []ListResponse `json:"data"`
	Total  int64          `json:"total"`
	Page   int            `json:"page"`
	Limit  int            `json:"limit"`
}

// List godoc
// @Summary      List users
// @Description  Get a paginated list of users
// @Tags         users
// @Produce      json
// @Param        page query int false "Page number"
// @Param        limit query int false "Items per page"
// @Param        companyId query string false "Company ID"
// @Param        search query string false "Search term"
// @Param        role query string false "Role"
// @Param        status query string false "Status"
// @Success      200  {object}  PaginatedUserResponse
// @Failure      500  {object}  handlers.ErrorResponse
// @Security     BearerAuth
// @Router       /users [get]
func (h *Handler) List(c *gin.Context) {
	req := &usercommands.ListRequest{
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
	if role == string(usermodel.RoleDeveloperKey) {
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

	if search := c.Query("search"); search != "" {
		req.Search = &search
	}
	if role := c.Query("role"); role != "" {
		req.Role = &role
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

	items := make([]ListResponse, len(res.Users))
	for i, u := range res.Users {
		items[i] = ListResponse{
			ID:        u.ID(),
			Firstname: u.Firstname(),
			Lastname:  u.Lastname(),
			Email:     u.Email(),
			Role:      u.Role().String(),
			Status:    u.Status().String(),
			CreatedAt: u.CreatedAt(),
		}
	}

	c.JSON(http.StatusOK, PaginatedUserResponse{
		Status: true,
		Data:   items,
		Total:  res.Total,
		Page:   req.Page,
		Limit:  req.Limit,
	})
}
