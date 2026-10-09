package companyresourcetype

import (
	"net/http"
	"time4book/internal/app/adapters/http/handlers"
	companyresourcetypecommands "time4book/internal/app/core/usecases/companyresourcetype"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ListItem struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	IconKey string    `json:"iconKey"`
}
type ListResponse struct {
	Status bool       `json:"status"`
	Data   []ListItem `json:"data"`
}

func (h *Handler) List(c *gin.Context) {
	initiatorID, err := userID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, handlers.ErrorResponse{Status: false, Error: "invalid user"})
		return
	}
	scope, err := companyID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, handlers.ErrorResponse{Status: false, Error: "invalid company scope"})
		return
	}
	if scope == uuid.Nil {
		c.JSON(http.StatusBadRequest, handlers.ErrorResponse{Status: false, Error: "company scope required"})
		return
	}
	result, err := h.commands.List.Execute(c.Request.Context(), &companyresourcetypecommands.ListRequest{InitiatorID: initiatorID, CompanyID: scope})
	if err != nil {
		writeError(c, err)
		return
	}
	items := make([]ListItem, len(result.Items))
	for i, item := range result.Items {
		items[i] = ListItem{ID: item.ID(), Name: item.Name(), IconKey: item.IconKey()}
	}
	c.JSON(http.StatusOK, ListResponse{Status: true, Data: items})
}
