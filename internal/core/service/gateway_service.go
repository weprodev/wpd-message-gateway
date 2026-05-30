// Package service contains application services that orchestrate domain logic.
// Services depend only on port interfaces — never on infrastructure implementations.
package service

import (
	"context"
	"fmt"

	"github.com/weprodev/wpd-message-gateway/internal/core/domain"
	"github.com/weprodev/wpd-message-gateway/internal/core/port"
	"github.com/weprodev/wpd-message-gateway/pkg/contracts"
)

const (
	channelEmail = "email"
	channelSMS   = "sms"
	channelPush  = "push"
	channelChat  = "chat"
	channelOTP   = "otp"

	// memoryProviderName is the sentinel name stored in the integrations table
	// when a workspace uses memory dispatch. Checked here to skip DB lookup.
	memoryProviderName = "memory"
)

// GatewayService orchestrates message dispatch for a single workspace.
//
// It resolves the workspace dispatch mode from WorkspaceSettingsRepository,
// then routes the message to the in-process InboxWriter, a real provider
// fetched via IntegrationRepository, or both — depending on the mode.
//
// Provider instances are cached per (integrationID, updatedAt) so HTTP clients
// are reused across requests while auto-invalidating on config changes.
type GatewayService struct {
	integrations port.IntegrationRepository
	templates    port.TemplateRepository
	settings     port.WorkspaceSettingsRepository
	inbox        port.InboxWriter

	emailCache *emailSenderCache
	smsCache   *smsSenderCache
	pushCache  *pushSenderCache
	chatCache  *chatSenderCache
	otpCache   *otpSenderCache
}

// NewGatewayService constructs a GatewayService.
//
// inbox is the in-process capture store (implements port.InboxWriter).
// Pass nil to disable memory capture — memory_only dispatch will then error.
func NewGatewayService(
	integrations port.IntegrationRepository,
	templates port.TemplateRepository,
	settings port.WorkspaceSettingsRepository,
	inbox port.InboxWriter,
) *GatewayService {
	return &GatewayService{
		integrations: integrations,
		templates:    templates,
		settings:     settings,
		inbox:        inbox,
		emailCache:   newProviderCache[port.EmailSender](),
		smsCache:     newProviderCache[port.SMSSender](),
		pushCache:    newProviderCache[port.PushSender](),
		chatCache:    newProviderCache[port.ChatSender](),
		otpCache:     newProviderCache[port.OTPSender](),
	}
}

// SendEmail dispatches an email for workspaceID according to its dispatch mode.
func (s *GatewayService) SendEmail(ctx context.Context, workspaceID string, email *contracts.Email) (*contracts.SendResult, error) {
	mode := s.resolveDispatchMode(ctx, workspaceID)
	return s.dispatch(ctx, workspaceID, mode, channelEmail,
		func() (*contracts.SendResult, error) { return s.writeEmailToInbox(ctx, workspaceID, email) },
		func(intg *domain.Integration) (*contracts.SendResult, error) {
			sender, err := resolveEmailSender(s.emailCache, intg)
			if err != nil {
				return nil, err
			}
			return sender.Send(ctx, email)
		},
	)
}

// SendSMS dispatches an SMS for workspaceID according to its dispatch mode.
func (s *GatewayService) SendSMS(ctx context.Context, workspaceID string, sms *contracts.SMS) (*contracts.SendResult, error) {
	mode := s.resolveDispatchMode(ctx, workspaceID)
	return s.dispatch(ctx, workspaceID, mode, channelSMS,
		func() (*contracts.SendResult, error) { return s.writeSMSToInbox(ctx, workspaceID, sms) },
		func(intg *domain.Integration) (*contracts.SendResult, error) {
			sender, err := resolveSMSSender(s.smsCache, intg)
			if err != nil {
				return nil, err
			}
			return sender.Send(ctx, sms)
		},
	)
}

// SendPush dispatches a push notification for workspaceID according to its dispatch mode.
func (s *GatewayService) SendPush(ctx context.Context, workspaceID string, push *contracts.PushNotification) (*contracts.SendResult, error) {
	mode := s.resolveDispatchMode(ctx, workspaceID)
	return s.dispatch(ctx, workspaceID, mode, channelPush,
		func() (*contracts.SendResult, error) { return s.writePushToInbox(ctx, workspaceID, push) },
		func(intg *domain.Integration) (*contracts.SendResult, error) {
			sender, err := resolvePushSender(s.pushCache, intg)
			if err != nil {
				return nil, err
			}
			return sender.Send(ctx, push)
		},
	)
}

