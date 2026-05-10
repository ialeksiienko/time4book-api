package auth

import (
	"net/http"
	"time4book/internal/app/adapters/in/http/handlers"
	authcommands "time4book/internal/app/core/usecases/auth"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MeResponse struct {
	Status bool `json:"status"`
	User   struct {
		ID        uuid.UUID `json:"id"`
		Firstname string    `json:"firstName"`
		Lastname  string    `json:"lastName"`
		Email     string    `json:"email"`
		Role      string    `json:"role"`
		CompanyID uuid.UUID `json:"companyId,omitempty"`
	} `json:"user"`
	Company *struct {
		ID       uuid.UUID `json:"id"`
		Name     string    `json:"name"`
		NIP      *string   `json:"nip,omitempty"`
		Address  *string   `json:"address,omitempty"`
		Industry *string   `json:"industry,omitempty"`
	} `json:"company,omitempty"`
}

// Me godoc
// @Summary      Get current user
// @Description  Returns the profile of the currently authenticated user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      200  {object}  MeResponse
// @Failure      401  {object}  handlers.ErrorResponse
// @Failure      500  {object}  handlers.ErrorResponse
// @Security     BearerAuth
// @Router       /auth/me [get]
func (h *Handler) Me(c *gin.Context) {
	userIDStr := c.GetString("userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, handlers.ErrorResponse{
			Status: false,
			Error:  "invalid user id in token",
		})
		return
	}

	req := &authcommands.MeRequest{
		UserID: userID,
	}

	res, err := h.commands.Me.Execute(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, handlers.ErrorResponse{
			Status: false,
			Error:  err.Error(),
		})
		return
	}

	respData := MeResponse{Status: true}
	respData.User.ID = res.User.ID()
	respData.User.Firstname = res.User.Firstname()
	respData.User.Lastname = res.User.Lastname()
	respData.User.Email = res.User.Email()
	respData.User.Role = res.User.Role().String()
	respData.User.CompanyID = res.User.CompanyID()

	if res.Company != nil {
		respData.Company = &struct {
			ID       uuid.UUID `json:"id"`
			Name     string    `json:"name"`
			NIP      *string   `json:"nip,omitempty"`
			Address  *string   `json:"address,omitempty"`
			Industry *string   `json:"industry,omitempty"`
		}{
			ID:       res.Company.ID(),
			Name:     res.Company.Name(),
			NIP:      res.Company.NIP(),
			Address:  res.Company.Address(),
			Industry: res.Company.Industry(),
		}
	}

	c.JSON(http.StatusOK, respData)
}
