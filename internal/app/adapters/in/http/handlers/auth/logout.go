package auth

import (
	"net/http"
	"time4book/internal/app/adapters/in/http/handlers"
	authcommands "time4book/internal/app/core/usecases/auth"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Logout godoc
// @Summary      Logout user
// @Description  Logs out the current user by invalidating the refresh token
// @Tags         auth
// @Produce      json
// @Success      200  {object}  handlers.SuccessResponse
// @Failure      401  {object}  handlers.ErrorResponse
// @Failure      500  {object}  handlers.ErrorResponse
// @Security     BearerAuth
// @Router       /auth/logout [post]
func (h *Handler) Logout(c *gin.Context) {
	userIDStr := c.GetString("userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, handlers.ErrorResponse{
			Status: false,
			Error:  "invalid user id in token",
		})
		return
	}

	req := &authcommands.LogoutRequest{
		UserID: userID,
	}

	_, err = h.commands.Logout.Execute(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, handlers.ErrorResponse{
			Status: false,
			Error:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, handlers.SuccessResponse{
		Status:  true,
		Message: "logged out successfully",
	})
}
