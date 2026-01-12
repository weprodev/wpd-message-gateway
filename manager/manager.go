package manager

import (
	"fmt"
	"sync"

	"github.com/weprodev/wpd-message-gateway/config"
	"github.com/weprodev/wpd-message-gateway/contracts"
	"github.com/weprodev/wpd-message-gateway/providers/email/mailgun"
)

// Manager handles provider registration and message dispatching.
// It is the central gateway for all message types (Email, SMS, Push, Chat).
type Manager struct {
	config         *config.Config
	emailProviders map[string]contracts.EmailSender
	smsProviders   map[string]contracts.SMSSender
	pushProviders  map[string]contracts.PushSender
	chatProviders  map[string]contracts.ChatSender
	mu             sync.RWMutex
}

// New creates a new Manager with the given configuration.
// It automatically initializes all configured providers.
func New(cfg *config.Config) (*Manager, error) {
	m := &Manager{
		config:         cfg,
		emailProviders: make(map[string]contracts.EmailSender),
		smsProviders:   make(map[string]contracts.SMSSender),
		pushProviders:  make(map[string]contracts.PushSender),
		chatProviders:  make(map[string]contracts.ChatSender),
	}

	if err := m.initializeProviders(); err != nil {
		return nil, err
	}

	return m, nil
}

// initializeProviders creates provider instances from configuration.
func (m *Manager) initializeProviders() error {
	// Initialize Email Providers
	for name, cfg := range m.config.EmailProviders {
		switch name {
		case "mailgun":
			provider, err := mailgun.New(cfg)
			if err != nil {
				return fmt.Errorf("failed to initialize %s: %w", name, err)
			}
			m.emailProviders[name] = provider
		}
	}

	// Future: Initialize SMS, Push, Chat providers here

	return nil
}

// Config returns the manager's configuration (read-only access)
func (m *Manager) Config() *config.Config {
	return m.config
}
