package postgres

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/weprodev/go-pkg/pgsql"
)

func newMockPgClient(t *testing.T) (*pgsql.PgClient, sqlmock.Sqlmock, *sql.DB) {
	t.Helper()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	return &pgsql.PgClient{DB: db}, mock, db
}
