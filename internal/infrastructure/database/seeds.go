package database

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
)

type seedDB interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

const defaultSeedsDir = "database/seeds"

// SeedsDir returns the directory of idempotent SQL seed files.
func SeedsDir() string {
	if dir := os.Getenv("DATABASE_SEEDS_DIR"); dir != "" {
		return dir
	}
	return defaultSeedsDir
}

// ApplySeeds runs all *.sql seed files in dir in lexicographic order.
// Seed scripts must be idempotent (ON CONFLICT DO NOTHING, etc.).
func ApplySeeds(ctx context.Context, db seedDB, dir string) error {
	files, err := filepath.Glob(filepath.Join(dir, "*.sql"))
	if err != nil {
		return fmt.Errorf("list seed files: %w", err)
	}
	if len(files) == 0 {
		return fmt.Errorf("no seed files in %s", dir)
	}
	sort.Strings(files)

	for _, file := range files {
		name := filepath.Base(file)
		sqlBytes, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("read seed %s: %w", name, err)
		}
		slog.InfoContext(ctx, "applying database seed", "file", name)
		if _, err := db.ExecContext(ctx, string(sqlBytes)); err != nil {
			return fmt.Errorf("apply seed %s: %w", name, err)
		}
	}
	return nil
}
