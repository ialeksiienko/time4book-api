package company

import (
	"errors"
	"net/http"
	"time4book/internal/app/adapters/in/http/handlers"
	"time4book/internal/app/core/domain/model/user"
	"time4book/internal/app/core/usecases/company"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Block godoc
// @Summary      Block company
// @Description  Block a company account
// @Tags         companies
// @Produce      json
// @Param        id path string true "Company ID"
// @Success      200  {object}  handlers.SuccessResponse
// @Failure      400  {object}  handlers.ErrorResponse
// @Failure      500  {object}  handlers.ErrorResponse
// @Security     BearerAuth
// @Router       /companies/{id}/block [post]
func (h *Handler) Block(c *gin.Context) {
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

	req := &companycommands.BlockRequest{
		InitiatorID: initiatorID,
		CompanyID:   id,
	}

	_, err = h.commands.Block.Execute(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, user.ErrUnauthorized) {
			c.JSON(http.StatusForbidden, handlers.ErrorResponse{
				Status: false,
				Error:  err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, handlers.ErrorResponse{
			Status: false,
			Error:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, handlers.SuccessResponse{
		Status:  true,
		Message: "company blocked",
	})
}
