package manager

import (
	"context"

	"github.com/weprodev/wpd-message-gateway/contracts"
	msgerrors "github.com/weprodev/wpd-message-gateway/errors"
)

// Chat returns the default chat sender
func (m *Manager) Chat() (contracts.ChatSender, error) {
	if m.config.DefaultChatProvider == "" {
		return nil, msgerrors.NewProviderNotFoundError("chat", "default (none configured)")
	}
	return m.ChatProvider(m.config.DefaultChatProvider)
}

// ChatProvider returns a specific chat provider by name
func (m *Manager) ChatProvider(name string) (contracts.ChatSender, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	provider, ok := m.chatProviders[name]
	if !ok {
		return nil, msgerrors.NewProviderNotFoundError("chat", name)
	}
	return provider, nil
}

// SendChat sends a chat message using the default provider
func (m *Manager) SendChat(ctx context.Context, message *contracts.ChatMessage) (*contracts.SendResult, error) {
	provider, err := m.Chat()
	if err != nil {
		return nil, err
	}
	return provider.Send(ctx, message)
}

// SendChatWith sends a chat message using a specific provider
func (m *Manager) SendChatWith(ctx context.Context, providerName string, message *contracts.ChatMessage) (*contracts.SendResult, error) {
	provider, err := m.ChatProvider(providerName)
	if err != nil {
		return nil, err
	}
	return provider.Send(ctx, message)
}

// RegisterChatProvider registers a custom chat provider
func (m *Manager) RegisterChatProvider(name string, provider contracts.ChatSender) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.chatProviders[name] = provider
}

// AvailableChatProviders returns the names of all registered chat providers
func (m *Manager) AvailableChatProviders() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	names := make([]string, 0, len(m.chatProviders))
	for name := range m.chatProviders {
		names = append(names, name)
	}
	return names
}
