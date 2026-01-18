package service

import (
	"context"
	"fmt"

	"github.com/weprodev/wpd-message-gateway/internal/core/port"
	"github.com/weprodev/wpd-message-gateway/pkg/contracts"
)

// GatewayService handles provider registration and message dispatching.
type GatewayService struct {
	config   GatewayConfig
	registry *Registry
}

// GatewayConfig holds the configuration needed by the service.
type GatewayConfig interface {
	DefaultEmailProvider() string
	DefaultSMSProvider() string
	DefaultPushProvider() string
	DefaultChatProvider() string
}

// NewGatewayService creates a new GatewayService.
func NewGatewayService(cfg GatewayConfig, registry *Registry) *GatewayService {
	return &GatewayService{
		config:   cfg,
		registry: registry,
	}
}

// SendEmail sends an email using the default provider.
func (s *GatewayService) SendEmail(ctx context.Context, email *contracts.Email) (*contracts.SendResult, error) {
	provider, err := s.Email()
	if err != nil {
		return nil, err
	}
	return provider.Send(ctx, email)
}

// SendEmailWith sends an email using a specific provider.
func (s *GatewayService) SendEmailWith(ctx context.Context, providerName string, email *contracts.Email) (*contracts.SendResult, error) {
	provider, err := s.EmailProvider(providerName)
	if err != nil {
		return nil, err
	}
	return provider.Send(ctx, email)
}

// Email returns the default email provider.
func (s *GatewayService) Email() (port.EmailSender, error) {
	providerName := s.config.DefaultEmailProvider()
	if providerName == "" {
		return nil, NewProviderNotFoundError("email", "default (none configured)")
	}
	return s.EmailProvider(providerName)
}

// EmailProvider returns a specific email provider by name.
func (s *GatewayService) EmailProvider(name string) (port.EmailSender, error) {
	provider, ok := s.registry.GetEmailProvider(name)
	if !ok {
		return nil, NewProviderNotFoundError("email", name)
	}
	return provider, nil
}

// RegisterEmailProvider registers a custom email provider.
func (s *GatewayService) RegisterEmailProvider(name string, provider port.EmailSender) {
	s.registry.RegisterEmailProvider(name, provider)
}

// SendSMS sends an SMS using the default provider.
func (s *GatewayService) SendSMS(ctx context.Context, sms *contracts.SMS) (*contracts.SendResult, error) {
	provider, err := s.SMS()
	if err != nil {
		return nil, err
	}
	return provider.Send(ctx, sms)
}

// SendSMSWith sends an SMS using a specific provider.
func (s *GatewayService) SendSMSWith(ctx context.Context, providerName string, sms *contracts.SMS) (*contracts.SendResult, error) {
	provider, err := s.SMSProvider(providerName)
	if err != nil {
		return nil, err
	}
	return provider.Send(ctx, sms)
}

// SMS returns the default SMS provider.
func (s *GatewayService) SMS() (port.SMSSender, error) {
	providerName := s.config.DefaultSMSProvider()
	if providerName == "" {
		return nil, NewProviderNotFoundError("sms", "default (none configured)")
	}
	return s.SMSProvider(providerName)
}

// SMSProvider returns a specific SMS provider by name.
func (s *GatewayService) SMSProvider(name string) (port.SMSSender, error) {
	provider, ok := s.registry.GetSMSProvider(name)
	if !ok {
		return nil, NewProviderNotFoundError("sms", name)
	}
	return provider, nil
}

// RegisterSMSProvider registers a custom SMS provider.
func (s *GatewayService) RegisterSMSProvider(name string, provider port.SMSSender) {
	s.registry.RegisterSMSProvider(name, provider)
}

// SendPush sends a push notification using the default provider.
func (s *GatewayService) SendPush(ctx context.Context, notification *contracts.PushNotification) (*contracts.SendResult, error) {
	provider, err := s.Push()
	if err != nil {
		return nil, err
	}
	return provider.Send(ctx, notification)
}

// SendPushWith sends a push notification using a specific provider.
func (s *GatewayService) SendPushWith(ctx context.Context, providerName string, notification *contracts.PushNotification) (*contracts.SendResult, error) {
	provider, err := s.PushProvider(providerName)
	if err != nil {
		return nil, err
	}
	return provider.Send(ctx, notification)
}

// Push returns the default push provider.
func (s *GatewayService) Push() (port.PushSender, error) {
	providerName := s.config.DefaultPushProvider()
	if providerName == "" {
		return nil, NewProviderNotFoundError("push", "default (none configured)")
	}
	return s.PushProvider(providerName)
}

// PushProvider returns a specific push provider by name.
func (s *GatewayService) PushProvider(name string) (port.PushSender, error) {
	provider, ok := s.registry.GetPushProvider(name)
	if !ok {
		return nil, NewProviderNotFoundError("push", name)
	}
	return provider, nil
}

// RegisterPushProvider registers a custom push provider.
func (s *GatewayService) RegisterPushProvider(name string, provider port.PushSender) {
	s.registry.RegisterPushProvider(name, provider)
}

// SendChat sends a chat message using the default provider.
func (s *GatewayService) SendChat(ctx context.Context, message *contracts.ChatMessage) (*contracts.SendResult, error) {
	provider, err := s.Chat()
	if err != nil {
		return nil, err
	}
	return provider.Send(ctx, message)
}

// SendChatWith sends a chat message using a specific provider.
func (s *GatewayService) SendChatWith(ctx context.Context, providerName string, message *contracts.ChatMessage) (*contracts.SendResult, error) {
	provider, err := s.ChatProvider(providerName)
	if err != nil {
		return nil, err
	}
	return provider.Send(ctx, message)
}

// Chat returns the default chat provider.
func (s *GatewayService) Chat() (port.ChatSender, error) {
	providerName := s.config.DefaultChatProvider()
	if providerName == "" {
		return nil, NewProviderNotFoundError("chat", "default (none configured)")
	}
	return s.ChatProvider(providerName)
}

// ChatProvider returns a specific chat provider by name.
func (s *GatewayService) ChatProvider(name string) (port.ChatSender, error) {
	provider, ok := s.registry.GetChatProvider(name)
	if !ok {
		return nil, NewProviderNotFoundError("chat", name)
	}
	return provider, nil
}

// RegisterChatProvider registers a custom chat provider.
func (s *GatewayService) RegisterChatProvider(name string, provider port.ChatSender) {
	s.registry.RegisterChatProvider(name, provider)
}

type ProviderNotFoundError struct {
	ProviderType string
	ProviderName string
}

func (e *ProviderNotFoundError) Error() string {
	return fmt.Sprintf("%s provider '%s' not found", e.ProviderType, e.ProviderName)
}

func NewProviderNotFoundError(providerType, providerName string) *ProviderNotFoundError {
	return &ProviderNotFoundError{
		ProviderType: providerType,
		ProviderName: providerName,
	}
}
