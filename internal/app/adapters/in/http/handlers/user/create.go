package user

import (
	"errors"
	"net/http"
	"time4book/internal/app/adapters/in/http/handlers"
	usercommands "time4book/internal/app/core/usecases/user"
	"time4book/pkg/validator"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateRequest struct {
	Firstname string `json:"firstName" binding:"required"`
	Lastname  string `json:"lastName" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	Role      string `json:"role" binding:"required"`
}

type CreateResponse struct {
	Status bool      `json:"status"`
	UserID uuid.UUID `json:"userId"`
}

// Create godoc
// @Summary      Create user
// @Description  Create a new user within a company
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request body CreateRequest true "User details"
// @Success      201  {object}  CreateResponse
// @Failure      400  {object}  handlers.ErrorResponse
// @Failure      500  {object}  handlers.ErrorResponse
// @Security     BearerAuth
// @Router       /users [post]
func (h *Handler) Create(c *gin.Context) {
	var body CreateRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, handlers.ErrorResponse{
			Status: false,
			Error:  err.Error(),
		})
		return
	}

	initiatorIDStr := c.GetString("userID")
	initiatorID, _ := uuid.Parse(initiatorIDStr)

	companyIDStr := c.GetString("companyID")
	companyID, _ := uuid.Parse(companyIDStr)

	req := &usercommands.CreateRequest{
		InitiatorID: initiatorID,
		CompanyID:   companyID,
		Firstname:   body.Firstname,
		Lastname:    body.Lastname,
		Email:       body.Email,
		Password:    body.Password,
		Role:        body.Role,
	}

	res, err := h.commands.Create.Execute(c.Request.Context(), req)
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

	c.JSON(http.StatusCreated, CreateResponse{
		Status: true,
		UserID: res.UserID,
	})
}