// SendChat dispatches a chat message for workspaceID according to its dispatch mode.
func (s *GatewayService) SendChat(ctx context.Context, workspaceID string, chat *contracts.ChatMessage) (*contracts.SendResult, error) {
	mode := s.resolveDispatchMode(ctx, workspaceID)
	return s.dispatch(ctx, workspaceID, mode, channelChat,
		func() (*contracts.SendResult, error) { return s.writeChatToInbox(ctx, workspaceID, chat) },
		func(intg *domain.Integration) (*contracts.SendResult, error) {
			sender, err := resolveChatSender(s.chatCache, intg)
			if err != nil {
				return nil, err
			}
			return sender.Send(ctx, chat)
		},
	)
}

// SendOTP dispatches an OTP message for workspaceID according to its dispatch mode.
func (s *GatewayService) SendOTP(ctx context.Context, workspaceID string, otp *contracts.OTP) (*contracts.SendResult, error) {
	mode := s.resolveDispatchMode(ctx, workspaceID)
	return s.dispatch(ctx, workspaceID, mode, channelOTP,
		func() (*contracts.SendResult, error) { return s.writeOTPToInbox(ctx, workspaceID, otp) },
		func(intg *domain.Integration) (*contracts.SendResult, error) {
			sender, err := resolveOTPSender(s.otpCache, intg)
			if err != nil {
				return nil, err
			}
			return sender.Send(ctx, otp)
		},
	)
}

// CheckOTPStatus queries the delivery status of a previously sent OTP message
// identified by its request_id. It resolves the workspace's active OTP integration
// and delegates to the provider's CheckStatus if supported.
func (s *GatewayService) CheckOTPStatus(ctx context.Context, workspaceID, requestID string) (*contracts.VerificationStatus, error) {
	intg, err := s.activeIntegration(ctx, workspaceID, channelOTP)
	if err != nil {
		return nil, err
	}

	sender, err := resolveOTPSender(s.otpCache, intg)
	if err != nil {
		return nil, err
	}

	checker, ok := sender.(port.OTPStatusChecker)
	if !ok {
		return nil, fmt.Errorf("OTP provider %q does not support status checking", intg.ProviderName)
	}

	return checker.CheckStatus(ctx, requestID)
}

// RevokeOTP revokes a previously sent OTP verification message identified by
// request_id. It resolves the workspace's active OTP integration and delegates
// to the provider's Revoke if supported.
func (s *GatewayService) RevokeOTP(ctx context.Context, workspaceID, requestID string) (*contracts.SendResult, error) {
	intg, err := s.activeIntegration(ctx, workspaceID, channelOTP)
	if err != nil {
		return nil, err
	}

	sender, err := resolveOTPSender(s.otpCache, intg)
	if err != nil {
		return nil, err
	}

	revoker, ok := sender.(port.OTPRevoker)
	if !ok {
		return nil, fmt.Errorf("OTP provider %q does not support revocation", intg.ProviderName)
	}

	return revoker.Revoke(ctx, requestID)
}

// dispatch is the single entry point for all channel dispatch logic.
// It applies the workspace dispatch mode and calls the appropriate fn(s).
//
//   - writeToInbox: captures to in-process RAM (always available)
//   - sendViaProvider: instantiates provider from DB integration + sends
func (s *GatewayService) dispatch(
	ctx context.Context,
	workspaceID string,
	mode domain.MessageDispatchMode,
	channel string,
	writeToInbox func() (*contracts.SendResult, error),
	sendViaProvider func(*domain.Integration) (*contracts.SendResult, error),
) (*contracts.SendResult, error) {
	switch mode {
	case domain.DispatchMemoryOnly:
		r, err := writeToInbox()
		if err != nil {
			return nil, err
		}
		attachMeta(r, mode, channel, "")
		return r, nil

	case domain.DispatchProviderOnly:
		intg, err := s.activeIntegration(ctx, workspaceID, channel)
		if err != nil {
			return nil, err
		}
		// If the stored integration IS the memory provider, fall through to inbox.
		if intg.ProviderName == memoryProviderName {
			r, err := writeToInbox()
			if err != nil {
				return nil, err
			}
			attachMeta(r, mode, channel, intg.ID)
			return r, nil
		}
		r, err := sendViaProvider(intg)
		if err != nil {
			return nil, err
		}
		attachMeta(r, mode, channel, intg.ID)
		return r, nil

	case domain.DispatchMemoryAndProvider:
		intg, err := s.activeIntegration(ctx, workspaceID, channel)
		if err != nil {
			return nil, err
		}
		// If the integration is already memory, a single write is enough.
		if intg.ProviderName == memoryProviderName {
			r, err := writeToInbox()
			if err != nil {
				return nil, err
			}
			attachMeta(r, mode, channel, intg.ID)
			return r, nil
		}
		// Both paths: capture to inbox first (non-fatal), then send via provider.
		inboxResult, _ := writeToInbox() // inbox failure does not abort the send
		provResult, err := sendViaProvider(intg)
		if err != nil {
			return nil, err
		}
		if inboxResult != nil {
			if provResult.Meta == nil {
				provResult.Meta = make(map[string]string)
			}
			provResult.Meta["inbox_message_id"] = inboxResult.ID
		}
		attachMeta(provResult, mode, channel, intg.ID)
		return provResult, nil

	default:
		// Undefined modes fall back to memory_only (safe default).
		r, err := writeToInbox()
		if err != nil {
			return nil, err
		}
		attachMeta(r, domain.DispatchMemoryOnly, channel, "")
		return r, nil
	}
}

