package memory

import (
	"fmt"

	"github.com/weprodev/wpd-message-gateway/internal/app/registry"
	"github.com/weprodev/wpd-message-gateway/internal/core/port"
)

func init() {
	registry.RegisterEmailProvider("memory", func(cfg registry.EmailConfig, store port.MessageStore, mailpit registry.MailpitConfig) (port.EmailSender, error) {
		memStore, ok := store.(*Store)
		if !ok {
			return nil, fmt.Errorf("memory provider requires *memory.Store, got %T", store)
		}
		mailpitCfg := MailpitConfig{Enabled: mailpit.Enabled}
		return NewEmailProvider(memStore, mailpitCfg), nil
	})

	registry.RegisterSMSProvider("memory", func(cfg registry.SMSConfig, store port.MessageStore) (port.SMSSender, error) {
		memStore, ok := store.(*Store)
		if !ok {
			return nil, fmt.Errorf("memory provider requires *memory.Store, got %T", store)
		}
		return NewSMSProvider(memStore), nil
	})

	registry.RegisterPushProvider("memory", func(cfg registry.PushConfig, store port.MessageStore) (port.PushSender, error) {
		memStore, ok := store.(*Store)
		if !ok {
			return nil, fmt.Errorf("memory provider requires *memory.Store, got %T", store)
		}
		return NewPushProvider(memStore), nil
	})

	registry.RegisterChatProvider("memory", func(cfg registry.ChatConfig, store port.MessageStore) (port.ChatSender, error) {
		memStore, ok := store.(*Store)
		if !ok {
			return nil, fmt.Errorf("memory provider requires *memory.Store, got %T", store)
		}
		return NewChatProvider(memStore), nil
	})
}
