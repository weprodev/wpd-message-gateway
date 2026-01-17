package memory

import (
	"sync"
	"time"

	"github.com/weprodev/wpd-message-gateway/pkg/contracts"
)

const ProviderName = "memory"

// --- Stored Message Wrappers ---

// StoredEmail wraps an email with metadata for storage.
type StoredEmail struct {
	ID        string           `json:"id"`
	CreatedAt time.Time        `json:"created_at"`
	Email     *contracts.Email `json:"email"`
}

// StoredSMS wraps an SMS with metadata for storage.
type StoredSMS struct {
	ID        string         `json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	SMS       *contracts.SMS `json:"sms"`
}

// StoredPush wraps a push notification with metadata for storage.
type StoredPush struct {
	ID        string                      `json:"id"`
	CreatedAt time.Time                   `json:"created_at"`
	Push      *contracts.PushNotification `json:"push"`
}

// StoredChat wraps a chat message with metadata for storage.
type StoredChat struct {
	ID        string                 `json:"id"`
	CreatedAt time.Time              `json:"created_at"`
	Chat      *contracts.ChatMessage `json:"chat"`
}

// Store implements an in-memory message store for all message types.
type Store struct {
	mu     sync.RWMutex
	emails []*StoredEmail
	sms    []*StoredSMS
	pushes []*StoredPush
	chats  []*StoredChat
}

// NewStore creates a new in-memory store.
func NewStore() *Store {
	return &Store{
		emails: make([]*StoredEmail, 0),
		sms:    make([]*StoredSMS, 0),
		pushes: make([]*StoredPush, 0),
		chats:  make([]*StoredChat, 0),
	}
}

// --- Email Methods ---

// Emails returns a copy of all stored emails.
func (s *Store) Emails() []*StoredEmail {
	s.mu.RLock()
	defer s.mu.RUnlock()
	msgs := make([]*StoredEmail, len(s.emails))
	copy(msgs, s.emails)
	return msgs
}

// EmailByID returns a stored email by its ID, or nil if not found.
func (s *Store) EmailByID(id string) *StoredEmail {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, e := range s.emails {
		if e.ID == id {
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
		if e.ID == id {
			s.emails = append(s.emails[:i], s.emails[i+1:]...)
			return true
		}
	}
	return false
}

// AddEmail adds a stored email.
func (s *Store) AddEmail(stored *StoredEmail) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.emails = append(s.emails, stored)
}

// --- SMS Methods ---

// AllSMS returns a copy of all stored SMS messages.
func (s *Store) AllSMS() []*StoredSMS {
	s.mu.RLock()
	defer s.mu.RUnlock()
	msgs := make([]*StoredSMS, len(s.sms))
	copy(msgs, s.sms)
	return msgs
}

// SMSByID returns a stored SMS by its ID, or nil if not found.
func (s *Store) SMSByID(id string) *StoredSMS {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, msg := range s.sms {
		if msg.ID == id {
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
		if msg.ID == id {
			s.sms = append(s.sms[:i], s.sms[i+1:]...)
			return true
		}
	}
	return false
}

// AddSMS adds a stored SMS.
func (s *Store) AddSMS(stored *StoredSMS) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sms = append(s.sms, stored)
}

// --- Push Methods ---

// Pushes returns a copy of all stored push notifications.
func (s *Store) Pushes() []*StoredPush {
	s.mu.RLock()
	defer s.mu.RUnlock()
	msgs := make([]*StoredPush, len(s.pushes))
	copy(msgs, s.pushes)
	return msgs
}

// PushByID returns a stored push notification by its ID, or nil if not found.
func (s *Store) PushByID(id string) *StoredPush {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, push := range s.pushes {
		if push.ID == id {
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
		if push.ID == id {
			s.pushes = append(s.pushes[:i], s.pushes[i+1:]...)
			return true
		}
	}
	return false
}

// AddPush adds a stored push notification.
func (s *Store) AddPush(stored *StoredPush) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.pushes = append(s.pushes, stored)
}

// --- Chat Methods ---

// Chats returns a copy of all stored chat messages.
func (s *Store) Chats() []*StoredChat {
	s.mu.RLock()
	defer s.mu.RUnlock()
	msgs := make([]*StoredChat, len(s.chats))
	copy(msgs, s.chats)
	return msgs
}

// ChatByID returns a stored chat message by its ID, or nil if not found.
func (s *Store) ChatByID(id string) *StoredChat {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, c := range s.chats {
		if c.ID == id {
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
		if c.ID == id {
			s.chats = append(s.chats[:i], s.chats[i+1:]...)
			return true
		}
	}
	return false
}

// AddChat adds a stored chat message.
func (s *Store) AddChat(stored *StoredChat) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.chats = append(s.chats, stored)
}

// --- General Methods ---

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
	s.emails = make([]*StoredEmail, 0)
	s.sms = make([]*StoredSMS, 0)
	s.pushes = make([]*StoredPush, 0)
	s.chats = make([]*StoredChat, 0)
}
