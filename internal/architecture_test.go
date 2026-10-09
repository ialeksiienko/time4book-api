package internal_test

import (
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestArchitectureDependencies(t *testing.T) {
	err := filepath.WalkDir("app", func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		file, err := parser.ParseFile(token.NewFileSet(), path, nil, parser.ImportsOnly)
		if err != nil {
			return err
		}
		path = filepath.ToSlash(path)
		for _, spec := range file.Imports {
			dependency, err := strconv.Unquote(spec.Path.Value)
			if err != nil {
				return err
			}
			local := strings.TrimPrefix(dependency, "time4book/internal/app/")
			forbidden := false
			switch {
			case strings.HasPrefix(path, "app/core/domain/"):
				forbidden = strings.HasPrefix(local, "adapters/") || strings.HasPrefix(local, "core/") && !strings.HasPrefix(local, "core/domain/")
			case strings.HasPrefix(path, "app/core/ports/"):
				forbidden = strings.HasPrefix(local, "adapters/") || strings.HasPrefix(local, "core/usecases/") || strings.HasPrefix(local, "core/queries/")
			case strings.HasPrefix(path, "app/core/"):
				forbidden = strings.HasPrefix(local, "adapters/")
			case strings.HasPrefix(path, "app/adapters/in/"):
				forbidden = strings.HasPrefix(local, "adapters/out/") || strings.HasPrefix(local, "core/usecases/") || strings.HasPrefix(local, "core/queries/") || local == "core"
			case strings.HasPrefix(path, "app/adapters/out/"):
				forbidden = strings.HasPrefix(local, "adapters/in/") || strings.HasPrefix(local, "core/usecases/") || strings.HasPrefix(local, "core/queries/") || local == "core"
			}
			if strings.HasPrefix(path, "app/core/") {
				for _, prefix := range []string{"github.com/gin-gonic/", "github.com/jackc/", "github.com/golang-jwt/"} {
					forbidden = forbidden || strings.HasPrefix(dependency, prefix)
				}
			}
			if forbidden {
				t.Errorf("%s imports %s across the architecture boundary", path, dependency)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
