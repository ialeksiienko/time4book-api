package user

import (
	"net/http"
	"time4book/internal/app/adapters/in/http/handlers"
	"time4book/internal/app/core/usecases/user"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UpdateRequest struct {
	Firstname *string `json:"firstName"`
	Lastname  *string `json:"lastName"`
	Role      *string `json:"role"`
	Status    *string `json:"status"`
}

// Update godoc
// @Summary      Update user
// @Description  Update user profile, role, or status
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id path string true "User ID"
// @Param        request body UpdateRequest true "Update parameters"
// @Success      200  {object}  handlers.SuccessResponse
// @Failure      400  {object}  handlers.ErrorResponse
// @Failure      500  {object}  handlers.ErrorResponse
// @Security     BearerAuth
// @Router       /users/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	var body UpdateRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, handlers.ErrorResponse{
			Status: false,
			Error:  err.Error(),
		})
		return
	}

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

	req := &usercommands.UpdateRequest{
		InitiatorID: initiatorID,
		TargetID:    targetID,
		Firstname:   body.Firstname,
		Lastname:    body.Lastname,
		Role:        body.Role,
		Status:      body.Status,
	}

	_, err = h.commands.Update.Execute(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, handlers.ErrorResponse{
			Status: false,
			Error:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, handlers.SuccessResponse{
		Status:  true,
		Message: "user updated",
	})
}
