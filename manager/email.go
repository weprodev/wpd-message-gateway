package manager

import (
	"context"

	"github.com/weprodev/wpd-message-gateway/contracts"
	msgerrors "github.com/weprodev/wpd-message-gateway/errors"
)

// Email returns the default email sender.
func (m *Manager) Email() (contracts.EmailSender, error) {
	providerName := m.config.DefaultEmailProvider()
	if providerName == "" {
		return nil, msgerrors.NewProviderNotFoundError("email", "default (none configured)")
	}
	return m.EmailProvider(providerName)
}

// EmailProvider returns a specific email provider by name.
// Providers are created lazily if they don't exist yet.
func (m *Manager) EmailProvider(name string) (contracts.EmailSender, error) {
	// Try to get from registry first
	provider, ok := m.registry.GetEmailProvider(name)
	if ok {
		return provider, nil
	}

	// Provider doesn't exist, try to create it lazily
	memStore := m.registry.GetMemoryStore()
	if err := m.ensureEmailProvider(name, memStore); err != nil {
		return nil, msgerrors.NewProviderNotFoundError("email", name)
	}

	// Get the newly created provider
	provider, ok = m.registry.GetEmailProvider(name)
	if !ok {
		return nil, msgerrors.NewProviderNotFoundError("email", name)
	}

	return provider, nil
}

// SendEmail sends an email using the default provider.
// Routing logic: checks config to determine which provider to use.
func (m *Manager) SendEmail(ctx context.Context, email *contracts.Email) (*contracts.SendResult, error) {
	provider, err := m.Email()
	if err != nil {
		return nil, err
	}
	return provider.Send(ctx, email)
}

// SendEmailWith sends an email using a specific provider.
func (m *Manager) SendEmailWith(ctx context.Context, providerName string, email *contracts.Email) (*contracts.SendResult, error) {
	provider, err := m.EmailProvider(providerName)
	if err != nil {
		return nil, err
	}
	return provider.Send(ctx, email)
}

// RegisterEmailProvider registers a custom email provider.
func (m *Manager) RegisterEmailProvider(name string, provider contracts.EmailSender) {
	m.registry.RegisterEmailProvider(name, provider)
}

// AvailableEmailProviders returns the names of all registered email providers.
func (m *Manager) AvailableEmailProviders() []string {
	// Note: This only returns providers that have been initialized.
	// To get all configured providers, check m.config.EmailProviders
	providers := make([]string, 0)

	// Check registry for initialized providers
	commonProviders := []string{"memory", "mailgun"}
	for _, name := range commonProviders {
		if _, ok := m.registry.GetEmailProvider(name); ok {
			providers = append(providers, name)
		}
	}

	// Also check configured providers
	for name := range m.config.EmailProviders {
		if _, ok := m.registry.GetEmailProvider(name); ok {
			// Avoid duplicates
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
