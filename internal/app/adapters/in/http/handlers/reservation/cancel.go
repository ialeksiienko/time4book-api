package reservation

import (
	"net/http"
	"time4book/internal/app/adapters/in/http/handlers"
	"time4book/internal/app/core/usecases/reservation"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Cancel godoc
// @Summary      Cancel reservation
// @Description  Cancel an active reservation
// @Tags         reservations
// @Produce      json
// @Param        id path string true "Reservation ID"
// @Success      200  {object}  handlers.SuccessResponse
// @Failure      400  {object}  handlers.ErrorResponse
// @Failure      500  {object}  handlers.ErrorResponse
// @Security     BearerAuth
// @Router       /reservations/{id}/cancel [post]
func (h *Handler) Cancel(c *gin.Context) {
	initiatorIDStr := c.GetString("userID")
	initiatorID, _ := uuid.Parse(initiatorIDStr)

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, handlers.ErrorResponse{
			Status: false,
			Error:  "invalid reservation id",
		})
		return
	}

	companyIDStr := c.GetString("companyID")
	companyID, _ := uuid.Parse(companyIDStr)

	req := &reservationcommands.CancelRequest{
		InitiatorID:   initiatorID,
		CompanyID:     companyID,
		ReservationID: id,
	}

	_, err = h.commands.Cancel.Execute(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, handlers.ErrorResponse{
			Status: false,
			Error:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, handlers.SuccessResponse{
		Status:  true,
		Message: "reservation cancelled",
	})
}

