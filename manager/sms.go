package manager

import (
	"context"

	"github.com/weprodev/wpd-message-gateway/contracts"
	msgerrors "github.com/weprodev/wpd-message-gateway/errors"
)

// SMS returns the default SMS sender
func (m *Manager) SMS() (contracts.SMSSender, error) {
	if m.config.DefaultSMSProvider == "" {
		return nil, msgerrors.NewProviderNotFoundError("sms", "default (none configured)")
	}
	return m.SMSProvider(m.config.DefaultSMSProvider)
}

// SMSProvider returns a specific SMS provider by name
func (m *Manager) SMSProvider(name string) (contracts.SMSSender, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	provider, ok := m.smsProviders[name]
	if !ok {
		return nil, msgerrors.NewProviderNotFoundError("sms", name)
	}
	return provider, nil
}

// SendSMS sends an SMS using the default provider
func (m *Manager) SendSMS(ctx context.Context, sms *contracts.SMS) (*contracts.SendResult, error) {
	provider, err := m.SMS()
	if err != nil {
		return nil, err
	}
	return provider.Send(ctx, sms)
}

// SendSMSWith sends an SMS using a specific provider
func (m *Manager) SendSMSWith(ctx context.Context, providerName string, sms *contracts.SMS) (*contracts.SendResult, error) {
	provider, err := m.SMSProvider(providerName)
	if err != nil {
		return nil, err
	}
	return provider.Send(ctx, sms)
}

// RegisterSMSProvider registers a custom SMS provider
func (m *Manager) RegisterSMSProvider(name string, provider contracts.SMSSender) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.smsProviders[name] = provider
}

// AvailableSMSProviders returns the names of all registered SMS providers
func (m *Manager) AvailableSMSProviders() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	names := make([]string, 0, len(m.smsProviders))
	for name := range m.smsProviders {
		names = append(names, name)
	}
	return names
}
