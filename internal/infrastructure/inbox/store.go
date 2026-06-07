package inbox

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/weprodev/wpd-message-gateway/internal/core/port"
	"github.com/weprodev/wpd-message-gateway/pkg/contracts"
)

var _ port.InboxReader = (*Store)(nil)
var _ port.InboxWriter = (*Store)(nil)

// Store implements an in-memory message store that satisfies port.InboxReader and port.InboxWriter.
type Store struct {
	mu     sync.RWMutex
	emails []port.StoredEmail
	sms    []port.StoredSMS
	pushes []port.StoredPush
	chats  []port.StoredChat
}

// NewStore creates a new in-memory store.
func NewStore() *Store {
	return &Store{}
}

// StatsForWorkspace returns counts for messages belonging to workspace.
func (s *Store) StatsForWorkspace(workspaceID string) map[string]int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	nEmail := countWorkspace(s.emails, workspaceID)
	nSMS := countWorkspaceSMS(s.sms, workspaceID)
	nPush := countWorkspacePush(s.pushes, workspaceID)
	nChat := countWorkspaceChat(s.chats, workspaceID)
	return map[string]int{
		"emails": nEmail,
		"sms":    nSMS,
		"push":   nPush,
		"chat":   nChat,
		"total":  nEmail + nSMS + nPush + nChat,
	}
}

func countWorkspace(items []port.StoredEmail, workspaceID string) int {
	n := 0
	for _, e := range items {
		if e.WorkspaceID == workspaceID {
			n++
		}
	}
	return n
}

func countWorkspaceSMS(items []port.StoredSMS, workspaceID string) int {
	n := 0
	for _, e := range items {
		if e.WorkspaceID == workspaceID {
			n++
		}
	}
	return n
}

func countWorkspacePush(items []port.StoredPush, workspaceID string) int {
	n := 0
	for _, e := range items {
		if e.WorkspaceID == workspaceID {
			n++
		}
	}
	return n
}

func countWorkspaceChat(items []port.StoredChat, workspaceID string) int {
	n := 0
	for _, e := range items {
		if e.WorkspaceID == workspaceID {
			n++
		}
	}
	return n
}

func (s *Store) EmailsForWorkspace(workspaceID string) []port.StoredEmail {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]port.StoredEmail, 0)
	for _, e := range s.emails {
		if e.WorkspaceID == workspaceID {
			out = append(out, e)
		}
	}
	return out
}

func (s *Store) EmailByIDForWorkspace(id, workspaceID string) (port.StoredEmail, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, e := range s.emails {
		if e.ID == id && e.WorkspaceID == workspaceID {
			return e, true
		}
	}
	return port.StoredEmail{}, false
}

func (s *Store) DeleteEmailByIDForWorkspace(id, workspaceID string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, e := range s.emails {
		if e.ID == id && e.WorkspaceID == workspaceID {
			s.emails = append(s.emails[:i], s.emails[i+1:]...)
			return true
		}
	}
	return false
}

func (s *Store) SMSForWorkspace(workspaceID string) []port.StoredSMS {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]port.StoredSMS, 0)
	for _, e := range s.sms {
		if e.WorkspaceID == workspaceID {
			out = append(out, e)
		}
	}
	return out
}

func (s *Store) SMSByIDForWorkspace(id, workspaceID string) (port.StoredSMS, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, e := range s.sms {
		if e.ID == id && e.WorkspaceID == workspaceID {
			return e, true
		}
	}
	return port.StoredSMS{}, false
}

func (s *Store) DeleteSMSByIDForWorkspace(id, workspaceID string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, e := range s.sms {
		if e.ID == id && e.WorkspaceID == workspaceID {
			s.sms = append(s.sms[:i], s.sms[i+1:]...)
			return true
		}
	}
	return false
}

func (s *Store) PushForWorkspace(workspaceID string) []port.StoredPush {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]port.StoredPush, 0)
	for _, e := range s.pushes {
		if e.WorkspaceID == workspaceID {
			out = append(out, e)
		}
	}
	return out
}

