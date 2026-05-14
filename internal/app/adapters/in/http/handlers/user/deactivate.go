package user

import (
	"net/http"
	"time4book/internal/app/adapters/in/http/handlers"
	"time4book/internal/app/core/usecases/user"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Deactivate godoc
// @Summary      Deactivate user
// @Description  Deactivate a user account
// @Tags         users
// @Produce      json
// @Param        id path string true "User ID"
// @Success      200  {object}  handlers.SuccessResponse
// @Failure      400  {object}  handlers.ErrorResponse
// @Failure      500  {object}  handlers.ErrorResponse
// @Security     BearerAuth
// @Router       /users/{id} [delete]
func (h *Handler) Deactivate(c *gin.Context) {
	initiatorIDStr := c.GetString("userID")
	initiatorID, _ := uuid.Parse(initiatorIDStr)

	targetIDStr := c.Param("id")
	targetID, err := uuid.Parse(targetIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, handlers.ErrorResponse{
			Status: false,
			Error:  "invalid user id",
		})
		return
	}

	req := &usercommands.DeactivateRequest{
		InitiatorID: initiatorID,
		TargetID:    targetID,
	}

	_, err = h.commands.Deactivate.Execute(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, handlers.ErrorResponse{
			Status: false,
			Error:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, handlers.SuccessResponse{
		Status:  true,
		Message: "user deactivated",
	})
}
