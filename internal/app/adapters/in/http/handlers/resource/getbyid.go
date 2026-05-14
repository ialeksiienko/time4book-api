package resource

import (
	"net/http"
	"time4book/internal/app/adapters/in/http/handlers"
	"time4book/internal/app/core/usecases/resource"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetByIDResponse keeps all resource fields at the top level so existing JSON clients remain compatible.
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

	req := &resourcecommands.GetByIDRequest{
		ResourceID: id,
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
		Status:       true,
		ResourceBody: toResourceBody(r),
	})
}