func (s *Store) PushByIDForWorkspace(id, workspaceID string) (port.StoredPush, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, e := range s.pushes {
		if e.ID == id && e.WorkspaceID == workspaceID {
			return e, true
		}
	}
	return port.StoredPush{}, false
}

func (s *Store) DeletePushByIDForWorkspace(id, workspaceID string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, e := range s.pushes {
		if e.ID == id && e.WorkspaceID == workspaceID {
			s.pushes = append(s.pushes[:i], s.pushes[i+1:]...)
			return true
		}
	}
	return false
}

func (s *Store) ChatForWorkspace(workspaceID string) []port.StoredChat {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]port.StoredChat, 0)
	for _, e := range s.chats {
		if e.WorkspaceID == workspaceID {
			out = append(out, e)
		}
	}
	return out
}

func (s *Store) ChatByIDForWorkspace(id, workspaceID string) (port.StoredChat, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, e := range s.chats {
		if e.ID == id && e.WorkspaceID == workspaceID {
			return e, true
		}
	}
	return port.StoredChat{}, false
}

func (s *Store) DeleteChatByIDForWorkspace(id, workspaceID string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, e := range s.chats {
		if e.ID == id && e.WorkspaceID == workspaceID {
			s.chats = append(s.chats[:i], s.chats[i+1:]...)
			return true
		}
	}
	return false
}

func (s *Store) ClearWorkspace(workspaceID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.emails = filterKeepOtherWorkspaces(s.emails, workspaceID)
	s.sms = filterKeepOtherWorkspacesSMS(s.sms, workspaceID)
	s.pushes = filterKeepOtherWorkspacesPush(s.pushes, workspaceID)
	s.chats = filterKeepOtherWorkspacesChat(s.chats, workspaceID)
}

func filterKeepOtherWorkspaces(items []port.StoredEmail, workspaceID string) []port.StoredEmail {
	out := items[:0]
	for _, e := range items {
		if e.WorkspaceID != workspaceID {
			out = append(out, e)
		}
	}
	return out
}

func filterKeepOtherWorkspacesSMS(items []port.StoredSMS, workspaceID string) []port.StoredSMS {
	out := items[:0]
	for _, e := range items {
		if e.WorkspaceID != workspaceID {
			out = append(out, e)
		}
	}
	return out
}

func filterKeepOtherWorkspacesPush(items []port.StoredPush, workspaceID string) []port.StoredPush {
	out := items[:0]
	for _, e := range items {
		if e.WorkspaceID != workspaceID {
			out = append(out, e)
		}
	}
	return out
}

func filterKeepOtherWorkspacesChat(items []port.StoredChat, workspaceID string) []port.StoredChat {
	out := items[:0]
	for _, e := range items {
		if e.WorkspaceID != workspaceID {
			out = append(out, e)
		}
	}
	return out
}

func (s *Store) WriteEmail(_ context.Context, workspaceID string, email contracts.Email) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := uuid.New().String()
	s.emails = append(s.emails, port.StoredEmail{
		ID:          id,
		WorkspaceID: workspaceID,
		CreatedAt:   time.Now(),
		Email:       email,
	})
	return id, nil
}

func (s *Store) WriteSMS(_ context.Context, workspaceID string, sms contracts.SMS) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := uuid.New().String()
	s.sms = append(s.sms, port.StoredSMS{
		ID:          id,
		WorkspaceID: workspaceID,
		CreatedAt:   time.Now(),
		SMS:         sms,
	})
	return id, nil
}

func (s *Store) WritePush(_ context.Context, workspaceID string, push contracts.PushNotification) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := uuid.New().String()
	s.pushes = append(s.pushes, port.StoredPush{
		ID:          id,
		WorkspaceID: workspaceID,
		CreatedAt:   time.Now(),
		Push:        push,
	})
	return id, nil
}

func (s *Store) WriteChat(_ context.Context, workspaceID string, chat contracts.ChatMessage) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := uuid.New().String()
	s.chats = append(s.chats, port.StoredChat{
		ID:          id,
		WorkspaceID: workspaceID,
		CreatedAt:   time.Now(),
		Chat:        chat,
	})
	return id, nil
}
