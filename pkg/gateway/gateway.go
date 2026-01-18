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

	"github.com/weprodev/wpd-message-gateway/internal/core/port"
	"github.com/weprodev/wpd-message-gateway/internal/core/service"
	"github.com/weprodev/wpd-message-gateway/internal/infrastructure/provider/mailgun"
	"github.com/weprodev/wpd-message-gateway/internal/infrastructure/provider/memory"
	"github.com/weprodev/wpd-message-gateway/pkg/contracts"
)

// Config holds the gateway configuration.
type Config struct {
	DefaultEmailProvider string
	DefaultSMSProvider   string
	DefaultPushProvider  string
	DefaultChatProvider  string

	Mailgun MailgunConfig
	Memory  MemoryConfig
}

// MailgunConfig holds Mailgun provider configuration.
type MailgunConfig struct {
	APIKey    string
	Domain    string
	BaseURL   string
	FromEmail string
	FromName  string
}

// MemoryConfig holds memory provider configuration.
type MemoryConfig struct {
	MailpitEnabled bool
}

// Gateway is the main entry point for sending messages.
type Gateway struct {
	service     *service.GatewayService
	memoryStore *memory.Store
}

// configAdapter adapts Gateway config to service.GatewayConfig interface.
type configAdapter struct {
	cfg Config
}

func (c *configAdapter) DefaultEmailProvider() string { return c.cfg.DefaultEmailProvider }
func (c *configAdapter) DefaultSMSProvider() string   { return c.cfg.DefaultSMSProvider }
func (c *configAdapter) DefaultPushProvider() string  { return c.cfg.DefaultPushProvider }
func (c *configAdapter) DefaultChatProvider() string  { return c.cfg.DefaultChatProvider }

// New creates a new Gateway instance.
func New(cfg Config) (*Gateway, error) {
	registry := service.NewRegistry()
	memoryStore := memory.NewStore()

	gw := &Gateway{
		memoryStore: memoryStore,
	}

	if err := gw.initializeProviders(cfg, registry); err != nil {
		return nil, err
	}

	gw.service = service.NewGatewayService(&configAdapter{cfg}, registry)

	return gw, nil
}

func (g *Gateway) initializeProviders(cfg Config, registry *service.Registry) error {
	if cfg.DefaultEmailProvider != "" {
		provider, err := g.createEmailProvider(cfg, cfg.DefaultEmailProvider)
		if err != nil {
			return fmt.Errorf("failed to create email provider: %w", err)
		}
		registry.RegisterEmailProvider(cfg.DefaultEmailProvider, provider)
	}

	if cfg.DefaultSMSProvider != "" {
		provider, err := g.createSMSProvider(cfg.DefaultSMSProvider)
		if err != nil {
			return fmt.Errorf("failed to create SMS provider: %w", err)
		}
		registry.RegisterSMSProvider(cfg.DefaultSMSProvider, provider)
	}

	if cfg.DefaultPushProvider != "" {
		provider, err := g.createPushProvider(cfg.DefaultPushProvider)
		if err != nil {
			return fmt.Errorf("failed to create push provider: %w", err)
		}
		registry.RegisterPushProvider(cfg.DefaultPushProvider, provider)
	}

	if cfg.DefaultChatProvider != "" {
		provider, err := g.createChatProvider(cfg.DefaultChatProvider)
		if err != nil {
			return fmt.Errorf("failed to create chat provider: %w", err)
		}
		registry.RegisterChatProvider(cfg.DefaultChatProvider, provider)
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

func (g *Gateway) createEmailProvider(cfg Config, name string) (port.EmailSender, error) {
	switch name {
	case "memory":
		mailpitCfg := memory.MailpitConfig{Enabled: cfg.Memory.MailpitEnabled}
		return memory.NewEmailProvider(g.memoryStore, mailpitCfg), nil
	case "mailgun":
		return mailgun.New(mailgun.Config{
			APIKey:    cfg.Mailgun.APIKey,
			Domain:    cfg.Mailgun.Domain,
			BaseURL:   cfg.Mailgun.BaseURL,
			FromEmail: cfg.Mailgun.FromEmail,
			FromName:  cfg.Mailgun.FromName,
		})
	default:
		return nil, fmt.Errorf("unknown email provider: %s", name)
	}
}

func (g *Gateway) createSMSProvider(name string) (port.SMSSender, error) {
	switch name {
	case "memory":
		return memory.NewSMSProvider(g.memoryStore), nil
	default:
		return nil, fmt.Errorf("unknown SMS provider: %s", name)
	}
}

func (g *Gateway) createPushProvider(name string) (port.PushSender, error) {
	switch name {
	case "memory":
		return memory.NewPushProvider(g.memoryStore), nil
	default:
		return nil, fmt.Errorf("unknown push provider: %s", name)
	}
}

func (g *Gateway) createChatProvider(name string) (port.ChatSender, error) {
	switch name {
	case "memory":
		return memory.NewChatProvider(g.memoryStore), nil
	default:
		return nil, fmt.Errorf("unknown chat provider: %s", name)
	}
}
