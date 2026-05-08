package company

import (
	"net/http"
	"time4book/internal/app/adapters/in/http/handlers"
	"time4book/internal/app/core/usecases/company"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UpdateRequest struct {
	Name     string  `json:"name" binding:"required"`
	NIP      *string `json:"nip,omitempty"`
	Address  *string `json:"address,omitempty"`
	Industry *string `json:"industry,omitempty"`
}

// Update godoc
// @Summary      Update company
// @Description  Update company details
// @Tags         companies
// @Accept       json
// @Produce      json
// @Param        id path string true "Company ID"
// @Param        request body UpdateRequest true "Update parameters"
// @Success      200  {object}  handlers.SuccessResponse
// @Failure      400  {object}  handlers.ErrorResponse
// @Failure      500  {object}  handlers.ErrorResponse
// @Security     BearerAuth
// @Router       /companies/{id} [put]
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

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, handlers.ErrorResponse{
			Status: false,
			Error:  "invalid company id",
		})
		return
	}

	req := &companycommands.UpdateRequest{
		InitiatorID: initiatorID,
		CompanyID:   id,
		Name:        body.Name,
		NIP:         body.NIP,
		Address:     body.Address,
		Industry:    body.Industry,
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
		Message: "company updated",
	})
}

