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
// This is used by the memory provider internally and by handlers that need access to stored messages.
func GetStore() *Store {
	globalStoreOnce.Do(func() {
		globalStore = NewStore()
	})
	return globalStore
}

// StoredEmail wraps an email with metadata for storage.
type StoredEmail struct {
	ID          string           `json:"id"`
	WorkspaceID string           `json:"workspace_id,omitempty"`
	CreatedAt   time.Time        `json:"created_at"`
	Email       *contracts.Email `json:"email"`
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

// StoredOTP wraps an OTP with metadata for storage.
type StoredOTP struct {
	ID          string         `json:"id"`
	WorkspaceID string         `json:"workspace_id,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	OTP         *contracts.OTP `json:"otp"`
}

// Store implements an in-memory message store for all message types.
type Store struct {
	mu     sync.RWMutex
	emails []*StoredEmail
	sms    []*StoredSMS
	pushes []*StoredPush
	chats  []*StoredChat
	otps   []*StoredOTP
}

// NewStore creates a new in-memory store.
func NewStore() *Store {
	return &Store{
		emails: make([]*StoredEmail, 0),
		sms:    make([]*StoredSMS, 0),
		pushes: make([]*StoredPush, 0),
		chats:  make([]*StoredChat, 0),
		otps:   make([]*StoredOTP, 0),
	}
}

// Emails returns a copy of all stored emails.
func (s *Store) Emails() []*StoredEmail {
	s.mu.RLock()
	defer s.mu.RUnlock()
	msgs := make([]*StoredEmail, len(s.emails))
	copy(msgs, s.emails)
	return msgs
}

// EmailsForWorkspace returns stored emails tagged with the given workspace ID.
func (s *Store) EmailsForWorkspace(workspaceID string) []*StoredEmail {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var out []*StoredEmail
	for _, e := range s.emails {
		if e != nil && e.WorkspaceID == workspaceID {
			out = append(out, e)
		}
	}
	return out
}

// EmailByIDForWorkspace returns an email by ID only if it belongs to the workspace.
func (s *Store) EmailByIDForWorkspace(id, workspaceID string) *StoredEmail {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, e := range s.emails {
		if e != nil && e.ID == id && e.WorkspaceID == workspaceID {
			return e
		}
	}
	return nil
}

// DeleteEmailByIDForWorkspace deletes an email if it belongs to the workspace.
func (s *Store) DeleteEmailByIDForWorkspace(id, workspaceID string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, e := range s.emails {
		if e != nil && e.ID == id && e.WorkspaceID == workspaceID {
			s.emails = append(s.emails[:i], s.emails[i+1:]...)
			return true
		}
	}
	return false
}

// StatsForWorkspace returns counts for messages belonging to workspace.
func (s *Store) StatsForWorkspace(workspaceID string) map[string]int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	emails, otps := 0, 0
	for _, e := range s.emails {
		if e != nil && e.WorkspaceID == workspaceID {
			emails++
		}
	}
	for _, o := range s.otps {
		if o != nil && o.WorkspaceID == workspaceID {
			otps++
		}
	}
	return map[string]int{
		"emails": emails,
		"sms":    0,
		"push":   0,
		"chat":   0,
		"otp":    otps,
		"total":  emails + otps,
	}
}

// ClearWorkspace removes in-memory messages for a single workspace.
func (s *Store) ClearWorkspace(workspaceID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	emails := s.emails[:0]
	for _, e := range s.emails {
		if e == nil || e.WorkspaceID != workspaceID {
			emails = append(emails, e)
		}
	}
	s.emails = emails
	otps := s.otps[:0]
	for _, o := range s.otps {
		if o == nil || o.WorkspaceID != workspaceID {
			otps = append(otps, o)
		}
	}
	s.otps = otps
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

// OTPs returns a copy of all stored OTPs.
func (s *Store) OTPs() []*StoredOTP {
	s.mu.RLock()
	defer s.mu.RUnlock()
	msgs := make([]*StoredOTP, len(s.otps))
	copy(msgs, s.otps)
	return msgs
}

// OTPByID returns a stored OTP by its ID, or nil if not found.
func (s *Store) OTPByID(id string) *StoredOTP {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, o := range s.otps {
		if o.ID == id {
			return o
		}
	}
	return nil
}

// DeleteOTPByID deletes an OTP by ID. Returns true if deleted.
func (s *Store) DeleteOTPByID(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, o := range s.otps {
		if o.ID == id {
			s.otps = append(s.otps[:i], s.otps[i+1:]...)
			return true
		}
	}
	return false
}

// AddOTP adds a stored OTP.
func (s *Store) AddOTP(stored *StoredOTP) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.otps = append(s.otps, stored)
}

// OTPsForWorkspace returns stored OTPs tagged with the given workspace ID.
func (s *Store) OTPsForWorkspace(workspaceID string) []*StoredOTP {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var out []*StoredOTP
	for _, o := range s.otps {
		if o != nil && o.WorkspaceID == workspaceID {
			out = append(out, o)
		}
	}
	return out
}

// OTPByIDForWorkspace returns an OTP by ID only if it belongs to the workspace.
func (s *Store) OTPByIDForWorkspace(id, workspaceID string) *StoredOTP {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, o := range s.otps {
		if o != nil && o.ID == id && o.WorkspaceID == workspaceID {
			return o
		}
	}
	return nil
}

// DeleteOTPByIDForWorkspace deletes an OTP if it belongs to the workspace.
func (s *Store) DeleteOTPByIDForWorkspace(id, workspaceID string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, o := range s.otps {
		if o != nil && o.ID == id && o.WorkspaceID == workspaceID {
			s.otps = append(s.otps[:i], s.otps[i+1:]...)
			return true
		}
	}
	return false
}

// Count returns the total number of stored messages across all types.
func (s *Store) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.emails) + len(s.sms) + len(s.pushes) + len(s.chats) + len(s.otps)
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
		"otp":    len(s.otps),
		"total":  len(s.emails) + len(s.sms) + len(s.pushes) + len(s.chats) + len(s.otps),
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
	s.otps = make([]*StoredOTP, 0)
}