// resolveDispatchMode reads the workspace setting, defaulting gracefully.
func (s *GatewayService) resolveDispatchMode(ctx context.Context, workspaceID string) domain.MessageDispatchMode {
	if s.settings == nil {
		return domain.DefaultMessageDispatchMode()
	}
	v, err := s.settings.Get(ctx, workspaceID, domain.SettingKeyMessageDispatchMode)
	if err != nil || v == "" {
		return domain.DefaultMessageDispatchMode()
	}
	if m, ok := domain.ParseMessageDispatchMode(v); ok {
		return m
	}
	return domain.DefaultMessageDispatchMode()
}

// activeIntegration fetches the active (connected) integration for a workspace+channel.
func (s *GatewayService) activeIntegration(ctx context.Context, workspaceID, channel string) (*domain.Integration, error) {
	intg, err := s.integrations.GetActiveByWorkspaceAndChannel(ctx, workspaceID, channel)
	if err != nil {
		return nil, fmt.Errorf("no active %s integration for workspace %s: %w", channel, workspaceID, err)
	}
	return intg, nil
}

// inbox write helpers — each converts a domain error into meaningful context.

func (s *GatewayService) writeEmailToInbox(ctx context.Context, workspaceID string, email *contracts.Email) (*contracts.SendResult, error) {
	if s.inbox == nil {
		return nil, fmt.Errorf("inbox writer not configured")
	}
	id, err := s.inbox.WriteEmail(ctx, workspaceID, email)
	if err != nil {
		return nil, fmt.Errorf("write email to inbox: %w", err)
	}
	return &contracts.SendResult{ID: id, StatusCode: 200, Message: "captured in memory"}, nil
}

func (s *GatewayService) writeSMSToInbox(ctx context.Context, workspaceID string, sms *contracts.SMS) (*contracts.SendResult, error) {
	if s.inbox == nil {
		return nil, fmt.Errorf("inbox writer not configured")
	}
	id, err := s.inbox.WriteSMS(ctx, workspaceID, sms)
	if err != nil {
		return nil, fmt.Errorf("write SMS to inbox: %w", err)
	}
	return &contracts.SendResult{ID: id, StatusCode: 200, Message: "captured in memory"}, nil
}

func (s *GatewayService) writePushToInbox(ctx context.Context, workspaceID string, push *contracts.PushNotification) (*contracts.SendResult, error) {
	if s.inbox == nil {
		return nil, fmt.Errorf("inbox writer not configured")
	}
	id, err := s.inbox.WritePush(ctx, workspaceID, push)
	if err != nil {
		return nil, fmt.Errorf("write push to inbox: %w", err)
	}
	return &contracts.SendResult{ID: id, StatusCode: 200, Message: "captured in memory"}, nil
}

func (s *GatewayService) writeChatToInbox(ctx context.Context, workspaceID string, chat *contracts.ChatMessage) (*contracts.SendResult, error) {
	if s.inbox == nil {
		return nil, fmt.Errorf("inbox writer not configured")
	}
	id, err := s.inbox.WriteChat(ctx, workspaceID, chat)
	if err != nil {
		return nil, fmt.Errorf("write chat to inbox: %w", err)
	}
	return &contracts.SendResult{ID: id, StatusCode: 200, Message: "captured in memory"}, nil
}

func (s *GatewayService) writeOTPToInbox(ctx context.Context, workspaceID string, otp *contracts.OTP) (*contracts.SendResult, error) {
	if s.inbox == nil {
		return nil, fmt.Errorf("inbox writer not configured")
	}
	id, err := s.inbox.WriteOTP(ctx, workspaceID, otp)
	if err != nil {
		return nil, fmt.Errorf("write OTP to inbox: %w", err)
	}
	return &contracts.SendResult{ID: id, StatusCode: 200, Message: "captured in memory"}, nil
}

// attachMeta stamps standard dispatch metadata onto a result without allocating
// if Meta is already populated.
func attachMeta(r *contracts.SendResult, mode domain.MessageDispatchMode, channel, integrationID string) {
	if r == nil {
		return
	}
	if r.Meta == nil {
		r.Meta = make(map[string]string, 3)
	}
	r.Meta["dispatch_mode"] = string(mode)
	r.Meta["channel"] = channel
	if integrationID != "" {
		r.Meta["integration_id"] = integrationID
	}
}
