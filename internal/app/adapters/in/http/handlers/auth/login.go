package auth

import (
	"errors"
	"net/http"
	"time4book/internal/app/adapters/in/http/handlers"
	authcommands "time4book/internal/app/core/usecases/auth"
	"time4book/pkg/validator"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Status       bool      `json:"status"`
	UserID       uuid.UUID `json:"userId"`
	AccessToken  string    `json:"accessToken"`
	RefreshToken string    `json:"refreshToken"`
}

// Login godoc
// @Summary      Login user
// @Description  Authenticates a user and returns access and refresh tokens
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body LoginRequest true "Login credentials"
// @Success      200  {object}  LoginResponse
// @Failure      400  {object}  handlers.ErrorResponse
// @Failure      401  {object}  handlers.ErrorResponse
// @Failure      500  {object}  handlers.ErrorResponse
// @Router       /auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var body LoginRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, handlers.ErrorResponse{
			Status: false,
			Error:  err.Error(),
		})
		return
	}

	req := &authcommands.LoginRequest{
		Email:    body.Email,
		Password: body.Password,
	}

	res, err := h.commands.Login.Execute(c.Request.Context(), req)
	if err != nil {
		var validationErr validator.ValidationErrors
		if errors.As(err, &validationErr) {
			c.JSON(http.StatusBadRequest, handlers.ErrorResponse{
				Status: false,
				Error:  validationErr.Error(),
			})
			return
		}

		if errors.Is(err, authcommands.ErrCompanyBlocked) {
			c.JSON(http.StatusForbidden, handlers.ErrorResponse{
				Status: false,
				Error:  err.Error(),
			})
			return
		}

		c.JSON(http.StatusUnauthorized, handlers.ErrorResponse{
			Status: false,
			Error:  "invalid email or password",
		})
		return
	}

	c.SetCookie("refresh_token", res.RefreshToken.Value, int(res.RefreshToken.ExpiresAt.Unix()), "/", "", false, true)

	c.JSON(http.StatusOK, LoginResponse{
		Status:       true,
		UserID:       res.UserID,
		AccessToken:  res.AccessToken.Value,
		RefreshToken: res.RefreshToken.Value,
	})
}
