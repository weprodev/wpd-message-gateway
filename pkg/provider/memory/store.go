package memory

import (
	"sync"
	"time"

	"github.com/weprodev/wpd-message-gateway/pkg/contracts"
)

const ProviderName = "memory"

var (
	globalStore     *Store
	globalStoreOnce sync.Once
)

// GetStore returns the singleton memory store instance.
func GetStore() *Store {
	globalStoreOnce.Do(func() {
		globalStore = NewStore()
	})
	return globalStore
}

// SentEmail wraps a sent email for the in-memory SDK store.
type SentEmail struct {
	ID        string          `json:"id"`
	CreatedAt time.Time       `json:"created_at"`
	Email     contracts.Email `json:"email"`
}

// SentSMS wraps a sent SMS for the in-memory SDK store.
type SentSMS struct {
	ID        string        `json:"id"`
	CreatedAt time.Time     `json:"created_at"`
	SMS       contracts.SMS `json:"sms"`
}

// SentPush wraps a sent push notification for the in-memory SDK store.
type SentPush struct {
	ID        string                     `json:"id"`
	CreatedAt time.Time                  `json:"created_at"`
	Push      contracts.PushNotification `json:"push"`
}

// SentChat wraps a sent chat message for the in-memory SDK store.
type SentChat struct {
	ID        string                `json:"id"`
	CreatedAt time.Time             `json:"created_at"`
	Chat      contracts.ChatMessage `json:"chat"`
}

// Store implements an in-memory message store for SDK-level verification/testing.
type Store struct {
	mu     sync.RWMutex
	emails []*SentEmail
	sms    []*SentSMS
	pushes []*SentPush
	chats  []*SentChat
}

// NewStore creates a new in-memory store.
func NewStore() *Store {
	return &Store{
		emails: make([]*SentEmail, 0),
		sms:    make([]*SentSMS, 0),
		pushes: make([]*SentPush, 0),
		chats:  make([]*SentChat, 0),
	}
}

// Emails returns a copy of all stored emails.
func (s *Store) Emails() []*SentEmail {
	s.mu.RLock()
	defer s.mu.RUnlock()
	msgs := make([]*SentEmail, len(s.emails))
	copy(msgs, s.emails)
	return msgs
}

// EmailByID returns a stored email by its ID, or nil if not found.
func (s *Store) EmailByID(id string) *SentEmail {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, e := range s.emails {
		if e != nil && e.ID == id {
			return e
		}
	}
	return nil
}

// DeleteEmailByID deletes an email by ID. Returns true if deleted.
func (s *Store) DeleteEmailByID(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, e := range s.emails {
		if e != nil && e.ID == id {
			s.emails = append(s.emails[:i], s.emails[i+1:]...)
			return true
		}
	}
	return false
}

// AddEmail adds a stored email.
func (s *Store) AddEmail(stored SentEmail) {
	s.mu.Lock()
	defer s.mu.Unlock()
	heap := new(SentEmail)
	*heap = stored
	s.emails = append(s.emails, heap)
}

// AllSMS returns a copy of all stored SMS messages.
func (s *Store) AllSMS() []*SentSMS {
	s.mu.RLock()
	defer s.mu.RUnlock()
	msgs := make([]*SentSMS, len(s.sms))
	copy(msgs, s.sms)
	return msgs
}

// SMSByID returns a stored SMS by its ID, or nil if not found.
func (s *Store) SMSByID(id string) *SentSMS {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, msg := range s.sms {
		if msg != nil && msg.ID == id {
			return msg
		}
	}
	return nil
}

// DeleteSMSByID deletes an SMS by ID. Returns true if deleted.
func (s *Store) DeleteSMSByID(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, msg := range s.sms {
		if msg != nil && msg.ID == id {
			s.sms = append(s.sms[:i], s.sms[i+1:]...)
			return true
		}
	}
	return false
}

// AddSMS adds a stored SMS.
func (s *Store) AddSMS(stored SentSMS) {
	s.mu.Lock()
	defer s.mu.Unlock()
	heap := new(SentSMS)
	*heap = stored
	s.sms = append(s.sms, heap)
}

// Pushes returns a copy of all stored push notifications.
func (s *Store) Pushes() []*SentPush {
	s.mu.RLock()
	defer s.mu.RUnlock()
	msgs := make([]*SentPush, len(s.pushes))
	copy(msgs, s.pushes)
	return msgs
}

// PushByID returns a stored push notification by its ID, or nil if not found.
func (s *Store) PushByID(id string) *SentPush {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, push := range s.pushes {
		if push != nil && push.ID == id {
			return push
		}
	}
	return nil
}

// DeletePushByID deletes a push notification by ID. Returns true if deleted.
func (s *Store) DeletePushByID(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, push := range s.pushes {
		if push != nil && push.ID == id {
			s.pushes = append(s.pushes[:i], s.pushes[i+1:]...)
			return true
		}
	}
	return false
}

// AddPush adds a stored push notification.
func (s *Store) AddPush(stored SentPush) {
	s.mu.Lock()
	defer s.mu.Unlock()
	heap := new(SentPush)
	*heap = stored
	s.pushes = append(s.pushes, heap)
}

// Chats returns a copy of all stored chat messages.
func (s *Store) Chats() []*SentChat {
	s.mu.RLock()
	defer s.mu.RUnlock()
	msgs := make([]*SentChat, len(s.chats))
	copy(msgs, s.chats)
	return msgs
}

// ChatByID returns a stored chat message by its ID, or nil if not found.
func (s *Store) ChatByID(id string) *SentChat {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, c := range s.chats {
		if c != nil && c.ID == id {
			return c
		}
	}
	return nil
}

// DeleteChatByID deletes a chat message by ID. Returns true if deleted.
func (s *Store) DeleteChatByID(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, c := range s.chats {
		if c != nil && c.ID == id {
			s.chats = append(s.chats[:i], s.chats[i+1:]...)
			return true
		}
	}
	return false
}

// AddChat adds a stored chat message.
func (s *Store) AddChat(stored SentChat) {
	s.mu.Lock()
	defer s.mu.Unlock()
	heap := new(SentChat)
	*heap = stored
	s.chats = append(s.chats, heap)
}

// Count returns the total number of stored messages across all types.
func (s *Store) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.emails) + len(s.sms) + len(s.pushes) + len(s.chats)
}

// Stats returns message counts by type.
func (s *Store) Stats() map[string]int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return map[string]int{
		"emails": len(s.emails),
		"sms":    len(s.sms),
		"push":   len(s.pushes),
		"chat":   len(s.chats),
		"total":  len(s.emails) + len(s.sms) + len(s.pushes) + len(s.chats),
	}
}

// Clear removes all stored messages.
func (s *Store) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.emails = make([]*SentEmail, 0)
	s.sms = make([]*SentSMS, 0)
	s.pushes = make([]*SentPush, 0)
	s.chats = make([]*SentChat, 0)
}
