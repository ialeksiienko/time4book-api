package docs

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"

	"github.com/swaggo/swag"
)

func TestRegisteredSwaggerMatchesJSON(t *testing.T) {
	registered, err := swag.ReadDoc()
	if err != nil {
		t.Fatal(err)
	}
	file, err := os.ReadFile("swagger.json")
	if err != nil {
		t.Fatal(err)
	}
	var active, saved map[string]any
	if err := json.Unmarshal([]byte(registered), &active); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(file, &saved); err != nil {
		t.Fatal(err)
	}
	delete(active, "schemes")
	if !reflect.DeepEqual(active, saved) {
		t.Fatal("registered Swagger differs from swagger.json")
	}
}
