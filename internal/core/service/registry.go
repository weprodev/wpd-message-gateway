package service

import (
	"sync"

	"github.com/weprodev/wpd-message-gateway/internal/core/port"
)

// Registry manages provider instances in a thread-safe manner.
type Registry struct {
	mu sync.RWMutex

	emailProviders map[string]port.EmailSender
	smsProviders   map[string]port.SMSSender
	pushProviders  map[string]port.PushSender
	chatProviders  map[string]port.ChatSender
}

// NewRegistry creates a new Registry.
func NewRegistry() *Registry {
	return &Registry{
		emailProviders: make(map[string]port.EmailSender),
		smsProviders:   make(map[string]port.SMSSender),
		pushProviders:  make(map[string]port.PushSender),
		chatProviders:  make(map[string]port.ChatSender),
	}
}

// GetEmailProvider returns an email provider by name.
func (r *Registry) GetEmailProvider(name string) (port.EmailSender, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	provider, ok := r.emailProviders[name]
	return provider, ok
}

// RegisterEmailProvider registers an email provider with the given name.
func (r *Registry) RegisterEmailProvider(name string, provider port.EmailSender) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.emailProviders[name] = provider
}

// GetSMSProvider returns an SMS provider by name.
func (r *Registry) GetSMSProvider(name string) (port.SMSSender, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	provider, ok := r.smsProviders[name]
	return provider, ok
}

// RegisterSMSProvider registers an SMS provider with the given name.
func (r *Registry) RegisterSMSProvider(name string, provider port.SMSSender) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.smsProviders[name] = provider
}

// GetPushProvider returns a push provider by name.
func (r *Registry) GetPushProvider(name string) (port.PushSender, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	provider, ok := r.pushProviders[name]
	return provider, ok
}

// RegisterPushProvider registers a push provider with the given name.
func (r *Registry) RegisterPushProvider(name string, provider port.PushSender) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.pushProviders[name] = provider
}

// GetChatProvider returns a chat provider by name.
func (r *Registry) GetChatProvider(name string) (port.ChatSender, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	provider, ok := r.chatProviders[name]
	return provider, ok
}

// RegisterChatProvider registers a chat provider with the given name.
func (r *Registry) RegisterChatProvider(name string, provider port.ChatSender) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.chatProviders[name] = provider
}
