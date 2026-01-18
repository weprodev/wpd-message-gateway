// Package gateway provides an embedded SDK for the message gateway.
//
// Usage:
//
//	gw, err := gateway.New(gateway.Config{
//	    DefaultEmailProvider: "memory",
//	})
//
//	result, err := gw.SendEmail(ctx, &contracts.Email{
//	    To:      []string{"user@example.com"},
//	    Subject: "Hello",
//	    HTML:    "<p>World</p>",
//	})
package gateway

import (
	"context"
	"fmt"

	"github.com/weprodev/wpd-message-gateway/internal/app/registry"
	"github.com/weprodev/wpd-message-gateway/internal/core/port"
	"github.com/weprodev/wpd-message-gateway/internal/core/service"
	"github.com/weprodev/wpd-message-gateway/pkg/contracts"
)

func (c *configAdapter) DefaultEmailProvider() string { return c.cfg.DefaultEmailProvider }
func (c *configAdapter) DefaultSMSProvider() string   { return c.cfg.DefaultSMSProvider }
func (c *configAdapter) DefaultPushProvider() string  { return c.cfg.DefaultPushProvider }
func (c *configAdapter) DefaultChatProvider() string  { return c.cfg.DefaultChatProvider }

// New creates a new Gateway instance.
func New(cfg Config) (*Gateway, error) {
	serviceRegistry := service.NewRegistry()

	gw := &Gateway{cfg: cfg}

	if err := gw.initializeProviders(serviceRegistry); err != nil {
		return nil, err
	}

	gw.service = service.NewGatewayService(&configAdapter{cfg}, serviceRegistry)

	return gw, nil
}

func (g *Gateway) initializeProviders(serviceRegistry *service.Registry) error {
	if name := g.cfg.DefaultEmailProvider; name != "" {
		provider, err := g.createEmailProvider(name)
		if err != nil {
			return fmt.Errorf("failed to create email provider %s: %w", name, err)
		}
		serviceRegistry.RegisterEmailProvider(name, provider)
	}

	if name := g.cfg.DefaultSMSProvider; name != "" {
		provider, err := g.createSMSProvider(name)
		if err != nil {
			return fmt.Errorf("failed to create SMS provider %s: %w", name, err)
		}
		serviceRegistry.RegisterSMSProvider(name, provider)
	}

	if name := g.cfg.DefaultPushProvider; name != "" {
		provider, err := g.createPushProvider(name)
		if err != nil {
			return fmt.Errorf("failed to create push provider %s: %w", name, err)
		}
		serviceRegistry.RegisterPushProvider(name, provider)
	}

	if name := g.cfg.DefaultChatProvider; name != "" {
		provider, err := g.createChatProvider(name)
		if err != nil {
			return fmt.Errorf("failed to create chat provider %s: %w", name, err)
		}
		serviceRegistry.RegisterChatProvider(name, provider)
	}

	return nil
}

// SendEmail sends an email using the default provider.
func (g *Gateway) SendEmail(ctx context.Context, email *contracts.Email) (*contracts.SendResult, error) {
	return g.service.SendEmail(ctx, email)
}

// SendEmailWith sends an email using a specific provider.
func (g *Gateway) SendEmailWith(ctx context.Context, provider string, email *contracts.Email) (*contracts.SendResult, error) {
	return g.service.SendEmailWith(ctx, provider, email)
}

// SendSMS sends an SMS using the default provider.
func (g *Gateway) SendSMS(ctx context.Context, sms *contracts.SMS) (*contracts.SendResult, error) {
	return g.service.SendSMS(ctx, sms)
}

// SendSMSWith sends an SMS using a specific provider.
func (g *Gateway) SendSMSWith(ctx context.Context, provider string, sms *contracts.SMS) (*contracts.SendResult, error) {
	return g.service.SendSMSWith(ctx, provider, sms)
}

// SendPush sends a push notification using the default provider.
func (g *Gateway) SendPush(ctx context.Context, push *contracts.PushNotification) (*contracts.SendResult, error) {
	return g.service.SendPush(ctx, push)
}

// SendPushWith sends a push notification using a specific provider.
func (g *Gateway) SendPushWith(ctx context.Context, provider string, push *contracts.PushNotification) (*contracts.SendResult, error) {
	return g.service.SendPushWith(ctx, provider, push)
}

// SendChat sends a chat message using the default provider.
func (g *Gateway) SendChat(ctx context.Context, chat *contracts.ChatMessage) (*contracts.SendResult, error) {
	return g.service.SendChat(ctx, chat)
}

// SendChatWith sends a chat message using a specific provider.
func (g *Gateway) SendChatWith(ctx context.Context, provider string, chat *contracts.ChatMessage) (*contracts.SendResult, error) {
	return g.service.SendChatWith(ctx, provider, chat)
}

func (g *Gateway) createEmailProvider(name string) (port.EmailSender, error) {
	factory, err := registry.GetEmailFactory(name)
	if err != nil {
		return nil, err
	}

	cfg := g.cfg.EmailProviders[name]
	mailpit := registry.MailpitConfig{Enabled: g.cfg.MailpitEnabled}
	return factory(cfg, mailpit)
}

func (g *Gateway) createSMSProvider(name string) (port.SMSSender, error) {
	factory, err := registry.GetSMSFactory(name)
	if err != nil {
		return nil, err
	}

	return factory(g.cfg.SMSProviders[name])
}

func (g *Gateway) createPushProvider(name string) (port.PushSender, error) {
	factory, err := registry.GetPushFactory(name)
	if err != nil {
		return nil, err
	}

	return factory(g.cfg.PushProviders[name])
}

func (g *Gateway) createChatProvider(name string) (port.ChatSender, error) {
	factory, err := registry.GetChatFactory(name)
	if err != nil {
		return nil, err
	}

	return factory(g.cfg.ChatProviders[name])
}
