// Package service contains application services that orchestrate domain logic.
// Services depend only on port interfaces — never on infrastructure implementations.
package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/weprodev/wpd-message-gateway/internal/core/domain"
	"github.com/weprodev/wpd-message-gateway/internal/core/port"
	applogger "github.com/weprodev/wpd-message-gateway/internal/infrastructure/logger"
	"github.com/weprodev/wpd-message-gateway/pkg/contracts"
)

const (
	channelEmail = "email"
	channelSMS   = "sms"
	channelPush  = "push"
	channelChat  = "chat"

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
	logs         port.MessageRequestLogRepository

	emailCache *emailSenderCache
	smsCache   *smsSenderCache
	pushCache  *pushSenderCache
	chatCache  *chatSenderCache
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
	logs port.MessageRequestLogRepository,
) *GatewayService {
	return &GatewayService{
		integrations: integrations,
		templates:    templates,
		settings:     settings,
		inbox:        inbox,
		logs:         logs,
		emailCache:   newProviderCache[contracts.EmailSender](),
		smsCache:     newProviderCache[contracts.SMSSender](),
		pushCache:    newProviderCache[contracts.PushSender](),
		chatCache:    newProviderCache[contracts.ChatSender](),
	}
}

// SendEmail dispatches an email for workspaceID according to its dispatch mode.
func (s *GatewayService) SendEmail(ctx context.Context, workspaceID string, email contracts.Email) (*contracts.SendResult, error) {
	mode := s.resolveDispatchMode(ctx, workspaceID)
	return s.dispatch(ctx, workspaceID, mode, channelEmail,
		func() (*contracts.SendResult, error) { return s.writeEmailToInbox(ctx, workspaceID, email) },
		func(sendCtx context.Context, intg *domain.Integration) (*contracts.SendResult, error) {
			sender, err := resolveEmailSender(s.emailCache, intg)
			if err != nil {
				return nil, err
			}
			return sender.Send(sendCtx, email)
		},
	)
}

// SendSMS dispatches an SMS for workspaceID according to its dispatch mode.
func (s *GatewayService) SendSMS(ctx context.Context, workspaceID string, sms contracts.SMS) (*contracts.SendResult, error) {
	mode := s.resolveDispatchMode(ctx, workspaceID)
	return s.dispatch(ctx, workspaceID, mode, channelSMS,
		func() (*contracts.SendResult, error) { return s.writeSMSToInbox(ctx, workspaceID, sms) },
		func(sendCtx context.Context, intg *domain.Integration) (*contracts.SendResult, error) {
			sender, err := resolveSMSSender(s.smsCache, intg)
			if err != nil {
				return nil, err
			}
			return sender.Send(sendCtx, sms)
		},
	)
}

// SendPush dispatches a push notification for workspaceID according to its dispatch mode.
func (s *GatewayService) SendPush(ctx context.Context, workspaceID string, push contracts.PushNotification) (*contracts.SendResult, error) {
	mode := s.resolveDispatchMode(ctx, workspaceID)
	return s.dispatch(ctx, workspaceID, mode, channelPush,
		func() (*contracts.SendResult, error) { return s.writePushToInbox(ctx, workspaceID, push) },
		func(sendCtx context.Context, intg *domain.Integration) (*contracts.SendResult, error) {
			sender, err := resolvePushSender(s.pushCache, intg)
			if err != nil {
				return nil, err
			}
			return sender.Send(sendCtx, push)
		},
	)
}

