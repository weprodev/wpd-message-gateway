package database

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestSeedsDir_default(t *testing.T) {
	t.Setenv("DATABASE_SEEDS_DIR", "")
	if got := SeedsDir(); got != defaultSeedsDir {
		t.Fatalf("SeedsDir() = %q, want %q", got, defaultSeedsDir)
	}
}

func TestSeedsDir_override(t *testing.T) {
	t.Setenv("DATABASE_SEEDS_DIR", "/custom/seeds")
	if got := SeedsDir(); got != "/custom/seeds" {
		t.Fatalf("SeedsDir() = %q", got)
	}
}

func TestApplySeeds_runsInOrder(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"002_second.sql", "001_first.sql"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("SELECT 1;"), 0o600); err != nil {
			t.Fatalf("write seed: %v", err)
		}
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock: %v", err)
	}
	defer db.Close()

	mock.ExpectExec("SELECT 1").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("SELECT 1").WillReturnResult(sqlmock.NewResult(0, 0))

	if err := ApplySeeds(context.Background(), db, dir); err != nil {
		t.Fatalf("ApplySeeds: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestApplySeeds_noFiles(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock: %v", err)
	}
	defer db.Close()

	err = ApplySeeds(context.Background(), db, t.TempDir())
	if err == nil {
		t.Fatal("expected error for empty seeds dir")
	}
}
