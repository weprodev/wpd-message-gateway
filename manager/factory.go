package manager

import (
	"fmt"

	"github.com/weprodev/wpd-message-gateway/config"
	"github.com/weprodev/wpd-message-gateway/contracts"
	"github.com/weprodev/wpd-message-gateway/providers/email/mailgun"
	"github.com/weprodev/wpd-message-gateway/providers/memory"
)

// ProviderFactory creates provider instances based on name and configuration.
// This interface allows for dependency inversion - Manager depends on Factory interface,
// not concrete provider implementations.
type ProviderFactory interface {
	CreateEmailProvider(name string, cfg config.EmailConfig, mailpitCfg config.MailpitConfig, memStore *memory.Provider) (contracts.EmailSender, error)
	CreateSMSProvider(name string, cfg config.SMSConfig, memStore *memory.Provider) (contracts.SMSSender, error)
	CreatePushProvider(name string, cfg config.PushConfig, memStore *memory.Provider) (contracts.PushSender, error)
	CreateChatProvider(name string, cfg config.ChatConfig, memStore *memory.Provider) (contracts.ChatSender, error)
}

// DefaultFactory is the default implementation of ProviderFactory.
// It creates providers based on their name using a switch statement.
// This can be extended with registration patterns in the future.
type DefaultFactory struct{}

// NewDefaultFactory creates a new DefaultFactory instance.
func NewDefaultFactory() *DefaultFactory {
	return &DefaultFactory{}
}

// CreateEmailProvider creates an email provider based on the provider name.
func (f *DefaultFactory) CreateEmailProvider(name string, cfg config.EmailConfig, mailpitCfg config.MailpitConfig, memStore *memory.Provider) (contracts.EmailSender, error) {
	switch name {
	case "memory":
		if memStore == nil {
			return nil, fmt.Errorf("memory store is required for memory provider")
		}
		return memStore.EmailProvider(mailpitCfg), nil

	case "mailgun":
		provider, err := mailgun.New(cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to create mailgun provider: %w", err)
		}
		return provider, nil

	default:
		return nil, fmt.Errorf("unknown email provider: %s", name)
	}
}

// CreateSMSProvider creates an SMS provider based on the provider name.
func (f *DefaultFactory) CreateSMSProvider(name string, cfg config.SMSConfig, memStore *memory.Provider) (contracts.SMSSender, error) {
	switch name {
	case "memory":
		if memStore == nil {
			return nil, fmt.Errorf("memory store is required for memory provider")
		}
		return memStore.SMSProvider(), nil

	default:
		return nil, fmt.Errorf("unknown SMS provider: %s", name)
	}
}

// CreatePushProvider creates a push notification provider based on the provider name.
func (f *DefaultFactory) CreatePushProvider(name string, cfg config.PushConfig, memStore *memory.Provider) (contracts.PushSender, error) {
	switch name {
	case "memory":
		if memStore == nil {
			return nil, fmt.Errorf("memory store is required for memory provider")
		}
		return memStore.PushProvider(), nil

	default:
		return nil, fmt.Errorf("unknown push provider: %s", name)
	}
}

// CreateChatProvider creates a chat provider based on the provider name.
func (f *DefaultFactory) CreateChatProvider(name string, cfg config.ChatConfig, memStore *memory.Provider) (contracts.ChatSender, error) {
	switch name {
	case "memory":
		if memStore == nil {
			return nil, fmt.Errorf("memory store is required for memory provider")
		}
		return memStore.ChatProvider(), nil

	default:
		return nil, fmt.Errorf("unknown chat provider: %s", name)
	}
}
