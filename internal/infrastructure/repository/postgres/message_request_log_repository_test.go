package postgres

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestScanMessageRequestLogWithSource_nullInboxMessageID(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	createdAt := time.Date(2026, 7, 4, 9, 0, 0, 0, time.UTC)
	rows := sqlmock.NewRows([]string{
		"id", "workspace_id", "api_key_id", "channel_type", "http_method", "status_code", "endpoint",
		"provider_name", "request_id", "duration_ms", "error_message", "inbox_message_id", "created_at",
		"source_name", "client_id", "total_count",
	}).AddRow(
		"log-1", "ws-1", nil, "email", "POST", 200, "/v1/email",
		"memory", nil, 12, nil, nil, createdAt,
		"Demo Key", "demo-client", int64(1),
	)

	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	resultRows, err := db.Query("SELECT")
	if err != nil {
		t.Fatalf("Query: %v", err)
	}
	t.Cleanup(func() { _ = resultRows.Close() })

	if !resultRows.Next() {
		t.Fatal("expected one row")
	}

	row, total, err := scanMessageRequestLogWithSource(resultRows)
	if err != nil {
		t.Fatalf("scanMessageRequestLogWithSource: %v", err)
	}

	if row.InboxMessageID != "" {
		t.Fatalf("InboxMessageID: got %q, want empty for SQL NULL", row.InboxMessageID)
	}
	if row.ProviderName != "memory" {
		t.Fatalf("ProviderName: got %q", row.ProviderName)
	}
	if total != 1 {
		t.Fatalf("total: got %d", total)
	}
}

func TestScanMessageRequestLogWithSource_withInboxMessageID(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	createdAt := time.Date(2026, 7, 4, 9, 0, 0, 0, time.UTC)
	inboxID := "inbox-msg-99"
	rows := sqlmock.NewRows([]string{
		"id", "workspace_id", "api_key_id", "channel_type", "http_method", "status_code", "endpoint",
		"provider_name", "request_id", "duration_ms", "error_message", "inbox_message_id", "created_at",
		"source_name", "client_id", "total_count",
	}).AddRow(
		"log-2", "ws-1", nil, "email", "POST", 200, "/v1/email",
		"memory", nil, 12, nil, inboxID, createdAt,
		"Demo Key", "demo-client", int64(1),
	)

	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	resultRows, err := db.Query("SELECT")
	if err != nil {
		t.Fatalf("Query: %v", err)
	}
	t.Cleanup(func() { _ = resultRows.Close() })

	if !resultRows.Next() {
		t.Fatal("expected one row")
	}

	row, _, err := scanMessageRequestLogWithSource(resultRows)
	if err != nil {
		t.Fatalf("scanMessageRequestLogWithSource: %v", err)
	}

	if row.InboxMessageID != inboxID {
		t.Fatalf("InboxMessageID: got %q, want %q", row.InboxMessageID, inboxID)
	}
}

func TestOptionalStringFromNull(t *testing.T) {
	t.Parallel()

	if got := optionalStringFromNull(sql.NullString{}); got != "" {
		t.Fatalf("NULL: got %q", got)
	}
	if got := optionalStringFromNull(sql.NullString{String: "x", Valid: true}); got != "x" {
		t.Fatalf("valid: got %q", got)
	}
}
