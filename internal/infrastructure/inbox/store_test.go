package inbox

import (
	"context"
	"testing"

	"github.com/weprodev/wpd-message-gateway/pkg/contracts"
)

func TestStoreWorkspaceIsolation(t *testing.T) {
	store := NewStore()
	ctx := context.Background()

	_, err := store.WriteEmail(ctx, "ws-a", contracts.Email{Subject: "a", To: []string{"a@test.com"}})
	if err != nil {
		t.Fatalf("WriteEmail: %v", err)
	}
	_, err = store.WriteSMS(ctx, "ws-a", contracts.SMS{To: []string{"+1"}, Message: "hi"})
	if err != nil {
		t.Fatalf("WriteSMS: %v", err)
	}
	_, err = store.WriteEmail(ctx, "ws-b", contracts.Email{Subject: "b", To: []string{"b@test.com"}})
	if err != nil {
		t.Fatalf("WriteEmail ws-b: %v", err)
	}

	stats := store.StatsForWorkspace("ws-a")
	if stats["emails"] != 1 || stats["sms"] != 1 || stats["total"] != 2 {
		t.Fatalf("ws-a stats: %v", stats)
	}
	if len(store.EmailsForWorkspace("ws-b")) != 1 {
		t.Fatalf("ws-b emails: %d", len(store.EmailsForWorkspace("ws-b")))
	}

	email, ok := store.EmailByIDForWorkspace(store.EmailsForWorkspace("ws-a")[0].ID, "ws-a")
	if !ok || email.Email.Subject != "a" {
		t.Fatalf("lookup email: ok=%v subject=%q", ok, email.Email.Subject)
	}

	store.ClearWorkspace("ws-a")
	if len(store.EmailsForWorkspace("ws-a")) != 0 {
		t.Fatal("expected ws-a cleared")
	}
	if len(store.EmailsForWorkspace("ws-b")) != 1 {
		t.Fatal("expected ws-b preserved")
	}
}
