package manager

import (
	"context"

	"github.com/weprodev/wpd-message-gateway/contracts"
	msgerrors "github.com/weprodev/wpd-message-gateway/errors"
)

// Push returns the default push notification sender
func (m *Manager) Push() (contracts.PushSender, error) {
	if m.config.DefaultPushProvider == "" {
		return nil, msgerrors.NewProviderNotFoundError("push", "default (none configured)")
	}
	return m.PushProvider(m.config.DefaultPushProvider)
}

// PushProvider returns a specific push notification provider by name
func (m *Manager) PushProvider(name string) (contracts.PushSender, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	provider, ok := m.pushProviders[name]
	if !ok {
		return nil, msgerrors.NewProviderNotFoundError("push", name)
	}
	return provider, nil
}

// SendPush sends a push notification using the default provider
func (m *Manager) SendPush(ctx context.Context, notification *contracts.PushNotification) (*contracts.SendResult, error) {
	provider, err := m.Push()
	if err != nil {
		return nil, err
	}
	return provider.Send(ctx, notification)
}

// SendPushWith sends a push notification using a specific provider
func (m *Manager) SendPushWith(ctx context.Context, providerName string, notification *contracts.PushNotification) (*contracts.SendResult, error) {
	provider, err := m.PushProvider(providerName)
	if err != nil {
		return nil, err
	}
	return provider.Send(ctx, notification)
}

// RegisterPushProvider registers a custom push notification provider
func (m *Manager) RegisterPushProvider(name string, provider contracts.PushSender) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.pushProviders[name] = provider
}

// AvailablePushProviders returns the names of all registered push providers
func (m *Manager) AvailablePushProviders() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	names := make([]string, 0, len(m.pushProviders))
	for name := range m.pushProviders {
		names = append(names, name)
	}
	return names
}
