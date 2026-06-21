// Package service contains application services that orchestrate domain logic.
// Services depend only on port interfaces — never on infrastructure implementations.
package service

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

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

	maxStoredDispatchErrorLen = 1024
)

// GatewayService orchestrates message dispatch for a single workspace.
//
// It resolves the workspace dispatch mode from WorkspaceSettingsRepository,
// then routes the message to the in-process InboxWriter, durable StoredMessageWriter,
// a real provider fetched via IntegrationRepository, or a combination — depending on the mode.
//
// Provider instances are cached per (integrationID, updatedAt) so HTTP clients
// are reused across requests while auto-invalidating on config changes.
type GatewayService struct {
	integrations   port.IntegrationRepository
	templates      port.TemplateRepository
	settings       port.WorkspaceSettingsRepository
	inbox          port.InboxWriter
	storedMessages port.StoredMessageWriter
	logs           port.MessageRequestLogRepository

	emailCache *emailSenderCache
	smsCache   *smsSenderCache
	pushCache  *pushSenderCache
	chatCache  *chatSenderCache
}

// NewGatewayService constructs a GatewayService.
//
// inbox is the in-process capture store (implements port.InboxWriter).
// storedMessages persists payloads to PostgreSQL for provider_and_database mode.
// Pass nil to disable durable capture — provider_and_database dispatch will then warn and skip DB write.
func NewGatewayService(
	integrations port.IntegrationRepository,
	templates port.TemplateRepository,
	settings port.WorkspaceSettingsRepository,
	inbox port.InboxWriter,
	storedMessages port.StoredMessageWriter,
	logs port.MessageRequestLogRepository,
) *GatewayService {
	return &GatewayService{
		integrations:   integrations,
		templates:      templates,
		settings:       settings,
		inbox:          inbox,
		storedMessages: storedMessages,
		logs:           logs,
		emailCache:     newProviderCache[contracts.EmailSender](),
		smsCache:       newProviderCache[contracts.SMSSender](),
		pushCache:      newProviderCache[contracts.PushSender](),
		chatCache:      newProviderCache[contracts.ChatSender](),
	}
}