// SendChat dispatches a chat message for workspaceID according to its dispatch mode.
func (s *GatewayService) SendChat(ctx context.Context, workspaceID string, chat contracts.ChatMessage) (*contracts.SendResult, error) {
	mode := s.resolveDispatchMode(ctx, workspaceID)
	return s.dispatch(ctx, workspaceID, mode, channelChat,
		func() (*contracts.SendResult, error) { return s.writeChatToInbox(ctx, workspaceID, chat) },
		func(sendCtx context.Context, intg *domain.Integration) (*contracts.SendResult, error) {
			sender, err := resolveChatSender(s.chatCache, intg)
			if err != nil {
				return nil, err
			}
			return sender.Send(sendCtx, chat)
		},
	)
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
	sendViaProvider func(context.Context, *domain.Integration) (*contracts.SendResult, error),
) (*contracts.SendResult, error) {
	slog.InfoContext(ctx, "dispatching message", "workspace_id", workspaceID, "dispatch_mode", mode, "channel", channel)
	switch mode {
	case domain.DispatchMemoryOnly:
		r, err := writeToInbox()
		if err != nil {
			slog.ErrorContext(ctx, "inbox write failed", "error", err, "workspace_id", workspaceID, "channel", channel)
			return nil, err
		}
		attachMeta(r, mode, channel, "", memoryProviderName)
		slog.InfoContext(ctx, "message dispatched via memory only", "workspace_id", workspaceID, "channel", channel, "message_id", r.ID)
		return r, nil

	case domain.DispatchProviderOnly, domain.DispatchProviderAndDatabase:
		intg, err := s.activeIntegration(ctx, workspaceID, channel)
		if err != nil {
			slog.ErrorContext(ctx, "provider lookup failed", "error", err, "workspace_id", workspaceID, "channel", channel)
			return nil, err
		}
		// If the stored integration IS the memory provider, fall through to inbox.
		if intg.ProviderName == memoryProviderName {
			r, err := writeToInbox()
			if err != nil {
				slog.ErrorContext(ctx, "inbox write failed (fallback)", "error", err, "workspace_id", workspaceID, "channel", channel)
				return nil, err
			}
			attachMeta(r, mode, channel, intg.ID, intg.ProviderName)
			slog.InfoContext(ctx, "message dispatched via memory fallback", "workspace_id", workspaceID, "channel", channel, "message_id", r.ID)
			return r, nil
		}
		providerCtx := applogger.WithProvider(ctx, intg.ProviderName)
		r, err := sendViaProvider(providerCtx, intg)
		if err != nil {
			slog.ErrorContext(ctx, "provider send failed", "error", err, "workspace_id", workspaceID, "channel", channel, "provider", intg.ProviderName)
			return nil, err
		}
		attachMeta(r, mode, channel, intg.ID, intg.ProviderName)
		slog.InfoContext(ctx, "message dispatched via provider", "workspace_id", workspaceID, "channel", channel, "provider", intg.ProviderName, "message_id", r.ID, "dispatch_mode", mode)
		return r, nil

	case domain.DispatchMemoryAndProvider:
		intg, err := s.activeIntegration(ctx, workspaceID, channel)
		if err != nil {
			slog.ErrorContext(ctx, "provider lookup failed", "error", err, "workspace_id", workspaceID, "channel", channel)
			return nil, err
		}
		// If the integration is already memory, a single write is enough.
		if intg.ProviderName == memoryProviderName {
			r, err := writeToInbox()
			if err != nil {
				slog.ErrorContext(ctx, "inbox write failed (fallback)", "error", err, "workspace_id", workspaceID, "channel", channel)
				return nil, err
			}
			attachMeta(r, mode, channel, intg.ID, intg.ProviderName)
			slog.InfoContext(ctx, "message dispatched via memory fallback", "workspace_id", workspaceID, "channel", channel, "message_id", r.ID)
			return r, nil
		}
		// Both paths: capture to inbox first (non-fatal), then send via provider.
		inboxResult, err := writeToInbox()
		if err != nil {
			slog.WarnContext(ctx, "inbox write failed (non-fatal)", "error", err, "workspace_id", workspaceID, "channel", channel)
		}
		providerCtx := applogger.WithProvider(ctx, intg.ProviderName)
		provResult, err := sendViaProvider(providerCtx, intg)
		if err != nil {
			slog.ErrorContext(ctx, "provider send failed", "error", err, "workspace_id", workspaceID, "channel", channel, "provider", intg.ProviderName)
			return nil, err
		}
		if inboxResult != nil {
			if provResult.Meta == nil {
				provResult.Meta = make(map[string]string)
			}
			provResult.Meta["inbox_message_id"] = inboxResult.ID
		}
		attachMeta(provResult, mode, channel, intg.ID, intg.ProviderName)
		slog.InfoContext(ctx, "message dispatched via memory and provider", "workspace_id", workspaceID, "channel", channel, "provider", intg.ProviderName, "message_id", provResult.ID)
		return provResult, nil

	default:
		// Undefined modes fall back to memory_only (safe default).
		slog.WarnContext(ctx, "unknown dispatch mode, falling back to memory only", "workspace_id", workspaceID, "mode", mode)
		r, err := writeToInbox()
		if err != nil {
			slog.ErrorContext(ctx, "inbox write failed (fallback)", "error", err, "workspace_id", workspaceID, "channel", channel)
			return nil, err
		}
		attachMeta(r, domain.DispatchMemoryOnly, channel, "", memoryProviderName)
		return r, nil
	}
}

// ResolveDispatchMode reads the workspace setting for outbound dispatch behavior.
func (s *GatewayService) ResolveDispatchMode(ctx context.Context, workspaceID string) domain.MessageDispatchMode {
	return s.resolveDispatchMode(ctx, workspaceID)
}

