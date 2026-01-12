package manager

import (
	"context"

	"github.com/weprodev/wpd-message-gateway/contracts"
	msgerrors "github.com/weprodev/wpd-message-gateway/errors"
)

// Email returns the default email sender
func (m *Manager) Email() (contracts.EmailSender, error) {
	if m.config.DefaultEmailProvider == "" {
		return nil, msgerrors.NewProviderNotFoundError("email", "default (none configured)")
	}
	return m.EmailProvider(m.config.DefaultEmailProvider)
}

// EmailProvider returns a specific email provider by name
func (m *Manager) EmailProvider(name string) (contracts.EmailSender, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	provider, ok := m.emailProviders[name]
	if !ok {
		return nil, msgerrors.NewProviderNotFoundError("email", name)
	}
	return provider, nil
}

// SendEmail sends an email using the default provider
func (m *Manager) SendEmail(ctx context.Context, email *contracts.Email) (*contracts.SendResult, error) {
	provider, err := m.Email()
	if err != nil {
		return nil, err
	}
	return provider.Send(ctx, email)
}

// SendEmailWith sends an email using a specific provider
func (m *Manager) SendEmailWith(ctx context.Context, providerName string, email *contracts.Email) (*contracts.SendResult, error) {
	provider, err := m.EmailProvider(providerName)
	if err != nil {
		return nil, err
	}
	return provider.Send(ctx, email)
}

// RegisterEmailProvider registers a custom email provider
func (m *Manager) RegisterEmailProvider(name string, provider contracts.EmailSender) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.emailProviders[name] = provider
}

// AvailableEmailProviders returns the names of all registered email providers
func (m *Manager) AvailableEmailProviders() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	names := make([]string, 0, len(m.emailProviders))
	for name := range m.emailProviders {
		names = append(names, name)
	}
	return names
}
