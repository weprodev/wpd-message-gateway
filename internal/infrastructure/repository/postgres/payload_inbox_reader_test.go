package postgres

import (
	"encoding/json"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/weprodev/wpd-message-gateway/pkg/contracts"
)

func TestPayloadInboxReader_ListEmailsForWorkspace(t *testing.T) {
	t.Parallel()

	client, mock, _ := newMockPgClient(t)
	repo := NewPayloadInboxReader(client, nil)
	workspaceID := "ws-1"
	email := contracts.Email{Subject: "Hello", To: []string{"a@test.com"}, HTML: "<p>Hi</p>"}
	requestBody, _ := json.Marshal(email)
	createdAt := time.Now().UTC().Truncate(time.Second)

	mock.ExpectQuery(regexp.QuoteMeta(`
			SELECT l.id, l.workspace_id, l.inbox_message_id, l.created_at, p.request_body
			FROM message_request_logs l
			INNER JOIN message_request_payloads p ON p.log_id = l.id
			WHERE l.workspace_id = $1 AND l.channel_type = $2
		`)).
		WithArgs(workspaceID, payloadInboxChannelEmail, payloadInboxDefaultPageSize+1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "workspace_id", "inbox_message_id", "created_at", "request_body"}).
			AddRow("log-1", workspaceID, "inbox-1", createdAt, string(requestBody)))

	page := repo.ListEmailsForWorkspace(workspaceID, 0, "")
	if len(page.Items) != 1 {
		t.Fatalf("items len = %d, want 1", len(page.Items))
	}
	if page.Items[0].ID != "inbox-1" {
		t.Fatalf("id = %q, want inbox-1", page.Items[0].ID)
	}
	if page.Items[0].Email.Subject != "Hello" {
		t.Fatalf("subject = %q", page.Items[0].Email.Subject)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestPayloadInboxReader_DeleteEmailByIDForWorkspace(t *testing.T) {
	t.Parallel()

	client, mock, _ := newMockPgClient(t)
	repo := NewPayloadInboxReader(client, nil)

	mock.ExpectExec(regexp.QuoteMeta(`
		DELETE FROM message_request_payloads
		WHERE log_id IN (
			SELECT l.id
			FROM message_request_logs l
			WHERE l.workspace_id = $1 AND l.channel_type = $2
			  AND (l.inbox_message_id = $3 OR l.id::text = $3)
		)
	`)).
		WithArgs("ws-1", payloadInboxChannelEmail, "inbox-1").
		WillReturnResult(sqlmock.NewResult(0, 1))

	if !repo.DeleteEmailByIDForWorkspace("inbox-1", "ws-1") {
		t.Fatal("expected delete success")
	}
}

func TestPayloadInboxReader_StatsForWorkspace(t *testing.T) {
	t.Parallel()

	client, mock, _ := newMockPgClient(t)
	repo := NewPayloadInboxReader(client, nil)

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT l.channel_type, COUNT(*)
		FROM message_request_logs l
		INNER JOIN message_request_payloads p ON p.log_id = l.id
		WHERE l.workspace_id = $1
		GROUP BY l.channel_type
	`)).
		WithArgs("ws-1").
		WillReturnRows(sqlmock.NewRows([]string{"channel_type", "count"}).
			AddRow(payloadInboxChannelEmail, 3))

	stats := repo.StatsForWorkspace("ws-1")
	if stats["emails"] != 3 || stats["total"] != 3 {
		t.Fatalf("stats = %v", stats)
	}
}
