// Package service contains application services that orchestrate domain logic.
// Services depend only on port interfaces — never on infrastructure implementations.
package service

import (
	"context"
	"errors"
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
// Pass nil to disable memory capture — memory dispatch will then error when inbox is required.
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
	config, err := s.resolveDispatchConfig(ctx, workspaceID)
	if err != nil {
		return nil, err
	}
	return s.dispatch(ctx, workspaceID, config, channelEmail,
		func() (*contracts.SendResult, error) { return s.writeEmailToInbox(ctx, workspaceID, email) },
		func(sendCtx context.Context, intg domain.Integration) (*contracts.SendResult, error) {
			sender, err := resolveEmailSender(s.emailCache, &intg)
			if err != nil {
				return nil, err
			}
			return sender.Send(sendCtx, email)
		},
	)
}

// SendSMS dispatches an SMS for workspaceID according to its dispatch mode.
func (s *GatewayService) SendSMS(ctx context.Context, workspaceID string, sms contracts.SMS) (*contracts.SendResult, error) {
	config, err := s.resolveDispatchConfig(ctx, workspaceID)
	if err != nil {
		return nil, err
	}
	return s.dispatch(ctx, workspaceID, config, channelSMS,
		func() (*contracts.SendResult, error) { return s.writeSMSToInbox(ctx, workspaceID, sms) },
		func(sendCtx context.Context, intg domain.Integration) (*contracts.SendResult, error) {
			sender, err := resolveSMSSender(s.smsCache, &intg)
			if err != nil {
				return nil, err
			}
			return sender.Send(sendCtx, sms)
		},
	)
}

// SendPush dispatches a push notification for workspaceID according to its dispatch mode.
func (s *GatewayService) SendPush(ctx context.Context, workspaceID string, push contracts.PushNotification) (*contracts.SendResult, error) {
	config, err := s.resolveDispatchConfig(ctx, workspaceID)
	if err != nil {
		return nil, err
	}
	return s.dispatch(ctx, workspaceID, config, channelPush,
		func() (*contracts.SendResult, error) { return s.writePushToInbox(ctx, workspaceID, push) },
		func(sendCtx context.Context, intg domain.Integration) (*contracts.SendResult, error) {
			sender, err := resolvePushSender(s.pushCache, &intg)
			if err != nil {
				return nil, err
			}
			return sender.Send(sendCtx, push)
		},
	)
}

