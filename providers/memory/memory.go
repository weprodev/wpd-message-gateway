package memory

import (
	"sync"
	"time"

	"github.com/weprodev/wpd-message-gateway/config"
	"github.com/weprodev/wpd-message-gateway/contracts"
)

const ProviderName = "memory"

// --- Stored Message Wrappers ---
// These wrap the original message types with ID and timestamp for devbox tracking.

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

// Provider implements an in-memory message interceptor for all message types.
// It is thread-safe and acts as a central store for the devbox.
type Provider struct {
	mu     sync.RWMutex
	emails []*StoredEmail
	sms    []*StoredSMS
	pushes []*StoredPush
	chats  []*StoredChat
}

// New creates a new universal Memory provider.
func New() *Provider {
	return &Provider{
		emails: make([]*StoredEmail, 0),
		sms:    make([]*StoredSMS, 0),
		pushes: make([]*StoredPush, 0),
		chats:  make([]*StoredChat, 0),
	}
}

// Name returns the provider name.
func (p *Provider) Name() string {
	return ProviderName
}

// --- Email Methods ---

// Emails returns a copy of all stored emails.
func (p *Provider) Emails() []*StoredEmail {
	p.mu.RLock()
	defer p.mu.RUnlock()
	msgs := make([]*StoredEmail, len(p.emails))
	copy(msgs, p.emails)
	return msgs
}

// EmailByID returns a stored email by its ID, or nil if not found.
func (p *Provider) EmailByID(id string) *StoredEmail {
	p.mu.RLock()
	defer p.mu.RUnlock()
	for _, e := range p.emails {
		if e.ID == id {
			return e
		}
	}
	return nil
}

// DeleteEmailByID deletes an email by ID. Returns true if deleted.
func (p *Provider) DeleteEmailByID(id string) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	for i, e := range p.emails {
		if e.ID == id {
			p.emails = append(p.emails[:i], p.emails[i+1:]...)
			return true
		}
	}
	return false
}

// addEmail adds a stored email (called by EmailProvider.Send).
func (p *Provider) addEmail(stored *StoredEmail) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.emails = append(p.emails, stored)
}

// --- SMS Methods ---

// SMS returns a copy of all stored SMS messages.
func (p *Provider) SMS() []*StoredSMS {
	p.mu.RLock()
	defer p.mu.RUnlock()
	msgs := make([]*StoredSMS, len(p.sms))
	copy(msgs, p.sms)
	return msgs
}

// SMSByID returns a stored SMS by its ID, or nil if not found.
func (p *Provider) SMSByID(id string) *StoredSMS {
	p.mu.RLock()
	defer p.mu.RUnlock()
	for _, s := range p.sms {
		if s.ID == id {
			return s
		}
	}
	return nil
}

// DeleteSMSByID deletes an SMS by ID. Returns true if deleted.
func (p *Provider) DeleteSMSByID(id string) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	for i, s := range p.sms {
		if s.ID == id {
			p.sms = append(p.sms[:i], p.sms[i+1:]...)
			return true
		}
	}
	return false
}

// addSMS adds a stored SMS (called by SMSProvider.Send).
func (p *Provider) addSMS(stored *StoredSMS) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.sms = append(p.sms, stored)
}

// --- Push Methods ---

// Pushes returns a copy of all stored push notifications.
func (p *Provider) Pushes() []*StoredPush {
	p.mu.RLock()
	defer p.mu.RUnlock()
	msgs := make([]*StoredPush, len(p.pushes))
	copy(msgs, p.pushes)
	return msgs
}

// PushByID returns a stored push notification by its ID, or nil if not found.
func (p *Provider) PushByID(id string) *StoredPush {
	p.mu.RLock()
	defer p.mu.RUnlock()
	for _, push := range p.pushes {
		if push.ID == id {
			return push
		}
	}
	return nil
}

// DeletePushByID deletes a push notification by ID. Returns true if deleted.
func (p *Provider) DeletePushByID(id string) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	for i, push := range p.pushes {
		if push.ID == id {
			p.pushes = append(p.pushes[:i], p.pushes[i+1:]...)
			return true
		}
	}
	return false
}

// addPush adds a stored push notification (called by PushProvider.Send).
func (p *Provider) addPush(stored *StoredPush) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.pushes = append(p.pushes, stored)
}

// --- Chat Methods ---

// Chats returns a copy of all stored chat messages.
func (p *Provider) Chats() []*StoredChat {
	p.mu.RLock()
	defer p.mu.RUnlock()
	msgs := make([]*StoredChat, len(p.chats))
	copy(msgs, p.chats)
	return msgs
}

// ChatByID returns a stored chat message by its ID, or nil if not found.
func (p *Provider) ChatByID(id string) *StoredChat {
	p.mu.RLock()
	defer p.mu.RUnlock()
	for _, c := range p.chats {
		if c.ID == id {
			return c
		}
	}
	return nil
}

// DeleteChatByID deletes a chat message by ID. Returns true if deleted.
func (p *Provider) DeleteChatByID(id string) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	for i, c := range p.chats {
		if c.ID == id {
			p.chats = append(p.chats[:i], p.chats[i+1:]...)
			return true
		}
	}
	return false
}

// addChat adds a stored chat message (called by ChatProvider.Send).
func (p *Provider) addChat(stored *StoredChat) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.chats = append(p.chats, stored)
}

// --- General Methods ---

// Count returns the total number of stored messages across all types.
func (p *Provider) Count() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.emails) + len(p.sms) + len(p.pushes) + len(p.chats)
}

// Stats returns message counts by type.
func (p *Provider) Stats() map[string]int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return map[string]int{
		"emails": len(p.emails),
		"sms":    len(p.sms),
		"push":   len(p.pushes),
		"chat":   len(p.chats),
		"total":  len(p.emails) + len(p.sms) + len(p.pushes) + len(p.chats),
	}
}

// Clear removes all stored messages.
func (p *Provider) Clear() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.emails = make([]*StoredEmail, 0)
	p.sms = make([]*StoredSMS, 0)
	p.pushes = make([]*StoredPush, 0)
	p.chats = make([]*StoredChat, 0)
}

// --- Factory Methods for Typed Adapters ---

// EmailProvider returns a typed adapter implementing contracts.EmailSender.
// Includes optional Mailpit forwarding when mailpit.enabled is true in config.
func (p *Provider) EmailProvider(mailpitCfg config.MailpitConfig) *EmailProvider {
	return &EmailProvider{
		store:         p,
		smtpForwarder: newSMTPForwarder(mailpitCfg),
	}
}

// SMSProvider returns a typed adapter implementing contracts.SMSSender.
func (p *Provider) SMSProvider() *SMSProvider {
	return &SMSProvider{store: p}
}

// PushProvider returns a typed adapter implementing contracts.PushSender.
func (p *Provider) PushProvider() *PushProvider {
	return &PushProvider{store: p}
}

// ChatProvider returns a typed adapter implementing contracts.ChatSender.
func (p *Provider) ChatProvider() *ChatProvider {
	return &ChatProvider{store: p}
}
