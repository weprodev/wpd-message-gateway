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
	cfg := s.resolveDispatchConfig(ctx, workspaceID)
	return s.dispatch(ctx, workspaceID, cfg, channelEmail,
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
	cfg := s.resolveDispatchConfig(ctx, workspaceID)
	return s.dispatch(ctx, workspaceID, cfg, channelSMS,
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
	cfg := s.resolveDispatchConfig(ctx, workspaceID)
	return s.dispatch(ctx, workspaceID, cfg, channelPush,
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
	cfg := s.resolveDispatchConfig(ctx, workspaceID)
	return s.dispatch(ctx, workspaceID, cfg, channelChat,
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
	cfg domain.MessageDispatchConfig,
	channel string,
	writeToInbox func() (*contracts.SendResult, error),
	sendViaProvider func(context.Context, *domain.Integration) (*contracts.SendResult, error),
) (*contracts.SendResult, error) {
	apiValue := cfg.APIValue()
	slog.InfoContext(ctx, "dispatching message", "workspace_id", workspaceID, "dispatch_mode", apiValue, "channel", channel)
	switch cfg.Mode {
	case domain.DispatchMemory:
		r, err := writeToInbox()
		if err != nil {
			slog.ErrorContext(ctx, "inbox write failed", "error", err, "workspace_id", workspaceID, "channel", channel)
			return attachMeta(nil, cfg, channel, "", memoryProviderName), err
		}
		attachMeta(r, cfg, channel, "", memoryProviderName)
		slog.InfoContext(ctx, "message dispatched via memory", "workspace_id", workspaceID, "channel", channel, "message_id", r.ID, "retain_request_log", cfg.RetainRequestLog)
		return r, nil

	case domain.DispatchProvider:
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
				return attachMeta(nil, cfg, channel, intg.ID, intg.ProviderName), err
			}
			attachMeta(r, cfg, channel, intg.ID, intg.ProviderName)
			slog.InfoContext(ctx, "message dispatched via memory fallback", "workspace_id", workspaceID, "channel", channel, "message_id", r.ID)
			return r, nil
		}
		providerCtx := applogger.WithProvider(ctx, intg.ProviderName)
		r, err := sendViaProvider(providerCtx, intg)
		if err != nil {
			slog.ErrorContext(ctx, "provider send failed", "error", err, "workspace_id", workspaceID, "channel", channel, "provider", intg.ProviderName)
			return attachMeta(nil, cfg, channel, intg.ID, intg.ProviderName), err
		}
		attachMeta(r, cfg, channel, intg.ID, intg.ProviderName)
		slog.InfoContext(ctx, "message dispatched via provider", "workspace_id", workspaceID, "channel", channel, "provider", intg.ProviderName, "message_id", r.ID, "dispatch_mode", apiValue)
		return r, nil

	default:
		// Undefined modes fall back to memory_only (safe default).
		slog.WarnContext(ctx, "unknown dispatch mode, falling back to memory only", "workspace_id", workspaceID, "mode", cfg.Mode)
		fallback := domain.DefaultMessageDispatchConfig()
		r, err := writeToInbox()
		if err != nil {
			slog.ErrorContext(ctx, "inbox write failed (fallback)", "error", err, "workspace_id", workspaceID, "channel", channel)
			return attachMeta(nil, fallback, channel, "", memoryProviderName), err
		}
		attachMeta(r, fallback, channel, "", memoryProviderName)
		return r, nil
	}
}

// ResolveDispatchConfig reads the workspace setting for outbound dispatch behavior.
func (s *GatewayService) ResolveDispatchConfig(ctx context.Context, workspaceID string) domain.MessageDispatchConfig {
	return s.resolveDispatchConfig(ctx, workspaceID)
}

// ShouldRetainForWorkspace reports whether request logs for the workspace should be marked retained.
func (s *GatewayService) ShouldRetainForWorkspace(ctx context.Context, workspaceID string) bool {
	return s.resolveDispatchConfig(ctx, workspaceID).RetainRequestLog
}

// resolveDispatchConfig reads message_dispatch_mode from workspace settings.
func (s *GatewayService) resolveDispatchConfig(ctx context.Context, workspaceID string) domain.MessageDispatchConfig {
	if s.settings == nil {
		return domain.DefaultMessageDispatchConfig()
	}
	v, err := s.settings.Get(ctx, workspaceID, domain.SettingKeyMessageDispatchMode)
	if err != nil || v == "" {
		return domain.DefaultMessageDispatchConfig()
	}
	if cfg, ok := domain.SettingValueToDispatchConfig(v); ok {
		return cfg
	}
	return domain.DefaultMessageDispatchConfig()
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

// attachMeta stamps standard dispatch metadata onto r. When r is nil (error paths
// with no provider payload), it allocates a minimal SendResult so callers can
// still read provider_name via contracts.ProviderNameFromResult.
func attachMeta(r *contracts.SendResult, cfg domain.MessageDispatchConfig, channel, integrationID, providerName string) *contracts.SendResult {
	if r == nil {
		r = &contracts.SendResult{}
	}
	if r.Meta == nil {
		r.Meta = make(map[string]string, 4)
	}
	r.Meta["dispatch_mode"] = string(cfg.APIValue())
	r.Meta["channel"] = channel
	if integrationID != "" {
		r.Meta["integration_id"] = integrationID
	}
	if providerName != "" {
		r.Meta["provider_name"] = providerName
	}
	return r
}

// RecordLog persists a MessageRequestLog entry. Failures are logged with slog.
func (s *GatewayService) RecordLog(ctx context.Context, entry *domain.MessageRequestLog) error {
	if s.logs == nil || entry.WorkspaceID == "" {
		return nil
	}
	return s.logs.Create(ctx, entry)
}