// SendChat dispatches a chat message for workspaceID according to its dispatch mode.
func (s *GatewayService) SendChat(ctx context.Context, workspaceID string, chat contracts.ChatMessage) (*contracts.SendResult, error) {
	config, err := s.resolveDispatchConfig(ctx, workspaceID)
	if err != nil {
		return nil, err
	}
	return s.dispatch(ctx, workspaceID, config, channelChat,
		func() (*contracts.SendResult, error) { return s.writeChatToInbox(ctx, workspaceID, chat) },
		func(sendCtx context.Context, intg domain.Integration) (*contracts.SendResult, error) {
			sender, err := resolveChatSender(s.chatCache, &intg)
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
//   - sendViaProvider: resolves provider from integration and sends
func (s *GatewayService) dispatch(
	ctx context.Context,
	workspaceID string,
	config domain.MessageDispatchConfig,
	channel string,
	writeToInbox func() (*contracts.SendResult, error),
	sendViaProvider func(context.Context, domain.Integration) (*contracts.SendResult, error),
) (*contracts.SendResult, error) {
	slog.InfoContext(ctx, "dispatching message", "dispatch_mode", config.Mode, "store_content", config.StoreMessageContent)

	var inboxResult *contracts.SendResult
	var provResult *contracts.SendResult
	var err error
	var integrationID, providerName string
	effectiveMode := config.Mode
	providerName = domain.ProviderNameMemory

	if config.RoutesViaProvider() {
		intg, errLookup := s.requireProviderIntegration(ctx, workspaceID, channel)
		if errLookup != nil {
			return attachMeta(nil, effectiveMode, config.StoreMessageContent, channel, integrationID, providerName), errLookup
		}

		integrationID = intg.ID
		providerName = intg.ProviderName

		providerCtx := applogger.WithProvider(ctx, providerName)
		provResult, err = sendViaProvider(providerCtx, intg)
		if err != nil {
			slog.WarnContext(providerCtx, "provider dispatch failed", "error", err, "integration_id", integrationID)
			return attachMeta(nil, effectiveMode, config.StoreMessageContent, channel, integrationID, providerName), err
		}
		slog.InfoContext(providerCtx, "message dispatched via provider", "message_id", provResult.ID)
	}

	if config.ShouldCaptureToInbox(effectiveMode) {
		inboxResult, err = writeToInbox()
		if err != nil {
			slog.ErrorContext(ctx, "inbox write failed", "error", err, "dispatch_mode", effectiveMode, "store_content", config.StoreMessageContent)
			if effectiveMode == domain.DispatchMemory {
				return attachMeta(nil, effectiveMode, config.StoreMessageContent, channel, integrationID, providerName), err
			}
			slog.WarnContext(ctx, "inbox capture failed after provider dispatch", "dispatch_mode", effectiveMode)
		} else {
			slog.InfoContext(ctx, "message captured in inbox", "message_id", inboxResult.ID, "store_content", config.StoreMessageContent)
		}
	}

	finalResult := provResult
	if finalResult == nil {
		finalResult = inboxResult
	} else if inboxResult != nil {
		if finalResult.Meta == nil {
			finalResult.Meta = make(map[string]string)
		}
		finalResult.Meta[contracts.MetaKeyInboxMessageID] = inboxResult.ID
	}

	return attachMeta(finalResult, effectiveMode, config.StoreMessageContent, channel, integrationID, providerName), nil
}

// ResolveDispatchConfig reads the workspace setting for outbound dispatch behavior.
func (s *GatewayService) ResolveDispatchConfig(ctx context.Context, workspaceID string) (domain.MessageDispatchConfig, error) {
	return s.resolveDispatchConfig(ctx, workspaceID)
}

// resolveDispatchConfig reads message_dispatch_mode and store_message_content from workspace settings.
func (s *GatewayService) resolveDispatchConfig(ctx context.Context, workspaceID string) (domain.MessageDispatchConfig, error) {
	if s.settings == nil {
		return domain.MessageDispatchConfig{
			Mode:                domain.DefaultMessageDispatchMode(),
			StoreMessageContent: false,
		}, nil
	}

	modeVal, err := s.settings.Get(ctx, workspaceID, domain.SettingKeyMessageDispatchMode)
	if err != nil && !errors.Is(err, port.ErrNotFound) {
		return domain.MessageDispatchConfig{}, fmt.Errorf("failed to lookup message dispatch mode: %w", err)
	}

	storeVal, err := s.settings.Get(ctx, workspaceID, domain.SettingKeyStoreMessageContent)
	if err != nil && !errors.Is(err, port.ErrNotFound) {
		return domain.MessageDispatchConfig{}, fmt.Errorf("failed to lookup store message content setting: %w", err)
	}

	return domain.ResolveMessageDispatchConfig(modeVal, storeVal), nil
}

// requireProviderIntegration loads and validates the workspace integration for provider dispatch.
func (s *GatewayService) requireProviderIntegration(ctx context.Context, workspaceID, channel string) (domain.Integration, error) {
	if s.integrations == nil {
		err := fmt.Errorf("integration repository not configured: %w", port.ErrInternal)
		slog.ErrorContext(ctx, "provider dispatch blocked",
			"reason", "repository_unavailable",
			"channel", channel,
			"error", err,
		)
		return domain.Integration{}, err
	}

	intg, err := s.integrations.GetActiveByWorkspaceAndChannel(ctx, workspaceID, channel)
	if err != nil {
		reason := "lookup_failed"
		log := slog.ErrorContext
		if errors.Is(err, port.ErrNotFound) {
			reason = "no_active_integration"
			log = slog.WarnContext
		}
		log(ctx, "provider dispatch blocked",
			"reason", reason,
			"channel", channel,
			"error", err,
		)
		return domain.Integration{}, fmt.Errorf("provider integration for %s channel: %w", channel, err)
	}
	if intg == nil {
		err := fmt.Errorf("provider integration for %s channel: %w", channel, port.ErrNotFound)
		slog.WarnContext(ctx, "provider dispatch blocked",
			"reason", "no_active_integration",
			"channel", channel,
			"error", err,
		)
		return domain.Integration{}, err
	}

	if err := domain.ValidateProviderIntegration(*intg); err != nil {
		slog.WarnContext(ctx, "provider dispatch blocked",
			"reason", "integration_not_ready",
			"channel", channel,
			"integration_id", intg.ID,
			"provider", intg.ProviderName,
			"status", intg.Status,
			"error", err,
		)
		return domain.Integration{}, err
	}

	return *intg, nil
}

// inbox write helpers — each converts a domain error into meaningful context.

func (s *GatewayService) writeEmailToInbox(ctx context.Context, workspaceID string, email contracts.Email) (*contracts.SendResult, error) {
	if s.inbox == nil {
		return nil, fmt.Errorf("inbox writer not configured")
	}
	email = s.enrichEmailFromIntegration(ctx, workspaceID, email)
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
func attachMeta(r *contracts.SendResult, effectiveMode domain.MessageDispatchMode, storeContent bool, channel, integrationID, providerName string) *contracts.SendResult {
	if r == nil {
		r = &contracts.SendResult{}
	}
	if r.Meta == nil {
		r.Meta = make(map[string]string, 5)
	}
	contracts.SetDispatchModeMeta(r, string(effectiveMode))
	contracts.SetStoreContentMeta(r, storeContent)
	r.Meta["channel"] = channel
	if integrationID != "" {
		r.Meta["integration_id"] = integrationID
	}
	if providerName != "" {
		r.Meta[contracts.MetaKeyProviderName] = providerName
	}
	return r
}

// RecordLog persists a MessageRequestLog entry. Repository errors are logged at the infra layer.
func (s *GatewayService) RecordLog(ctx context.Context, entry domain.MessageRequestLog) error {
	if s.logs == nil || entry.WorkspaceID == "" {
		return nil
	}
	return s.logs.Create(ctx, &entry)
}
