package manager

import (
	"fmt"

	"github.com/weprodev/wpd-message-gateway/config"
	"github.com/weprodev/wpd-message-gateway/providers/memory"
)

// Manager handles provider registration and message dispatching.
// It is the central gateway for all message types (Email, SMS, Push, Chat).
type Manager struct {
	config   *config.Config
	registry *Registry
	factory  ProviderFactory
}

// New creates a new Manager with the given configuration.
// Providers are created lazily on-demand via the factory.
func New(cfg *config.Config) (*Manager, error) {
	factory := NewDefaultFactory()
	registry := NewRegistry(factory)

	m := &Manager{
		config:   cfg,
		registry: registry,
		factory:  factory,
	}

	// Pre-initialize providers that are configured as defaults
	// This ensures they're available immediately
	if err := m.initializeDefaultProviders(); err != nil {
		return nil, fmt.Errorf("failed to initialize default providers: %w", err)
	}

	return m, nil
}

// initializeDefaultProviders pre-initializes providers that are set as defaults.
// Other providers are created lazily when first accessed.
// Returns error if a default provider fails to initialize (invalid config, etc.).
// Unknown providers (not in factory) are skipped - they can be registered manually.
func (m *Manager) initializeDefaultProviders() error {
	memStore := m.registry.GetMemoryStore()

	// Initialize default email provider if configured
	if m.config.DefaultEmailProvider() != "" {
		err := m.ensureEmailProvider(m.config.DefaultEmailProvider(), memStore)
		// Only fail if it's a known provider with invalid config
		// Unknown providers (not in factory) can be registered manually
		if err != nil && !isUnknownProviderError(err) {
			return fmt.Errorf("failed to initialize default email provider %s: %w", m.config.DefaultEmailProvider(), err)
		}
	}

	// Initialize default SMS provider if configured
	if m.config.DefaultSMSProvider() != "" {
		err := m.ensureSMSProvider(m.config.DefaultSMSProvider(), memStore)
		if err != nil && !isUnknownProviderError(err) {
			return fmt.Errorf("failed to initialize default SMS provider %s: %w", m.config.DefaultSMSProvider(), err)
		}
	}

	// Initialize default Push provider if configured
	if m.config.DefaultPushProvider() != "" {
		err := m.ensurePushProvider(m.config.DefaultPushProvider(), memStore)
		if err != nil && !isUnknownProviderError(err) {
			return fmt.Errorf("failed to initialize default push provider %s: %w", m.config.DefaultPushProvider(), err)
		}
	}

	// Initialize default Chat provider if configured
	if m.config.DefaultChatProvider() != "" {
		err := m.ensureChatProvider(m.config.DefaultChatProvider(), memStore)
		if err != nil && !isUnknownProviderError(err) {
			return fmt.Errorf("failed to initialize default chat provider %s: %w", m.config.DefaultChatProvider(), err)
		}
	}

	return nil
}

// isUnknownProviderError checks if an error indicates an unknown provider.
// Unknown providers can be registered manually, so we don't fail initialization.
func isUnknownProviderError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	// Check for errors indicating unknown/missing provider config
	return contains(errStr, "no configuration found") || contains(errStr, "unknown")
}

// contains checks if a string contains a substring.
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) && findSubstring(s, substr)))
}

// findSubstring finds a substring within a string.
func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// ensureEmailProvider ensures an email provider exists, creating it if necessary.
// Returns error if provider creation fails, nil if successful or provider already exists.
func (m *Manager) ensureEmailProvider(name string, memStore *memory.Provider) error {
	// Check if already registered
	if _, ok := m.registry.GetEmailProvider(name); ok {
		return nil
	}

	// Get config for this provider
	cfg, exists := m.config.EmailProviders[name]
	if !exists {
		// For memory provider, create empty config
		if name == "memory" {
			cfg = config.EmailConfig{}
		} else {
			// Unknown provider or missing config - skip (can be registered manually)
			return fmt.Errorf("no configuration found for email provider: %s", name)
		}
	}

	// Create provider via factory
	provider, err := m.factory.CreateEmailProvider(name, cfg, m.config.Mailpit, memStore)
	if err != nil {
		return err
	}

	// Register provider
	m.registry.RegisterEmailProvider(name, provider)
	return nil
}

// ensureSMSProvider ensures an SMS provider exists, creating it if necessary.
func (m *Manager) ensureSMSProvider(name string, memStore *memory.Provider) error {
	// Check if already registered
	if _, ok := m.registry.GetSMSProvider(name); ok {
		return nil
	}

	// Get config for this provider
	cfg, exists := m.config.SMSProviders[name]
	if !exists {
		// For memory provider, create empty config
		if name == "memory" {
			cfg = config.SMSConfig{}
		} else {
			return fmt.Errorf("no configuration found for SMS provider: %s", name)
		}
	}

	// Create provider via factory
	provider, err := m.factory.CreateSMSProvider(name, cfg, memStore)
	if err != nil {
		return err
	}

	// Register provider
	m.registry.RegisterSMSProvider(name, provider)
	return nil
}

// ensurePushProvider ensures a push provider exists, creating it if necessary.
func (m *Manager) ensurePushProvider(name string, memStore *memory.Provider) error {
	// Check if already registered
	if _, ok := m.registry.GetPushProvider(name); ok {
		return nil
	}

	// Get config for this provider
	cfg, exists := m.config.PushProviders[name]
	if !exists {
		// For memory provider, create empty config
		if name == "memory" {
			cfg = config.PushConfig{}
		} else {
			return fmt.Errorf("no configuration found for push provider: %s", name)
		}
	}

	// Create provider via factory
	provider, err := m.factory.CreatePushProvider(name, cfg, memStore)
	if err != nil {
		return err
	}

	// Register provider
	m.registry.RegisterPushProvider(name, provider)
	return nil
}

// ensureChatProvider ensures a chat provider exists, creating it if necessary.
func (m *Manager) ensureChatProvider(name string, memStore *memory.Provider) error {
	// Check if already registered
	if _, ok := m.registry.GetChatProvider(name); ok {
		return nil
	}

	// Get config for this provider
	cfg, exists := m.config.ChatProviders[name]
	if !exists {
		// For memory provider, create empty config
		if name == "memory" {
			cfg = config.ChatConfig{}
		} else {
			return fmt.Errorf("no configuration found for chat provider: %s", name)
		}
	}

	// Create provider via factory
	provider, err := m.factory.CreateChatProvider(name, cfg, memStore)
	if err != nil {
		return err
	}

	// Register provider
	m.registry.RegisterChatProvider(name, provider)
	return nil
}

// GetMemoryStore returns the shared memory store instance.
// This is used by DevBox UI to access stored messages.
func (m *Manager) GetMemoryStore() *memory.Provider {
	return m.registry.GetMemoryStore()
}

// Config returns the manager's configuration (read-only access)
func (m *Manager) Config() *config.Config {
	return m.config
}
