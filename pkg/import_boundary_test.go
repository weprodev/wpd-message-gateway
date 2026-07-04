package pkg_test

import (
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strings"
	"testing"
)

func TestPkgDoesNotImportInternal(t *testing.T) {
	t.Parallel()

	const forbidden = "github.com/weprodev/wpd-message-gateway/internal/"
	fset := token.NewFileSet()

	err := filepath.WalkDir(".", func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		file, err := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
		if err != nil {
			return err
		}

		for _, imp := range file.Imports {
			importPath := strings.Trim(imp.Path.Value, `"`)
			if strings.HasPrefix(importPath, forbidden) {
				t.Errorf("%s imports forbidden package %s", path, importPath)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk pkg: %v", err)
	}
}
