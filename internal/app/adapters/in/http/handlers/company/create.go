package company

import (
	"errors"
	"net/http"
	"time4book/internal/app/adapters/in/http/handlers"
	companycommands "time4book/internal/app/core/usecases/company"
	"time4book/pkg/validator"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateRequest struct {
	Name     string  `json:"name" binding:"required"`
	NIP      *string `json:"nip,omitempty"`
	Address  *string `json:"address,omitempty"`
	Industry *string `json:"industry,omitempty"`
}

type CreateResponse struct {
	Status    bool      `json:"status"`
	CompanyID uuid.UUID `json:"companyId"`
}

// Create godoc
// @Summary      Create company
// @Description  Create a new company
// @Tags         companies
// @Accept       json
// @Produce      json
// @Param        request body CreateRequest true "Company details"
// @Success      201  {object}  CreateResponse
// @Failure      400  {object}  handlers.ErrorResponse
// @Failure      500  {object}  handlers.ErrorResponse
// @Security     BearerAuth
// @Router       /companies [post]
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

	req := &companycommands.CreateRequest{
		OwnerID:  initiatorID,
		Name:     body.Name,
		NIP:      body.NIP,
		Address:  body.Address,
		Industry: body.Industry,
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
		Status:    true,
		CompanyID: res.CompanyID,
	})
}

