package inbox

import (
	"context"
	"fmt"
	"testing"

	"github.com/weprodev/wpd-message-gateway/pkg/contracts"
)

func TestStore_ListSMSForWorkspace_pagination(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	store := NewStore()
	workspaceID := "ws-paginate"

	var ids []string
	for i := 0; i < 5; i++ {
		id, err := store.WriteSMS(ctx, workspaceID, contracts.SMS{
			To:      []string{"+15550001111"},
			Message: fmt.Sprintf("msg-%d", i),
		})
		if err != nil {
			t.Fatalf("WriteSMS %d: %v", i, err)
		}
		ids = append(ids, id)
	}

	// Newest first: ids[4], ids[3], ids[2], ids[1], ids[0]
	page1 := store.ListSMSForWorkspace(workspaceID, 2, "")
	if len(page1.Items) != 2 {
		t.Fatalf("page1 len = %d, want 2", len(page1.Items))
	}
	if page1.Items[0].ID != ids[4] || page1.Items[1].ID != ids[3] {
		t.Fatalf("page1 order = [%s, %s], want [%s, %s]", page1.Items[0].ID, page1.Items[1].ID, ids[4], ids[3])
	}
	if !page1.HasMore || page1.NextCursor != ids[3] {
		t.Fatalf("page1 cursor/hasMore = (%q, %v), want (%q, true)", page1.NextCursor, page1.HasMore, ids[3])
	}

	page2 := store.ListSMSForWorkspace(workspaceID, 2, page1.NextCursor)
	if len(page2.Items) != 2 {
		t.Fatalf("page2 len = %d, want 2", len(page2.Items))
	}
	if page2.Items[0].ID != ids[2] || page2.Items[1].ID != ids[1] {
		t.Fatalf("page2 order = [%s, %s], want [%s, %s]", page2.Items[0].ID, page2.Items[1].ID, ids[2], ids[1])
	}
	if !page2.HasMore || page2.NextCursor != ids[1] {
		t.Fatalf("page2 cursor/hasMore = (%q, %v), want (%q, true)", page2.NextCursor, page2.HasMore, ids[1])
	}

	page3 := store.ListSMSForWorkspace(workspaceID, 2, page2.NextCursor)
	if len(page3.Items) != 1 {
		t.Fatalf("page3 len = %d, want 1", len(page3.Items))
	}
	if page3.Items[0].ID != ids[0] {
		t.Fatalf("page3 item = %s, want %s", page3.Items[0].ID, ids[0])
	}
	if page3.HasMore || page3.NextCursor != "" {
		t.Fatalf("page3 should be last page, got cursor=%q hasMore=%v", page3.NextCursor, page3.HasMore)
	}
}

func TestStore_ListSMSForWorkspace_defaultLimit(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	store := NewStore()
	workspaceID := "ws-default-limit"

	for i := 0; i < defaultInboxPageSize+1; i++ {
		if _, err := store.WriteSMS(ctx, workspaceID, contracts.SMS{To: []string{"+1"}, Message: "x"}); err != nil {
			t.Fatalf("WriteSMS: %v", err)
		}
	}

	page := store.ListSMSForWorkspace(workspaceID, 0, "")
	if len(page.Items) != defaultInboxPageSize {
		t.Fatalf("default page size = %d, want %d", len(page.Items), defaultInboxPageSize)
	}
	if !page.HasMore {
		t.Fatal("expected hasMore with default limit")
	}
}
