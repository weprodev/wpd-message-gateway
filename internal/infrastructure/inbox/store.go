package inbox

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/weprodev/wpd-message-gateway/internal/core/port"
	"github.com/weprodev/wpd-message-gateway/pkg/contracts"
)

const (
	defaultInboxPageSize = 50
	maxInboxPageSize     = 200
)

var _ port.InboxReader = (*Store)(nil)
var _ port.InboxWriter = (*Store)(nil)

// Store implements an in-memory message store indexed by workspace for O(1) bucket lookup.
type Store struct {
	mu       sync.RWMutex
	emails   map[string][]port.StoredEmail
	sms      map[string][]port.StoredSMS
	pushes   map[string][]port.StoredPush
	chats    map[string][]port.StoredChat
}

// NewStore creates a new in-memory store.
func NewStore() *Store {
	return &Store{
		emails: make(map[string][]port.StoredEmail),
		sms:    make(map[string][]port.StoredSMS),
		pushes: make(map[string][]port.StoredPush),
		chats:  make(map[string][]port.StoredChat),
	}
}

func clampInboxLimit(limit int) int {
	if limit <= 0 {
		return defaultInboxPageSize
	}
	if limit > maxInboxPageSize {
		return maxInboxPageSize
	}
	return limit
}

// StatsForWorkspace returns counts for messages belonging to workspace.
func (s *Store) StatsForWorkspace(workspaceID string) map[string]int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	nEmail := len(s.emails[workspaceID])
	nSMS := len(s.sms[workspaceID])
	nPush := len(s.pushes[workspaceID])
	nChat := len(s.chats[workspaceID])
	return map[string]int{
		"emails": nEmail,
		"sms":    nSMS,
		"push":   nPush,
		"chat":   nChat,
		"total":  nEmail + nSMS + nPush + nChat,
	}
}

func (s *Store) EmailsForWorkspace(workspaceID string) []port.StoredEmail {
	page := s.ListEmailsForWorkspace(workspaceID, 0, "")
	return page.Items
}

func (s *Store) ListEmailsForWorkspace(workspaceID string, limit int, cursor string) port.InboxEmailPage {
	s.mu.RLock()
	defer s.mu.RUnlock()

	all := s.emails[workspaceID]
	limit = clampInboxLimit(limit)
	start := 0
	if cursor != "" {
		for i, item := range all {
			if item.ID == cursor {
				start = i + 1
				break
			}
		}
	}

	end := start + limit
	hasMore := end < len(all)
	if end > len(all) {
		end = len(all)
	}

	items := make([]port.StoredEmail, 0, end-start)
	if start < len(all) {
		items = append(items, all[start:end]...)
	}
	nextCursor := ""
	if hasMore && len(items) > 0 {
		nextCursor = items[len(items)-1].ID
	}
	return port.InboxEmailPage{
		Items:      items,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}
}

func (s *Store) EmailByIDForWorkspace(id, workspaceID string) (port.StoredEmail, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, e := range s.emails[workspaceID] {
		if e.ID == id {
			return e, true
		}
	}
	return port.StoredEmail{}, false
}

func (s *Store) DeleteEmailByIDForWorkspace(id, workspaceID string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	items := s.emails[workspaceID]
	for i, e := range items {
		if e.ID == id {
			s.emails[workspaceID] = append(items[:i], items[i+1:]...)
			return true
		}
	}
	return false
}

func (s *Store) SMSForWorkspace(workspaceID string) []port.StoredSMS {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]port.StoredSMS(nil), s.sms[workspaceID]...)
}

func (s *Store) SMSByIDForWorkspace(id, workspaceID string) (port.StoredSMS, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, item := range s.sms[workspaceID] {
		if item.ID == id {
			return item, true
		}
	}
	return port.StoredSMS{}, false
}

func (s *Store) DeleteSMSByIDForWorkspace(id, workspaceID string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	items := s.sms[workspaceID]
	for i, e := range items {
		if e.ID == id {
			s.sms[workspaceID] = append(items[:i], items[i+1:]...)
			return true
		}
	}
	return false
}

func (s *Store) PushForWorkspace(workspaceID string) []port.StoredPush {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]port.StoredPush(nil), s.pushes[workspaceID]...)
}

func (s *Store) PushByIDForWorkspace(id, workspaceID string) (port.StoredPush, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, item := range s.pushes[workspaceID] {
		if item.ID == id {
			return item, true
		}
	}
	return port.StoredPush{}, false
}

func (s *Store) DeletePushByIDForWorkspace(id, workspaceID string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	items := s.pushes[workspaceID]
	for i, e := range items {
		if e.ID == id {
			s.pushes[workspaceID] = append(items[:i], items[i+1:]...)
			return true
		}
	}
	return false
}

func (s *Store) ChatForWorkspace(workspaceID string) []port.StoredChat {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]port.StoredChat(nil), s.chats[workspaceID]...)
}

func (s *Store) ChatByIDForWorkspace(id, workspaceID string) (port.StoredChat, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, item := range s.chats[workspaceID] {
		if item.ID == id {
			return item, true
		}
	}
	return port.StoredChat{}, false
}

func (s *Store) DeleteChatByIDForWorkspace(id, workspaceID string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	items := s.chats[workspaceID]
	for i, e := range items {
		if e.ID == id {
			s.chats[workspaceID] = append(items[:i], items[i+1:]...)
			return true
		}
	}
	return false
}

func (s *Store) ClearWorkspace(workspaceID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.emails, workspaceID)
	delete(s.sms, workspaceID)
	delete(s.pushes, workspaceID)
	delete(s.chats, workspaceID)
}

func (s *Store) WriteEmail(_ context.Context, workspaceID string, email contracts.Email) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := uuid.New().String()
	item := port.StoredEmail{
		ID:          id,
		WorkspaceID: workspaceID,
		CreatedAt:   time.Now(),
		Email:       email,
	}
	s.emails[workspaceID] = append([]port.StoredEmail{item}, s.emails[workspaceID]...)
	return id, nil
}

func (s *Store) WriteSMS(_ context.Context, workspaceID string, sms contracts.SMS) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := uuid.New().String()
	item := port.StoredSMS{
		ID:          id,
		WorkspaceID: workspaceID,
		CreatedAt:   time.Now(),
		SMS:         sms,
	}
	s.sms[workspaceID] = append([]port.StoredSMS{item}, s.sms[workspaceID]...)
	return id, nil
}

func (s *Store) WritePush(_ context.Context, workspaceID string, push contracts.PushNotification) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := uuid.New().String()
	item := port.StoredPush{
		ID:          id,
		WorkspaceID: workspaceID,
		CreatedAt:   time.Now(),
		Push:        push,
	}
	s.pushes[workspaceID] = append([]port.StoredPush{item}, s.pushes[workspaceID]...)
	return id, nil
}

func (s *Store) WriteChat(_ context.Context, workspaceID string, chat contracts.ChatMessage) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := uuid.New().String()
	item := port.StoredChat{
		ID:          id,
		WorkspaceID: workspaceID,
		CreatedAt:   time.Now(),
		Chat:        chat,
	}
	s.chats[workspaceID] = append([]port.StoredChat{item}, s.chats[workspaceID]...)
	return id, nil
}
