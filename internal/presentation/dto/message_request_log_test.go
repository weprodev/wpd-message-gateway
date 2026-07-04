package dto

import (
	"testing"
	"time"

	"github.com/weprodev/wpd-message-gateway/internal/core/domain"
)

func TestMessageRequestLogFromDomain(t *testing.T) {
	t.Parallel()

	createdAt := time.Date(2026, 7, 4, 10, 0, 0, 0, time.UTC)
	row := domain.MessageRequestLogWithSource{
		MessageRequestLog: domain.MessageRequestLog{
			ID:             "log-1",
			WorkspaceID:    "ws-1",
			ChannelType:    "email",
			HTTPMethod:     "POST",
			StatusCode:     200,
			Endpoint:       "/v1/email",
			ProviderName:   "memory",
			InboxMessageID: "inbox-1",
			CreatedAt:      createdAt,
		},
		SourceName: "Demo Key",
		ClientID:   "demo-client",
	}

	got := MessageRequestLogFromDomain(row)
	if got.ID != "log-1" || got.InboxMessageID != "inbox-1" || got.SourceName != "Demo Key" {
		t.Fatalf("unexpected mapping: %+v", got)
	}

	list := MessageRequestLogListFromDomain([]domain.MessageRequestLogWithSource{row}, 1)
	if len(list.Items) != 1 || list.Total != 1 {
		t.Fatalf("unexpected list mapping: %+v", list)
	}
}
