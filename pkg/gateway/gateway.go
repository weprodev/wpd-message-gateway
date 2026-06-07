// Package gateway provides the public embedded SDK for wpd-message-gateway.
//
// There are two ways to use this package:
//
//  1. Static config (no server, no DB) — call New() with provider credentials
//     declared in code. Providers must be imported to trigger self-registration:
//
//     import _ "github.com/weprodev/wpd-message-gateway/pkg/provider/mailgun"
package gateway

import (
	"context"
	"fmt"

	"github.com/weprodev/wpd-message-gateway/pkg/contracts"
)

// Gateway sends messages. Construct via New().
type Gateway struct {

	// Static config mode — set by New().
	emailSenders map[string]contracts.EmailSender
	smsSenders   map[string]contracts.SMSSender
	pushSenders  map[string]contracts.PushSender
	chatSenders  map[string]contracts.ChatSender
	defaultEmail string
	defaultSMS   string
	defaultPush  string
	defaultChat  string
}

// New constructs a Gateway from static config without a database or HTTP server.
// Providers must have been registered via init() before New() is called —
// import the provider packages with a blank import in your main package:
//
//	import _ "github.com/weprodev/wpd-message-gateway/pkg/provider/mailgun"
//	import _ "github.com/weprodev/wpd-message-gateway/pkg/provider/memory"
func New(cfg Config) (*Gateway, error) {
	g := &Gateway{
		emailSenders: make(map[string]contracts.EmailSender),
		smsSenders:   make(map[string]contracts.SMSSender),
		pushSenders:  make(map[string]contracts.PushSender),
		chatSenders:  make(map[string]contracts.ChatSender),
		defaultEmail: cfg.DefaultEmailProvider,
		defaultSMS:   cfg.DefaultSMSProvider,
		defaultPush:  cfg.DefaultPushProvider,
		defaultChat:  cfg.DefaultChatProvider,
	}
	if err := buildEmailSenders(g, cfg); err != nil {
		return nil, err
	}
	if err := buildSMSSenders(g, cfg); err != nil {
		return nil, err
	}
	if err := buildPushSenders(g, cfg); err != nil {
		return nil, err
	}
	if err := buildChatSenders(g, cfg); err != nil {
		return nil, err
	}
	return g, nil
}

// SendEmail sends an email using the configured default email provider.
func (g *Gateway) SendEmail(ctx context.Context, email contracts.Email) (*contracts.SendResult, error) {
	sender, err := g.emailSender()
	if err != nil {
		return nil, err
	}
	return sender.Send(ctx, email)
}

// SendSMS sends an SMS using the configured default SMS provider.
func (g *Gateway) SendSMS(ctx context.Context, sms contracts.SMS) (*contracts.SendResult, error) {
	sender, err := g.smsSender()
	if err != nil {
		return nil, err
	}
	return sender.Send(ctx, sms)
}

// SendPush sends a push notification using the configured default push provider.
func (g *Gateway) SendPush(ctx context.Context, push contracts.PushNotification) (*contracts.SendResult, error) {
	sender, err := g.pushSender()
	if err != nil {
		return nil, err
	}
	return sender.Send(ctx, push)
}

// SendChat sends a chat message using the configured default chat provider.
func (g *Gateway) SendChat(ctx context.Context, chat contracts.ChatMessage) (*contracts.SendResult, error) {
	sender, err := g.chatSender()
	if err != nil {
		return nil, err
	}
	return sender.Send(ctx, chat)
}

func (g *Gateway) emailSender() (contracts.EmailSender, error) {
	s, ok := g.emailSenders[g.defaultEmail]
	if !ok {
		return nil, fmt.Errorf("gateway: email provider %q not configured", g.defaultEmail)
	}
	return s, nil
}

func (g *Gateway) smsSender() (contracts.SMSSender, error) {
	s, ok := g.smsSenders[g.defaultSMS]
	if !ok {
		return nil, fmt.Errorf("gateway: SMS provider %q not configured", g.defaultSMS)
	}
	return s, nil
}

func (g *Gateway) pushSender() (contracts.PushSender, error) {
	s, ok := g.pushSenders[g.defaultPush]
	if !ok {
		return nil, fmt.Errorf("gateway: push provider %q not configured", g.defaultPush)
	}
	return s, nil
}

func (g *Gateway) chatSender() (contracts.ChatSender, error) {
	s, ok := g.chatSenders[g.defaultChat]
	if !ok {
		return nil, fmt.Errorf("gateway: chat provider %q not configured", g.defaultChat)
	}
	return s, nil
}
