package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/weprodev/wpd-message-gateway/internal/core/domain"
	"github.com/weprodev/wpd-message-gateway/internal/core/service"
	"github.com/weprodev/wpd-message-gateway/internal/infrastructure/logger"
	"github.com/weprodev/wpd-message-gateway/pkg/contracts"
)

// SendHelper coordinates decoding, dispatching, and logging for message sends.
type SendHelper struct {
	svc *service.GatewayService
}

// NewSendHelper creates a SendHelper.
func NewSendHelper(svc *service.GatewayService) *SendHelper {
	return &SendHelper{svc: svc}
}

// DispatchAndLog handles decoding, invoking the send function, and logging the result.
func (sh *SendHelper) DispatchAndLog(
	c echo.Context,
	channel string,
	workspaceID string,
	apiKeyID string,
	endpoint string,
	dst any,
	send func(ctx context.Context) (*contracts.SendResult, error),
) error {
	start := time.Now()
	ctx := c.Request().Context()

	// Enrich context with gateway-domain attributes for logger
	ctx = logger.WithWorkspace(ctx, workspaceID, apiKeyID)
	ctx = logger.WithChannel(ctx, channel)

	if workspaceID == "" {
		sh.RecordLog(ctx, workspaceID, apiKeyID, channel, c.Request().Method,
			http.StatusUnauthorized, endpoint, start, "missing workspace context", "")
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing workspace context"})
	}

	if err := json.NewDecoder(c.Request().Body).Decode(dst); err != nil {
		sh.RecordLog(ctx, workspaceID, apiKeyID, channel, c.Request().Method,
			http.StatusBadRequest, endpoint, start, "invalid JSON body", "")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid JSON body"})
	}

	result, err := send(ctx)
	if err != nil {
		// Log the full error internally; never echo infrastructure details to the caller.
		slog.ErrorContext(ctx, "send failed", "error", err)
		sh.RecordLog(ctx, workspaceID, apiKeyID, channel, c.Request().Method,
			http.StatusInternalServerError, endpoint, start, err.Error(), providerFromMeta(result))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "send failed"})
	}

	if provider := providerFromMeta(result); provider != "" {
		ctx = logger.WithProvider(ctx, provider)
	}
	slog.InfoContext(ctx, "send ok", "duration_ms", time.Since(start).Milliseconds())
	sh.RecordLog(ctx, workspaceID, apiKeyID, channel, c.Request().Method,
		http.StatusOK, endpoint, start, "", providerFromMeta(result))
	return c.JSON(http.StatusOK, result)
}

// RecordLog persists a MessageRequestLog entry. Failures are logged with slog.
func (sh *SendHelper) RecordLog(
	ctx context.Context,
	workspaceID, apiKeyID, channelType, method string,
	statusCode int,
	endpoint string,
	start time.Time,
	errMsg, providerName string,
) {
	if sh.svc == nil || workspaceID == "" {
		return
	}
	entry := &domain.MessageRequestLog{
		WorkspaceID:  workspaceID,
		APIKeyID:     apiKeyID,
		ChannelType:  channelType,
		HTTPMethod:   method,
		StatusCode:   statusCode,
		Endpoint:     endpoint,
		ProviderName: providerName,
		RequestID:    logger.GetRequestID(ctx),
		DurationMs:   int(time.Since(start).Milliseconds()),
		ErrorMessage: errMsg,
	}
	if err := sh.svc.RecordLog(ctx, entry); err != nil {
		slog.ErrorContext(ctx, "message_request_log insert failed", "error", err)
	}
}

func providerFromMeta(r *contracts.SendResult) string {
	if r == nil || r.Meta == nil {
		return ""
	}
	return r.Meta["provider_name"]
}
