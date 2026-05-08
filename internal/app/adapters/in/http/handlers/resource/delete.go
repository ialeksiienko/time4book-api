package resource

import (
	"net/http"
	"time4book/internal/app/adapters/in/http/handlers"
	"time4book/internal/app/core/usecases/resource"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Delete godoc
// @Summary      Delete resource
// @Description  Delete a resource
// @Tags         resources
// @Produce      json
// @Param        id path string true "Resource ID"
// @Success      200  {object}  handlers.SuccessResponse
// @Failure      400  {object}  handlers.ErrorResponse
// @Failure      500  {object}  handlers.ErrorResponse
// @Security     BearerAuth
// @Router       /resources/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	initiatorIDStr := c.GetString("userID")
	initiatorID, _ := uuid.Parse(initiatorIDStr)

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, handlers.ErrorResponse{
			Status: false,
			Error:  "invalid resource id",
		})
		return
	}

	req := &resourcecommands.DeleteRequest{
		InitiatorID: initiatorID,
		ResourceID:  id,
	}

	_, err = h.commands.Delete.Execute(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, handlers.ErrorResponse{
			Status: false,
			Error:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, handlers.SuccessResponse{
		Status:  true,
		Message: "resource deleted",
	})
}

