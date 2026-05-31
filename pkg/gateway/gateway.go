// Package gateway provides the public embedded SDK for wpd-message-gateway.
//
// There are two ways to use this package:
//
//  1. Static config (no server, no DB) — call New() with provider credentials
//     declared in code. Providers must be imported to trigger self-registration:
//
//     import _ "github.com/weprodev/wpd-message-gateway/internal/infrastructure/provider/mailgun"
//
//  2. DB-backed (server mode) — call NewWithService() with an already-wired
//     GatewayService. This is used internally by the HTTP server.
package gateway

import (
	"context"
	"fmt"

	"github.com/weprodev/wpd-message-gateway/internal/core/port"
	"github.com/weprodev/wpd-message-gateway/internal/core/service"
	"github.com/weprodev/wpd-message-gateway/pkg/contracts"
)

// Gateway sends messages. Construct via New() or NewWithService().
type Gateway struct {
	// DB-backed mode — set by NewWithService.
	svc         *service.GatewayService
	workspaceID string

	// Static config mode — set by New().
	emailSenders map[string]port.EmailSender
	smsSenders   map[string]port.SMSSender
	pushSenders  map[string]port.PushSender
	chatSenders  map[string]port.ChatSender
	otpSenders   map[string]port.OTPSender
	defaultEmail string
	defaultSMS   string
	defaultPush  string
	defaultChat  string
	defaultOTP   string
}

// NewWithService creates a Gateway scoped to a workspace using a wired
// GatewayService (database-backed integrations). Used internally by the server.
func NewWithService(svc *service.GatewayService, workspaceID string) *Gateway {
	return &Gateway{svc: svc, workspaceID: workspaceID}
}

// New constructs a Gateway from static config without a database or HTTP server.
// Providers must have been registered via init() before New() is called —
// import the provider packages with a blank import in your main package:
//
//	import _ "github.com/weprodev/wpd-message-gateway/internal/infrastructure/provider/mailgun"
//	import _ "github.com/weprodev/wpd-message-gateway/internal/infrastructure/provider/memory"
func New(cfg Config) (*Gateway, error) {
	g := &Gateway{
		emailSenders: make(map[string]port.EmailSender),
		smsSenders:   make(map[string]port.SMSSender),
		pushSenders:  make(map[string]port.PushSender),
		chatSenders:  make(map[string]port.ChatSender),
		otpSenders:   make(map[string]port.OTPSender),
		defaultEmail: cfg.DefaultEmailProvider,
		defaultSMS:   cfg.DefaultSMSProvider,
		defaultPush:  cfg.DefaultPushProvider,
		defaultChat:  cfg.DefaultChatProvider,
		defaultOTP:   cfg.DefaultOTPProvider,
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
	if err := buildOTPSenders(g, cfg); err != nil {
		return nil, err
	}
	return g, nil
}

// SendEmail sends an email. In static mode, uses the default email provider.
// In server mode, delegates to the workspace's active integration.
func (g *Gateway) SendEmail(ctx context.Context, email *contracts.Email) (*contracts.SendResult, error) {
	if g.svc != nil {
		return g.svc.SendEmail(ctx, g.workspaceID, email)
	}
	sender, err := g.emailSender()
	if err != nil {
		return nil, err
	}
	return sender.Send(ctx, email)
}

// SendSMS sends an SMS. In static mode, uses the default SMS provider.
func (g *Gateway) SendSMS(ctx context.Context, sms *contracts.SMS) (*contracts.SendResult, error) {
	if g.svc != nil {
		return g.svc.SendSMS(ctx, g.workspaceID, sms)
	}
	sender, err := g.smsSender()
	if err != nil {
		return nil, err
	}
	return sender.Send(ctx, sms)
}

// SendPush sends a push notification. In static mode, uses the default push provider.
func (g *Gateway) SendPush(ctx context.Context, push *contracts.PushNotification) (*contracts.SendResult, error) {
	if g.svc != nil {
		return g.svc.SendPush(ctx, g.workspaceID, push)
	}
	sender, err := g.pushSender()
	if err != nil {
		return nil, err
	}
	return sender.Send(ctx, push)
}

// SendChat sends a chat message. In static mode, uses the default chat provider.
func (g *Gateway) SendChat(ctx context.Context, chat *contracts.ChatMessage) (*contracts.SendResult, error) {
	if g.svc != nil {
		return g.svc.SendChat(ctx, g.workspaceID, chat)
	}
	sender, err := g.chatSender()
	if err != nil {
		return nil, err
	}
	return sender.Send(ctx, chat)
}

func (g *Gateway) SendOTP(ctx context.Context, otp *contracts.OTP) (*contracts.SendResult, error) {
	if g.svc != nil {
		return g.svc.SendOTP(ctx, g.workspaceID, otp)
	}
	sender, err := g.otpSender()
	if err != nil {
		return nil, err
	}
	return sender.Send(ctx, otp)
}

// RevokeOTP revokes a previously sent OTP message by request_id.
// In server mode, delegates to the workspace's active OTP integration.
func (g *Gateway) RevokeOTP(ctx context.Context, requestID string) (*contracts.SendResult, error) {
	if g.svc != nil {
		return g.svc.RevokeOTP(ctx, g.workspaceID, requestID)
	}
	sender, err := g.otpSender()
	if err != nil {
		return nil, err
	}
	revoker, ok := sender.(port.OTPRevoker)
	if !ok {
		return nil, fmt.Errorf("OTP provider %q does not support revocation", g.defaultOTP)
	}
	return revoker.Revoke(ctx, requestID)
}

func (g *Gateway) emailSender() (port.EmailSender, error) {
	s, ok := g.emailSenders[g.defaultEmail]
	if !ok {
		return nil, fmt.Errorf("gateway: email provider %q not configured", g.defaultEmail)
	}
	return s, nil
}

func (g *Gateway) smsSender() (port.SMSSender, error) {
	s, ok := g.smsSenders[g.defaultSMS]
	if !ok {
		return nil, fmt.Errorf("gateway: SMS provider %q not configured", g.defaultSMS)
	}
	return s, nil
}

func (g *Gateway) pushSender() (port.PushSender, error) {
	s, ok := g.pushSenders[g.defaultPush]
	if !ok {
		return nil, fmt.Errorf("gateway: push provider %q not configured", g.defaultPush)
	}
	return s, nil
}

func (g *Gateway) chatSender() (port.ChatSender, error) {
	s, ok := g.chatSenders[g.defaultChat]
	if !ok {
		return nil, fmt.Errorf("gateway: chat provider %q not configured", g.defaultChat)
	}
	return s, nil
}

func (g *Gateway) otpSender() (port.OTPSender, error) {
	s, ok := g.otpSenders[g.defaultOTP]
	if !ok {
		return nil, fmt.Errorf("gateway: OTP provider %q not configured", g.defaultOTP)
	}
	return s, nil
}
