package memory

import (
	"github.com/weprodev/wpd-message-gateway/pkg/contracts"
	"github.com/weprodev/wpd-message-gateway/pkg/registry"
)

func init() {
	registry.RegisterEmailProvider("memory", func(cfg registry.EmailConfig) (contracts.EmailSender, error) {
		return NewEmailProvider(GetStore()), nil
	})

	registry.RegisterSMSProvider("memory", func(cfg registry.SMSConfig) (contracts.SMSSender, error) {
		return NewSMSProvider(GetStore()), nil
	})

	registry.RegisterPushProvider("memory", func(cfg registry.PushConfig) (contracts.PushSender, error) {
		return NewPushProvider(GetStore()), nil
	})

	registry.RegisterChatProvider("memory", func(cfg registry.ChatConfig) (contracts.ChatSender, error) {
		return NewChatProvider(GetStore()), nil
	})
}
