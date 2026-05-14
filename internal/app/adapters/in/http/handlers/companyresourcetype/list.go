package companyresourcetype

import (
	"net/http"
	"time4book/internal/app/adapters/in/http/handlers"
	companyresourcetypecommands "time4book/internal/app/core/usecases/companyresourcetype"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ListItemBody struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	IconKey string    `json:"iconKey"`
}

func (h *Handler) List(c *gin.Context) {
	companyIDStr, ok := c.Get("companyID")
	if !ok {
		c.JSON(http.StatusBadRequest, handlers.ErrorResponse{Status: false, Error: "missing company scope; pass companyId (developer)"})
		return
	}
	companyID, ok := companyIDStr.(uuid.UUID)
	if !ok || companyID == uuid.Nil {
		c.JSON(http.StatusBadRequest, handlers.ErrorResponse{Status: false, Error: "invalid company scope"})
		return
	}

	initiatorID, err := uuid.Parse(c.GetString("userID"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, handlers.ErrorResponse{Status: false, Error: "invalid user"})
		return
	}

	res, err := h.commands.List.Execute(c.Request.Context(), &companyresourcetypecommands.ListRequest{
		InitiatorID: initiatorID,
		CompanyID:   companyID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, handlers.ErrorResponse{Status: false, Error: err.Error()})
		return
	}

	out := make([]ListItemBody, 0, len(res.Items))
	for _, t := range res.Items {
		out = append(out, ListItemBody{ID: t.ID(), Name: t.Name(), IconKey: t.IconKey()})
	}

	c.JSON(http.StatusOK, gin.H{"status": true, "data": out})
}
