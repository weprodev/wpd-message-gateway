package manager

import (
	"context"

	"github.com/weprodev/wpd-message-gateway/contracts"
	msgerrors "github.com/weprodev/wpd-message-gateway/errors"
)

// Push returns the default push notification sender.
func (m *Manager) Push() (contracts.PushSender, error) {
	providerName := m.config.DefaultPushProvider()
	if providerName == "" {
		return nil, msgerrors.NewProviderNotFoundError("push", "default (none configured)")
	}
	return m.PushProvider(providerName)
}

// PushProvider returns a specific push notification provider by name.
// PushProvider returns the configured default push provider.
func (m *Manager) PushProvider(providerName string) (contracts.PushSender, error) {
	if providerName == "" {
		providerName = m.config.DefaultPushProvider()
	}

	// Try to get from registry first
	provider, ok := m.registry.GetPushProvider(providerName)
	if ok {
		return provider, nil
	}

	// Provider doesn't exist, try to create it lazily
	memStore := m.registry.GetMemoryStore()
	if err := m.ensurePushProvider(providerName, memStore); err != nil {
		return nil, msgerrors.NewProviderNotFoundError("push", providerName)
	}

	// Get the newly created provider
	provider, ok = m.registry.GetPushProvider(providerName)
	if !ok {
		return nil, msgerrors.NewProviderNotFoundError("push", providerName)
	}

	return provider, nil
}

// SendPush sends a push notification using the default provider.
func (m *Manager) SendPush(ctx context.Context, notification *contracts.PushNotification) (*contracts.SendResult, error) {
	provider, err := m.Push()
	if err != nil {
		return nil, err
	}
	return provider.Send(ctx, notification)
}

// SendPushWith sends a push notification using a specific provider.
func (m *Manager) SendPushWith(ctx context.Context, providerName string, notification *contracts.PushNotification) (*contracts.SendResult, error) {
	provider, err := m.PushProvider(providerName)
	if err != nil {
		return nil, err
	}
	return provider.Send(ctx, notification)
}

// RegisterPushProvider registers a custom push notification provider.
func (m *Manager) RegisterPushProvider(name string, provider contracts.PushSender) {
	m.registry.RegisterPushProvider(name, provider)
}

// AvailablePushProviders returns the names of all registered push providers.
func (m *Manager) AvailablePushProviders() []string {
	providers := make([]string, 0)

	// Check common providers
	commonProviders := []string{"memory"}
	for _, name := range commonProviders {
		if _, ok := m.registry.GetPushProvider(name); ok {
			providers = append(providers, name)
		}
	}

	// Check configured providers
	for name := range m.config.PushProviders {
		if _, ok := m.registry.GetPushProvider(name); ok {
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
