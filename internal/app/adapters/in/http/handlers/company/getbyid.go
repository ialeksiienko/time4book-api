package company

import (
	"net/http"
	"time4book/internal/app/adapters/in/http/handlers"
	companycommands "time4book/internal/app/core/usecases/company"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GetByIDResponse struct {
	Status        bool      `json:"status"`
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	NIP           *string   `json:"nip,omitempty"`
	Address       *string   `json:"address,omitempty"`
	Industry      *string   `json:"industry,omitempty"`
	CompanyStatus string    `json:"companyStatus"`
}

// GetByID godoc
// @Summary      Get company by ID
// @Description  Get company details by ID
// @Tags         companies
// @Produce      json
// @Param        id path string true "Company ID"
// @Success      200  {object}  GetByIDResponse
// @Failure      400  {object}  handlers.ErrorResponse
// @Failure      500  {object}  handlers.ErrorResponse
// @Security     BearerAuth
// @Router       /companies/{id} [get]
func (h *Handler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, handlers.ErrorResponse{
			Status: false,
			Error:  "invalid company id",
		})
		return
	}

	req := &companycommands.GetByIDRequest{
		CompanyID: id,
	}

	res, err := h.commands.GetByID.Execute(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, handlers.ErrorResponse{
			Status: false,
			Error:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, GetByIDResponse{
		Status:        true,
		ID:            res.Company.ID(),
		Name:          res.Company.Name(),
		NIP:           res.Company.NIP(),
		Address:       res.Company.Address(),
		Industry:      res.Company.Industry(),
		CompanyStatus: res.Company.Status().String(),
	})
}
