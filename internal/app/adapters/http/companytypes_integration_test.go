package httpadapter

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
	"time4book/internal/app/adapters/postgres"
	authrepo "time4book/internal/app/adapters/postgres/repo/auth"
	companyrepo "time4book/internal/app/adapters/postgres/repo/company"
	companytyperepo "time4book/internal/app/adapters/postgres/repo/companyresourcetype"
	reservationrepo "time4book/internal/app/adapters/postgres/repo/reservation"
	resourcerepo "time4book/internal/app/adapters/postgres/repo/resource"
	userrepo "time4book/internal/app/adapters/postgres/repo/user"
	"time4book/internal/app/core/jwt"
	"time4book/internal/app/core/usecases"
	"time4book/pkg/validator"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

func TestCompanyResourceTypeFlow(t *testing.T) {
	databaseURL := os.Getenv("TIME4BOOK_TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("set TIME4BOOK_TEST_DATABASE_URL for PostgreSQL integration test")
	}
	ctx := context.Background()
	pool, err := pgxpool.Connect(ctx, databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	companyA, companyB := uuid.New(), uuid.New()
	ownerA, ownerB, developer := uuid.New(), uuid.New(), uuid.New()
	for _, company := range []struct{ id, owner uuid.UUID }{{companyA, ownerA}, {companyB, ownerB}} {
		if _, err := pool.Exec(ctx, "INSERT INTO companies (id,owner_id,name) VALUES ($1,$2,$3)", company.id, company.owner, company.id.String()); err != nil {
			t.Fatal(err)
		}
	}
	for _, account := range []struct {
		id, company uuid.UUID
		role        string
	}{{ownerA, companyA, "owner"}, {ownerB, companyB, "owner"}, {developer, companyA, "developer"}} {
		if _, err := pool.Exec(ctx, "INSERT INTO users (id,company_id,firstname,lastname,email,role_id) VALUES ($1,$2,'Test','User',$3,$4)", account.id, account.company, account.id.String()+"@example.com", account.role); err != nil {
			t.Fatal(err)
		}
	}

	db := postgres.New(pool)
	users := userrepo.New(db)
	types := companytyperepo.New(db)
	manager := jwt.NewManager("integration-secret", time.Hour, 24*time.Hour)
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	commands := usecases.New(users, authrepo.New(db), companyrepo.New(db), types, resourcerepo.New(db), reservationrepo.New(db), postgres.NewTxManager(db), manager, validator.New(), log)
	router := NewRouter(NewHandler(commands), JWTAuth(manager), RequireCompany(), ResolveCompany(users))

	request := func(id uuid.UUID, role, method, path string, body any, want int) map[string]any {
		t.Helper()
		var reader io.Reader
		if body != nil {
			payload, err := json.Marshal(body)
			if err != nil {
				t.Fatal(err)
			}
			reader = bytes.NewReader(payload)
		}
		req := httptest.NewRequest(method, path, reader)
		req.Header.Set("Content-Type", "application/json")
		tokens, err := manager.GenerateTokens(id, role)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Authorization", "Bearer "+tokens.AccessToken.Value)
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)
		if recorder.Code != want {
			t.Fatalf("%s %s: got %d, want %d: %s", method, path, recorder.Code, want, recorder.Body.String())
		}
		var result map[string]any
		if err := json.Unmarshal(recorder.Body.Bytes(), &result); err != nil {
			t.Fatal(err)
		}
		return result
	}

	typeBody := map[string]any{"companyId": companyA, "name": "Meeting room", "iconKey": "room"}
	created := request(ownerA, "owner", http.MethodPost, "/api/v1/company-resource-types", typeBody, http.StatusCreated)
	typeID, err := uuid.Parse(created["id"].(string))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, "INSERT INTO resources (id,company_id,name,type,company_resource_type_id) VALUES ($1,$2,'Invalid','custom',$3)", uuid.New(), companyB, typeID); err == nil {
		t.Fatal("database accepted a resource type from another company")
	}
	listed := request(ownerA, "owner", http.MethodGet, "/api/v1/company-resource-types", nil, http.StatusOK)
	items := listed["data"].([]any)
	if len(items) != 1 || items[0].(map[string]any)["id"] != typeID.String() {
		t.Fatalf("unexpected types: %v", listed)
	}
	request(ownerA, "owner", http.MethodPost, "/api/v1/company-resource-types", map[string]any{"companyId": companyA, "name": "MEETING ROOM", "iconKey": "room"}, http.StatusConflict)
	request(ownerA, "owner", http.MethodPost, "/api/v1/company-resource-types", map[string]any{"companyId": companyA, "name": "Invalid", "iconKey": "unknown"}, http.StatusBadRequest)
	request(ownerB, "owner", http.MethodGet, "/api/v1/company-resource-types?companyId="+companyA.String(), nil, http.StatusForbidden)
	request(ownerB, "owner", http.MethodPut, "/api/v1/company-resource-types/"+typeID.String(), map[string]any{"name": "Other room", "iconKey": "room"}, http.StatusForbidden)
	request(ownerB, "owner", http.MethodDelete, "/api/v1/company-resource-types/"+typeID.String(), nil, http.StatusForbidden)
	request(developer, "developer", http.MethodGet, "/api/v1/company-resource-types?companyId="+companyA.String(), nil, http.StatusOK)
	request(developer, "developer", http.MethodPut, "/api/v1/company-resource-types/"+typeID.String(), map[string]any{"name": "Conference room", "iconKey": "building"}, http.StatusOK)

	resourceBody := map[string]any{"companyId": companyA, "name": "Room A", "type": "custom", "companyResourceTypeId": typeID}
	request(ownerA, "owner", http.MethodPost, "/api/v1/resources", map[string]any{"companyId": companyA, "name": "Invalid", "type": "custom"}, http.StatusBadRequest)
	request(ownerA, "owner", http.MethodPost, "/api/v1/resources", map[string]any{"companyId": companyA, "name": "Invalid", "type": "unknown"}, http.StatusBadRequest)
	createdResource := request(ownerA, "owner", http.MethodPost, "/api/v1/resources", resourceBody, http.StatusCreated)
	resourceID := createdResource["resourceId"].(string)
	resourceResult := request(ownerA, "owner", http.MethodGet, "/api/v1/resources/"+resourceID, nil, http.StatusOK)
	reference := resourceResult["companyResourceType"].(map[string]any)
	if reference["id"] != typeID.String() || reference["name"] != "Conference room" {
		t.Fatalf("unexpected resource type: %v", reference)
	}
	listedResources := request(ownerA, "owner", http.MethodGet, "/api/v1/resources", nil, http.StatusOK)
	first := listedResources["data"].([]any)[0].(map[string]any)
	if first["companyId"] != companyA.String() || first["companyResourceType"].(map[string]any)["id"] != typeID.String() {
		t.Fatalf("unexpected resources: %v", listedResources)
	}
	request(ownerB, "owner", http.MethodPost, "/api/v1/resources", map[string]any{"companyId": companyB, "name": "Other room", "type": "custom", "companyResourceTypeId": typeID}, http.StatusForbidden)
	request(developer, "developer", http.MethodDelete, "/api/v1/company-resource-types/"+typeID.String(), nil, http.StatusConflict)
	request(ownerA, "owner", http.MethodPut, "/api/v1/resources/"+resourceID, map[string]any{"name": "Room A", "type": "room"}, http.StatusOK)
	request(developer, "developer", http.MethodDelete, "/api/v1/company-resource-types/"+typeID.String(), nil, http.StatusOK)
	legacy := request(ownerA, "owner", http.MethodPost, "/api/v1/resources", map[string]any{"companyId": companyA, "name": "Legacy room", "type": "room"}, http.StatusCreated)
	legacyResult := request(ownerA, "owner", http.MethodGet, "/api/v1/resources/"+legacy["resourceId"].(string), nil, http.StatusOK)
	if legacyResult["type"] != "room" || legacyResult["companyResourceType"] != nil {
		t.Fatalf("legacy resource changed: %v", legacyResult)
	}
}