// ShouldRetainForWorkspace reports whether request logs for the workspace should be marked retained.
func (s *GatewayService) ShouldRetainForWorkspace(ctx context.Context, workspaceID string) bool {
	return domain.ShouldRetainRequestLog(s.resolveDispatchMode(ctx, workspaceID))
}

// resolveDispatchMode reads message_dispatch_mode from workspace settings.
func (s *GatewayService) resolveDispatchMode(ctx context.Context, workspaceID string) domain.MessageDispatchMode {
	if s.settings == nil {
		return domain.DefaultMessageDispatchMode()
	}
	v, err := s.settings.Get(ctx, workspaceID, domain.SettingKeyMessageDispatchMode)
	if err != nil || v == "" {
		return domain.DefaultMessageDispatchMode()
	}
	if mode, ok := domain.SettingValueToDispatchMode(v); ok {
		return mode
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

func (s *GatewayService) writeEmailToInbox(ctx context.Context, workspaceID string, email contracts.Email) (*contracts.SendResult, error) {
	if s.inbox == nil {
		return nil, fmt.Errorf("inbox writer not configured")
	}
	id, err := s.inbox.WriteEmail(ctx, workspaceID, email)
	if err != nil {
		return nil, fmt.Errorf("write email to inbox: %w", err)
	}
	return &contracts.SendResult{ID: id, StatusCode: 200, Message: "captured in memory"}, nil
}

func (s *GatewayService) writeSMSToInbox(ctx context.Context, workspaceID string, sms contracts.SMS) (*contracts.SendResult, error) {
	if s.inbox == nil {
		return nil, fmt.Errorf("inbox writer not configured")
	}
	id, err := s.inbox.WriteSMS(ctx, workspaceID, sms)
	if err != nil {
		return nil, fmt.Errorf("write SMS to inbox: %w", err)
	}
	return &contracts.SendResult{ID: id, StatusCode: 200, Message: "captured in memory"}, nil
}

func (s *GatewayService) writePushToInbox(ctx context.Context, workspaceID string, push contracts.PushNotification) (*contracts.SendResult, error) {
	if s.inbox == nil {
		return nil, fmt.Errorf("inbox writer not configured")
	}
	id, err := s.inbox.WritePush(ctx, workspaceID, push)
	if err != nil {
		return nil, fmt.Errorf("write push to inbox: %w", err)
	}
	return &contracts.SendResult{ID: id, StatusCode: 200, Message: "captured in memory"}, nil
}

func (s *GatewayService) writeChatToInbox(ctx context.Context, workspaceID string, chat contracts.ChatMessage) (*contracts.SendResult, error) {
	if s.inbox == nil {
		return nil, fmt.Errorf("inbox writer not configured")
	}
	id, err := s.inbox.WriteChat(ctx, workspaceID, chat)
	if err != nil {
		return nil, fmt.Errorf("write chat to inbox: %w", err)
	}
	return &contracts.SendResult{ID: id, StatusCode: 200, Message: "captured in memory"}, nil
}

// attachMeta stamps standard dispatch metadata onto a result without allocating
// if Meta is already populated.
func attachMeta(r *contracts.SendResult, mode domain.MessageDispatchMode, channel, integrationID, providerName string) {
	if r == nil {
		return
	}
	if r.Meta == nil {
		r.Meta = make(map[string]string, 4)
	}
	r.Meta["dispatch_mode"] = string(mode)
	r.Meta["channel"] = channel
	if integrationID != "" {
		r.Meta["integration_id"] = integrationID
	}
	if providerName != "" {
		r.Meta["provider_name"] = providerName
	}
}

// ProviderNameForLog returns the provider name to persist on message_request_logs.
// It prefers dispatch metadata on the result and falls back to the active integration.
func (s *GatewayService) ProviderNameForLog(ctx context.Context, workspaceID, channel string, result *contracts.SendResult) string {
	if name := providerNameFromResult(result); name != "" {
		return name
	}
	switch s.resolveDispatchMode(ctx, workspaceID) {
	case domain.DispatchMemoryOnly:
		return memoryProviderName
	default:
		if s.integrations == nil {
			return ""
		}
		intg, err := s.activeIntegration(ctx, workspaceID, channel)
		if err != nil {
			return ""
		}
		return intg.ProviderName
	}
}

func providerNameFromResult(r *contracts.SendResult) string {
	if r == nil || r.Meta == nil {
		return ""
	}
	return r.Meta["provider_name"]
}

// RecordLog persists a MessageRequestLog entry. Failures are logged with slog.
func (s *GatewayService) RecordLog(ctx context.Context, entry *domain.MessageRequestLog) error {
	if s.logs == nil || entry.WorkspaceID == "" {
		return nil
	}
	return s.logs.Create(ctx, entry)
}
