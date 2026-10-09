package resource_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	handler "time4book/internal/app/adapters/in/http/handlers/resource"
	port "time4book/internal/app/core/ports/in/resource"
	"time4book/internal/app/core/queries"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type createCommand struct {
	request *port.CreateRequest
	id      uuid.UUID
}

func (c *createCommand) Execute(_ context.Context, request *port.CreateRequest) (*port.CreateResponse, error) {
	c.request = request
	return &port.CreateResponse{ResourceID: c.id}, nil
}

type getQuery struct {
	request  *port.GetRequest
	resource queries.Resource
}

func (q *getQuery) Execute(_ context.Context, request *port.GetRequest) (*port.GetResponse, error) {
	q.request = request
	return &port.GetResponse{Resources: []queries.Resource{q.resource}, Total: 1}, nil
}

func TestCreateUsesCommandPort(t *testing.T) {
	companyID, userID, resourceID := uuid.New(), uuid.New(), uuid.New()
	command := &createCommand{id: resourceID}
	h := handler.NewHandler(port.Commands{Create: command}, port.Queries{})
	router := gin.New()
	router.POST("/resources", func(c *gin.Context) {
		c.Set("companyID", companyID.String())
		c.Set("userID", userID.String())
		h.Create(c)
	})
	request := httptest.NewRequest(http.MethodPost, "/resources", strings.NewReader(`{"name":"Meeting room","type":"room"}`))
	request.Header.Set("Content-Type", "application/json")
	reply := httptest.NewRecorder()
	router.ServeHTTP(reply, request)
	if reply.Code != http.StatusCreated {
		t.Fatalf("status %d: %s", reply.Code, reply.Body.String())
	}
	if command.request == nil || command.request.CompanyID != companyID || command.request.InitiatorID != userID || command.request.Name != "Meeting room" {
		t.Fatalf("unexpected command request: %+v", command.request)
	}
	var body struct {
		ResourceID uuid.UUID `json:"resourceId"`
	}
	if err := json.Unmarshal(reply.Body.Bytes(), &body); err != nil || body.ResourceID != resourceID {
		t.Fatalf("unexpected response: %s (%v)", reply.Body.String(), err)
	}
}

func TestListUsesQueryPortAndPreservesResponse(t *testing.T) {
	companyID, resourceID, typeID := uuid.New(), uuid.New(), uuid.New()
	query := &getQuery{resource: queries.Resource{
		ID: resourceID, CompanyID: companyID, Name: "Room", Type: "custom", Status: "active",
		CompanyResourceType: &queries.CompanyResourceType{ID: typeID, Name: "Meeting room", IconKey: "room"},
	}}
	h := handler.NewHandler(port.Commands{}, port.Queries{Get: query})
	router := gin.New()
	router.GET("/resources", func(c *gin.Context) {
		c.Set("companyID", companyID.String())
		h.List(c)
	})
	reply := httptest.NewRecorder()
	router.ServeHTTP(reply, httptest.NewRequest(http.MethodGet, "/resources?page=2&limit=5&search=Room", nil))
	if reply.Code != http.StatusOK {
		t.Fatalf("status %d: %s", reply.Code, reply.Body.String())
	}
	if query.request == nil || query.request.CompanyID != companyID || query.request.Page != 2 || query.request.Limit != 5 || query.request.Search == nil || *query.request.Search != "Room" {
		t.Fatalf("unexpected query request: %+v", query.request)
	}
	var body struct {
		Data []struct {
			ID     uuid.UUID `json:"id"`
			Status string    `json:"resourceStatus"`
			Type   struct {
				ID   uuid.UUID `json:"id"`
				Icon string    `json:"iconKey"`
			} `json:"companyResourceType"`
		} `json:"data"`
		Total int64 `json:"total"`
	}
	if err := json.Unmarshal(reply.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if len(body.Data) != 1 || body.Total != 1 || body.Data[0].ID != resourceID || body.Data[0].Status != "active" || body.Data[0].Type.ID != typeID || body.Data[0].Type.Icon != "room" {
		t.Fatalf("unexpected response: %s", reply.Body.String())
	}
}
