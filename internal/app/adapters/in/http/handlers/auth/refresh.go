package auth

import (
	"errors"
	"net/http"
	"time4book/internal/app/adapters/in/http/handlers"
	authcommands "time4book/internal/app/core/usecases/auth"

	"github.com/gin-gonic/gin"
)

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type RefreshResponse struct {
	Status       bool   `json:"status"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// Refresh godoc
// @Summary      Refresh token
// @Description  Get a new access token using a refresh token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body RefreshRequest true "Refresh token"
// @Success      200  {object}  RefreshResponse
// @Failure      400  {object}  handlers.ErrorResponse
// @Failure      401  {object}  handlers.ErrorResponse
// @Router       /auth/refresh [post]
func (h *Handler) Refresh(c *gin.Context) {
	var body RefreshRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, handlers.ErrorResponse{
			Status: false,
			Error:  err.Error(),
		})
		return
	}

	req := &authcommands.RefreshRequest{
		RefreshToken: body.RefreshToken,
	}

	res, err := h.commands.Refresh.Execute(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, authcommands.ErrCompanyBlocked) {
			c.JSON(http.StatusForbidden, handlers.ErrorResponse{
				Status: false,
				Error:  err.Error(),
			})
			return
		}
		c.JSON(http.StatusUnauthorized, handlers.ErrorResponse{
			Status: false,
			Error:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, RefreshResponse{
		Status:       true,
		AccessToken:  res.AccessToken.Value,
		RefreshToken: res.RefreshToken.Value,
	})
}