// SendEmail dispatches an email for workspaceID according to its dispatch mode.
func (s *GatewayService) SendEmail(ctx context.Context, workspaceID string, email contracts.Email) (*contracts.SendResult, error) {
	mode := s.resolveDispatchMode(ctx, workspaceID)
	return s.dispatch(ctx, workspaceID, mode, channelEmail,
		func() (*contracts.SendResult, error) { return s.writeEmailToInbox(ctx, workspaceID, email) },
		func() (*contracts.SendResult, error) { return s.writeEmailToArchive(ctx, workspaceID, email) },
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
		func() (*contracts.SendResult, error) { return s.writeSMSToArchive(ctx, workspaceID, sms) },
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
		func() (*contracts.SendResult, error) { return s.writePushToArchive(ctx, workspaceID, push) },
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
		func() (*contracts.SendResult, error) { return s.writeChatToArchive(ctx, workspaceID, chat) },
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
//   - writeToInbox: captures to in-process RAM
//   - writeToArchive: persists payload to PostgreSQL
//   - sendViaProvider: instantiates provider from DB integration + sends
func (s *GatewayService) dispatch(
	ctx context.Context,
	workspaceID string,
	mode domain.MessageDispatchMode,
	channel string,
	writeToInbox func() (*contracts.SendResult, error),
	writeToArchive func() (*contracts.SendResult, error),
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

	case domain.DispatchProviderOnly:
		return s.sendViaActiveProvider(ctx, workspaceID, mode, channel, writeToInbox, sendViaProvider, nil, "", "message dispatched via provider only")

	case domain.DispatchMemoryAndProvider:
		return s.sendViaActiveProvider(ctx, workspaceID, mode, channel, writeToInbox, sendViaProvider, writeToInbox, "inbox_message_id", "message dispatched via memory and provider")

	case domain.DispatchProviderAndDatabase:
		return s.sendViaActiveProvider(ctx, workspaceID, mode, channel, writeToInbox, sendViaProvider, writeToArchive, "stored_message_id", "message dispatched via provider and database")

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

// sendViaActiveProvider sends via the workspace integration; persist failure aborts provider_and_database only.
func (s *GatewayService) sendViaActiveProvider(
	ctx context.Context,
	workspaceID string,
	mode domain.MessageDispatchMode,
	channel string,
	writeToInbox func() (*contracts.SendResult, error),
	sendViaProvider func(context.Context, *domain.Integration) (*contracts.SendResult, error),
	persist func() (*contracts.SendResult, error),
	persistMetaKey string,
	successLogMsg string,
) (*contracts.SendResult, error) {
	intg, err := s.activeIntegration(ctx, workspaceID, channel)
	if err != nil {
		slog.ErrorContext(ctx, "provider lookup failed", "error", err, "workspace_id", workspaceID, "channel", channel)
		return nil, err
	}
	if intg.ProviderName == memoryProviderName && mode != domain.DispatchProviderAndDatabase {
		r, err := writeToInbox()
		if err != nil {
			slog.ErrorContext(ctx, "inbox write failed (fallback)", "error", err, "workspace_id", workspaceID, "channel", channel)
			return nil, err
		}
		attachMeta(r, mode, channel, intg.ID, intg.ProviderName)
		slog.InfoContext(ctx, "message dispatched via memory fallback", "workspace_id", workspaceID, "channel", channel, "message_id", r.ID)
		return r, nil
	}

	if persist != nil {
		persistResult, err := persist()
		if err != nil {
			if mode == domain.DispatchProviderAndDatabase {
				slog.ErrorContext(ctx, "archive persist failed", "error", err, "workspace_id", workspaceID, "channel", channel, "dispatch_mode", mode)
				return nil, err
			}
			slog.WarnContext(ctx, "persist failed (non-fatal)", "error", err, "workspace_id", workspaceID, "channel", channel, "dispatch_mode", mode)
		}
		providerCtx := applogger.WithProvider(ctx, intg.ProviderName)
		provResult, err := sendViaProvider(providerCtx, intg)
		if persistMetaKey == "stored_message_id" && persistResult != nil {
			s.recordStoredMessageOutcome(ctx, persistResult.ID, provResult, err)
		}
		if err != nil {
			slog.ErrorContext(ctx, "provider send failed", "error", err, "workspace_id", workspaceID, "channel", channel, "provider", intg.ProviderName)
			return nil, err
		}
		if persistResult != nil && persistMetaKey != "" {
			if provResult.Meta == nil {
				provResult.Meta = make(map[string]string)
			}
			provResult.Meta[persistMetaKey] = persistResult.ID
		}
		attachMeta(provResult, mode, channel, intg.ID, intg.ProviderName)
		slog.InfoContext(ctx, successLogMsg, "workspace_id", workspaceID, "channel", channel, "provider", intg.ProviderName, "message_id", provResult.ID)
		return provResult, nil
	}

	providerCtx := applogger.WithProvider(ctx, intg.ProviderName)
	r, err := sendViaProvider(providerCtx, intg)
	if err != nil {
		slog.ErrorContext(ctx, "provider send failed", "error", err, "workspace_id", workspaceID, "channel", channel, "provider", intg.ProviderName)
		return nil, err
	}
	attachMeta(r, mode, channel, intg.ID, intg.ProviderName)
	slog.InfoContext(ctx, successLogMsg, "workspace_id", workspaceID, "channel", channel, "provider", intg.ProviderName, "message_id", r.ID)
	return r, nil
}

// resolveDispatchMode reads the workspace setting, defaulting gracefully.
func (s *GatewayService) resolveDispatchMode(ctx context.Context, workspaceID string) domain.MessageDispatchMode {
	if s.settings == nil {
		return domain.DefaultMessageDispatchMode()
	}
	if mode := s.dispatchModeFromSetting(ctx, workspaceID, domain.SettingKeyDataRetention); mode != "" {
		return mode
	}
	if mode := s.dispatchModeFromSetting(ctx, workspaceID, domain.SettingKeyMessageDispatchMode); mode != "" {
		return mode
	}
	return domain.DefaultMessageDispatchMode()
}

func (s *GatewayService) dispatchModeFromSetting(ctx context.Context, workspaceID, key string) domain.MessageDispatchMode {
	v, err := s.settings.Get(ctx, workspaceID, key)
	if err != nil || v == "" {
		return ""
	}
	if m, ok := domain.DataRetentionValueToDispatchMode(v); ok {
		return m
	}
	return ""
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

func (s *GatewayService) writeEmailToArchive(ctx context.Context, workspaceID string, email contracts.Email) (*contracts.SendResult, error) {
	if s.storedMessages == nil {
		return nil, fmt.Errorf("stored message writer not configured")
	}
	id, err := s.storedMessages.WriteEmail(ctx, workspaceID, email)
	if err != nil {
		return nil, fmt.Errorf("write email to archive: %w", err)
	}
	return &contracts.SendResult{ID: id, StatusCode: 200, Message: "stored in database"}, nil
}

func (s *GatewayService) writeSMSToArchive(ctx context.Context, workspaceID string, sms contracts.SMS) (*contracts.SendResult, error) {
	if s.storedMessages == nil {
		return nil, fmt.Errorf("stored message writer not configured")
	}
	id, err := s.storedMessages.WriteSMS(ctx, workspaceID, sms)
	if err != nil {
		return nil, fmt.Errorf("write SMS to archive: %w", err)
	}
	return &contracts.SendResult{ID: id, StatusCode: 200, Message: "stored in database"}, nil
}

func (s *GatewayService) writePushToArchive(ctx context.Context, workspaceID string, push contracts.PushNotification) (*contracts.SendResult, error) {
	if s.storedMessages == nil {
		return nil, fmt.Errorf("stored message writer not configured")
	}
	id, err := s.storedMessages.WritePush(ctx, workspaceID, push)
	if err != nil {
		return nil, fmt.Errorf("write push to archive: %w", err)
	}
	return &contracts.SendResult{ID: id, StatusCode: 200, Message: "stored in database"}, nil
}

func (s *GatewayService) writeChatToArchive(ctx context.Context, workspaceID string, chat contracts.ChatMessage) (*contracts.SendResult, error) {
	if s.storedMessages == nil {
		return nil, fmt.Errorf("stored message writer not configured")
	}
	id, err := s.storedMessages.WriteChat(ctx, workspaceID, chat)
	if err != nil {
		return nil, fmt.Errorf("write chat to archive: %w", err)
	}
	return &contracts.SendResult{ID: id, StatusCode: 200, Message: "stored in database"}, nil
}

func (s *GatewayService) recordStoredMessageOutcome(
	ctx context.Context,
	storedMessageID string,
	provResult *contracts.SendResult,
	sendErr error,
) {
	if s.storedMessages == nil || storedMessageID == "" {
		return
	}
	outcome := domain.StoredMessageDispatchOutcome{
		DispatchedAt: time.Now().UTC(),
	}
	if sendErr != nil {
		outcome.Status = domain.StoredMessageDispatchFailed
		outcome.DispatchError = truncateStoredDispatchError(sendErr.Error())
	} else {
		outcome.Status = domain.StoredMessageDispatchSent
		if provResult != nil {
			outcome.ProviderMessageID = provResult.ID
			outcome.ProviderStatusCode = provResult.StatusCode
		}
	}
	if err := s.storedMessages.RecordDispatchOutcome(ctx, storedMessageID, outcome); err != nil {
		slog.WarnContext(ctx, "stored message dispatch outcome update failed (non-fatal)",
			"error", err, "stored_message_id", storedMessageID, "dispatch_status", outcome.Status)
	}
}

func truncateStoredDispatchError(msg string) string {
	if len(msg) <= maxStoredDispatchErrorLen {
		return msg
	}
	return strings.TrimSpace(msg[:maxStoredDispatchErrorLen])
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

// RecordLog persists a MessageRequestLog entry. Failures are logged with slog.
func (s *GatewayService) RecordLog(ctx context.Context, entry *domain.MessageRequestLog) error {
	if s.logs == nil || entry.WorkspaceID == "" {
		return nil
	}
	return s.logs.Create(ctx, entry)
}
