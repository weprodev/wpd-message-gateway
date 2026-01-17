package manager

import (
	"context"

	"github.com/weprodev/wpd-message-gateway/contracts"
	msgerrors "github.com/weprodev/wpd-message-gateway/errors"
)

// SMS returns the default SMS sender.
func (m *Manager) SMS() (contracts.SMSSender, error) {
	providerName := m.config.DefaultSMSProvider()
	if providerName == "" {
		return nil, msgerrors.NewProviderNotFoundError("sms", "default (none configured)")
	}
	return m.SMSProvider(providerName)
}

// SMSProvider returns a specific SMS provider by name.
// Providers are created lazily if they don't exist yet.
func (m *Manager) SMSProvider(name string) (contracts.SMSSender, error) {
	// Try to get from registry first
	provider, ok := m.registry.GetSMSProvider(name)
	if ok {
		return provider, nil
	}

	// Provider doesn't exist, try to create it lazily
	memStore := m.registry.GetMemoryStore()
	if err := m.ensureSMSProvider(name, memStore); err != nil {
		return nil, msgerrors.NewProviderNotFoundError("sms", name)
	}

	// Get the newly created provider
	provider, ok = m.registry.GetSMSProvider(name)
	if !ok {
		return nil, msgerrors.NewProviderNotFoundError("sms", name)
	}

	return provider, nil
}

// SendSMS sends an SMS using the default provider.
func (m *Manager) SendSMS(ctx context.Context, sms *contracts.SMS) (*contracts.SendResult, error) {
	provider, err := m.SMS()
	if err != nil {
		return nil, err
	}
	return provider.Send(ctx, sms)
}

// SendSMSWith sends an SMS using a specific provider.
func (m *Manager) SendSMSWith(ctx context.Context, providerName string, sms *contracts.SMS) (*contracts.SendResult, error) {
	provider, err := m.SMSProvider(providerName)
	if err != nil {
		return nil, err
	}
	return provider.Send(ctx, sms)
}

// RegisterSMSProvider registers a custom SMS provider.
func (m *Manager) RegisterSMSProvider(name string, provider contracts.SMSSender) {
	m.registry.RegisterSMSProvider(name, provider)
}

// AvailableSMSProviders returns the names of all registered SMS providers.
func (m *Manager) AvailableSMSProviders() []string {
	providers := make([]string, 0)

	// Check common providers
	commonProviders := []string{"memory"}
	for _, name := range commonProviders {
		if _, ok := m.registry.GetSMSProvider(name); ok {
			providers = append(providers, name)
		}
	}

	// Check configured providers
	for name := range m.config.SMSProviders {
		if _, ok := m.registry.GetSMSProvider(name); ok {
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
