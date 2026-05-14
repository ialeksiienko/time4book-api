package company

import (
	"errors"
	"net/http"
	"time4book/internal/app/adapters/in/http/handlers"
	"time4book/internal/app/core/domain/model/user"
	companycommands "time4book/internal/app/core/usecases/company"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Delete godoc
// @Summary      Delete company
// @Description  Delete a company with all its resources, resource types and reservations
// @Tags         companies
// @Produce      json
// @Param        id path string true "Company ID"
// @Success      200  {object}  handlers.SuccessResponse
// @Failure      400  {object}  handlers.ErrorResponse
// @Failure      403  {object}  handlers.ErrorResponse
// @Failure      500  {object}  handlers.ErrorResponse
// @Security     BearerAuth
// @Router       /companies/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	initiatorID, err := uuid.Parse(c.GetString("userID"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, handlers.ErrorResponse{
			Status: false,
			Error:  "invalid user",
		})
		return
	}

	companyID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, handlers.ErrorResponse{
			Status: false,
			Error:  "invalid company id",
		})
		return
	}

	_, err = h.commands.Delete.Execute(c.Request.Context(), &companycommands.DeleteRequest{
		InitiatorID: initiatorID,
		CompanyID:   companyID,
	})
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
		Message: "company deleted",
	})
}
