package resource

import (
	"net/http"
	"time4book/internal/app/adapters/in/http/handlers"
	port "time4book/internal/app/core/ports/in/resource"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GetByIDResponse struct {
	Status bool `json:"status"`
	ResourceBody
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

	req := &port.RetrieveRequest{
		ResourceID: id,
	}

	if companyIDStr := c.GetString("companyID"); companyIDStr != "" {
		if cid, err := uuid.Parse(companyIDStr); err == nil {
			req.CompanyID = cid
		}
	}

	res, err := h.queries.Retrieve.Execute(c.Request.Context(), req)
	if err != nil {
		writeError(c, err)
		return
	}

	r := res.Resource
	c.JSON(http.StatusOK, GetByIDResponse{Status: true, ResourceBody: toResourceBody(r)})
}
