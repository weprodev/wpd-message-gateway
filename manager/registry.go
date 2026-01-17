package manager

import (
	"sync"

	"github.com/weprodev/wpd-message-gateway/contracts"
	"github.com/weprodev/wpd-message-gateway/providers/memory"
)

// Registry manages provider instances in a thread-safe manner.
// It provides lazy initialization of providers via a factory.
type Registry struct {
	factory ProviderFactory
	mu      sync.RWMutex

	emailProviders map[string]contracts.EmailSender
	smsProviders   map[string]contracts.SMSSender
	pushProviders  map[string]contracts.PushSender
	chatProviders  map[string]contracts.ChatSender

	// Shared memory store for DevBox UI access
	// This is created once and shared across all memory providers
	memoryStore *memory.Provider
}

// NewRegistry creates a new Registry with the given factory.
func NewRegistry(factory ProviderFactory) *Registry {
	return &Registry{
		factory:        factory,
		emailProviders: make(map[string]contracts.EmailSender),
		smsProviders:   make(map[string]contracts.SMSSender),
		pushProviders:  make(map[string]contracts.PushSender),
		chatProviders:  make(map[string]contracts.ChatSender),
		memoryStore:    memory.New(),
	}
}

// GetMemoryStore returns the shared memory store instance.
// This is used by DevBox UI to access stored messages.
func (r *Registry) GetMemoryStore() *memory.Provider {
	return r.memoryStore
}

// GetEmailProvider returns an email provider by name.
// If the provider doesn't exist, it returns nil and false.
func (r *Registry) GetEmailProvider(name string) (contracts.EmailSender, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	provider, ok := r.emailProviders[name]
	return provider, ok
}

// RegisterEmailProvider registers an email provider with the given name.
func (r *Registry) RegisterEmailProvider(name string, provider contracts.EmailSender) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.emailProviders[name] = provider
}

// GetSMSProvider returns an SMS provider by name.
// If the provider doesn't exist, it returns nil and false.
func (r *Registry) GetSMSProvider(name string) (contracts.SMSSender, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	provider, ok := r.smsProviders[name]
	return provider, ok
}

// RegisterSMSProvider registers an SMS provider with the given name.
func (r *Registry) RegisterSMSProvider(name string, provider contracts.SMSSender) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.smsProviders[name] = provider
}

// GetPushProvider returns a push provider by name.
// If the provider doesn't exist, it returns nil and false.
func (r *Registry) GetPushProvider(name string) (contracts.PushSender, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	provider, ok := r.pushProviders[name]
	return provider, ok
}

// RegisterPushProvider registers a push provider with the given name.
func (r *Registry) RegisterPushProvider(name string, provider contracts.PushSender) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.pushProviders[name] = provider
}

// GetChatProvider returns a chat provider by name.
// If the provider doesn't exist, it returns nil and false.
func (r *Registry) GetChatProvider(name string) (contracts.ChatSender, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	provider, ok := r.chatProviders[name]
	return provider, ok
}

// RegisterChatProvider registers a chat provider with the given name.
func (r *Registry) RegisterChatProvider(name string, provider contracts.ChatSender) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.chatProviders[name] = provider
}

// GetFactory returns the factory used by this registry.
func (r *Registry) GetFactory() ProviderFactory {
	return r.factory
}
