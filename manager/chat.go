package manager

import (
	"context"

	"github.com/weprodev/wpd-message-gateway/contracts"
	msgerrors "github.com/weprodev/wpd-message-gateway/errors"
)

// Chat returns the default chat sender.
func (m *Manager) Chat() (contracts.ChatSender, error) {
	providerName := m.config.DefaultChatProvider()
	if providerName == "" {
		return nil, msgerrors.NewProviderNotFoundError("chat", "default (none configured)")
	}
	return m.ChatProvider(providerName)
}

// ChatProvider returns a specific chat provider by name.
// Providers are created lazily if they don't exist yet.
func (m *Manager) ChatProvider(name string) (contracts.ChatSender, error) {
	// Try to get from registry first
	provider, ok := m.registry.GetChatProvider(name)
	if ok {
		return provider, nil
	}

	// Provider doesn't exist, try to create it lazily
	memStore := m.registry.GetMemoryStore()
	if err := m.ensureChatProvider(name, memStore); err != nil {
		return nil, msgerrors.NewProviderNotFoundError("chat", name)
	}

	// Get the newly created provider
	provider, ok = m.registry.GetChatProvider(name)
	if !ok {
		return nil, msgerrors.NewProviderNotFoundError("chat", name)
	}

	return provider, nil
}

// SendChat sends a chat message using the default provider.
func (m *Manager) SendChat(ctx context.Context, message *contracts.ChatMessage) (*contracts.SendResult, error) {
	provider, err := m.Chat()
	if err != nil {
		return nil, err
	}
	return provider.Send(ctx, message)
}

// SendChatWith sends a chat message using a specific provider.
func (m *Manager) SendChatWith(ctx context.Context, providerName string, message *contracts.ChatMessage) (*contracts.SendResult, error) {
	provider, err := m.ChatProvider(providerName)
	if err != nil {
		return nil, err
	}
	return provider.Send(ctx, message)
}

// RegisterChatProvider registers a custom chat provider.
func (m *Manager) RegisterChatProvider(name string, provider contracts.ChatSender) {
	m.registry.RegisterChatProvider(name, provider)
}

// AvailableChatProviders returns the names of all registered chat providers.
func (m *Manager) AvailableChatProviders() []string {
	providers := make([]string, 0)

	// Check common providers
	commonProviders := []string{"memory"}
	for _, name := range commonProviders {
		if _, ok := m.registry.GetChatProvider(name); ok {
			providers = append(providers, name)
		}
	}

	// Check configured providers
	for name := range m.config.ChatProviders {
		if _, ok := m.registry.GetChatProvider(name); ok {
			found := false
			for _, p := range providers {
				if p == name {
					found = true
					break
				}
			}
			if !found {
				providers = append(providers, name)
			}
		}
	}

	return providers
}
