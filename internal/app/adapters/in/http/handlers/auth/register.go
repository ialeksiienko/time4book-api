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

type RegisterRequest struct {
	Firstname       string  `json:"firstName" binding:"required"`
	Lastname        string  `json:"lastName" binding:"required"`
	Email           string  `json:"email" binding:"required,email"`
	Password        string  `json:"password" binding:"required,min=8"`
	CompanyName     string  `json:"companyName" binding:"required"`
	CompanyNIP      *string `json:"companyNip,omitempty"`
	CompanyAddress  *string `json:"companyAddress,omitempty"`
	CompanyIndustry *string `json:"companyIndustry,omitempty"`
}

type RegisterResponse struct {
	Status       bool      `json:"status"`
	UserID       uuid.UUID `json:"userId"`
	CompanyID    uuid.UUID `json:"companyId"`
	AccessToken  string    `json:"accessToken"`
	RefreshToken string    `json:"refreshToken"`
}

// Register godoc
// @Summary      Register a new user
// @Description  Registers a user and creates an associated company
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body RegisterRequest true "Registration details"
// @Success      201  {object}  RegisterResponse
// @Failure      400  {object}  handlers.ErrorResponse
// @Failure      500  {object}  handlers.ErrorResponse
// @Router       /auth/register [post]
func (h *Handler) Register(c *gin.Context) {
	var body RegisterRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, handlers.ErrorResponse{
			Status: false,
			Error:  err.Error(),
		})
		return
	}

	req := &authcommands.RegisterRequest{
		Firstname:       body.Firstname,
		Lastname:        body.Lastname,
		Email:           body.Email,
		Password:        body.Password,
		CompanyName:     body.CompanyName,
		CompanyNIP:      body.CompanyNIP,
		CompanyAddress:  body.CompanyAddress,
		CompanyIndustry: body.CompanyIndustry,
	}

	res, err := h.commands.Register.Execute(c.Request.Context(), req)
	if err != nil {
		var validationErr validator.ValidationErrors
		if errors.As(err, &validationErr) {
			c.JSON(http.StatusBadRequest, handlers.ErrorResponse{
				Status: false,
				Error:  validationErr.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, handlers.ErrorResponse{
			Status: false,
			Error:  err.Error(),
		})
		return
	}

	c.SetCookie("refresh_token", res.RefreshToken.Value, int(res.RefreshToken.ExpiresAt.Unix()), "/", "", false, true)

	c.JSON(http.StatusCreated, RegisterResponse{
		Status:       true,
		UserID:       res.UserID,
		CompanyID:    res.CompanyID,
		AccessToken:  res.AccessToken.Value,
		RefreshToken: res.RefreshToken.Value,
	})
}
