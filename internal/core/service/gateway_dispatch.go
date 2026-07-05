package service

import (
	"context"
	"log/slog"

	"github.com/weprodev/wpd-message-gateway/internal/core/domain"
	"github.com/weprodev/wpd-message-gateway/internal/core/logctx"
	"github.com/weprodev/wpd-message-gateway/pkg/contracts"
)

type dispatchMeta struct {
	provResult    *contracts.SendResult
	inboxResult   *contracts.SendResult
	integrationID string
	providerName  string
}

// routeMessage delivers the outbound message according to dispatch mode.
// Portal inbox capture is handled separately when store_message_content is enabled.
func (s *GatewayService) routeMessage(
	ctx context.Context,
	workspaceID string,
	config domain.MessageDispatchConfig,
	channel string,
	sendViaProvider func(context.Context, domain.Integration) (*contracts.SendResult, error),
	meta *dispatchMeta,
) error {
	if config.RoutesViaProvider() {
		return s.routeViaIntegration(ctx, workspaceID, channel, sendViaProvider, meta)
	}
	if config.ShouldCaptureToInbox() {
		return nil
	}
	return s.routeViaMemoryProvider(ctx, sendViaProvider, meta)
}

func (s *GatewayService) routeViaIntegration(
	ctx context.Context,
	workspaceID, channel string,
	sendViaProvider func(context.Context, domain.Integration) (*contracts.SendResult, error),
	meta *dispatchMeta,
) error {
	intg, err := s.requireProviderIntegration(ctx, workspaceID, channel)
	if err != nil {
		return err
	}

	meta.integrationID = intg.ID
	meta.providerName = intg.ProviderName

	providerCtx := logctx.WithProvider(ctx, meta.providerName)
	meta.provResult, err = sendViaProvider(providerCtx, intg)
	if err != nil {
		slog.WarnContext(providerCtx, "provider dispatch failed", "error", err, "integration_id", meta.integrationID)
		return err
	}
	slog.InfoContext(providerCtx, "message dispatched via provider", "message_id", meta.provResult.ID)
	return nil
}

func (s *GatewayService) routeViaMemoryProvider(
	ctx context.Context,
	sendViaProvider func(context.Context, domain.Integration) (*contracts.SendResult, error),
	meta *dispatchMeta,
) error {
	intg := memoryDispatchIntegration()
	providerCtx := logctx.WithProvider(ctx, domain.ProviderNameMemory)

	var err error
	meta.provResult, err = sendViaProvider(providerCtx, intg)
	if err != nil {
		slog.WarnContext(providerCtx, "memory provider dispatch failed", "error", err)
		return err
	}

	meta.providerName = domain.ProviderNameMemory
	slog.InfoContext(providerCtx, "message dispatched via memory provider", "message_id", meta.provResult.ID)
	return nil
}

func (s *GatewayService) captureToInboxIfEnabled(
	ctx context.Context,
	config domain.MessageDispatchConfig,
	writeToInbox func() (*contracts.SendResult, error),
	meta *dispatchMeta,
) error {
	if !config.ShouldCaptureToInbox() {
		return nil
	}

	result, err := writeToInbox()
	if err != nil {
		slog.ErrorContext(ctx, "inbox write failed", "error", err, "dispatch_mode", config.Mode)
		if config.RoutesViaProvider() {
			slog.WarnContext(ctx, "inbox capture failed after provider dispatch", "dispatch_mode", config.Mode)
			return nil
		}
		return err
	}

	meta.inboxResult = result
	slog.InfoContext(ctx, "message captured in inbox", "message_id", result.ID)
	return nil
}

func mergeDispatchResults(meta dispatchMeta) *contracts.SendResult {
	if meta.provResult == nil {
		return meta.inboxResult
	}
	if meta.inboxResult == nil {
		return meta.provResult
	}
	if meta.provResult.Meta == nil {
		meta.provResult.Meta = make(map[string]string)
	}
	meta.provResult.Meta[contracts.MetaKeyInboxMessageID] = meta.inboxResult.ID
	return meta.provResult
}

func attachDispatchMeta(
	result *contracts.SendResult,
	config domain.MessageDispatchConfig,
	channel string,
	meta dispatchMeta,
) *contracts.SendResult {
	if result == nil {
		result = &contracts.SendResult{}
	}
	if result.Meta == nil {
		result.Meta = make(map[string]string, 5)
	}
	contracts.SetDispatchModeMeta(result, string(config.Mode))
	contracts.SetStoreContentMeta(result, config.StoreMessageContent)
	result.Meta["channel"] = channel
	if meta.integrationID != "" {
		result.Meta["integration_id"] = meta.integrationID
	}
	if meta.providerName != "" {
		result.Meta[contracts.MetaKeyProviderName] = meta.providerName
	}
	return result
}
