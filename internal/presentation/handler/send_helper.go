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

	ctx = logger.WithWorkspace(ctx, workspaceID, apiKeyID)
	ctx = logger.WithChannel(ctx, channel)

	if workspaceID == "" {
		slog.WarnContext(ctx, "send rejected: missing workspace context", "endpoint", endpoint)
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing workspace context"})
	}

	if err := json.NewDecoder(c.Request().Body).Decode(dst); err != nil {
		slog.WarnContext(ctx, "send rejected: invalid JSON body", "endpoint", endpoint, "error", err)
		sh.recordLog(ctx, workspaceID, apiKeyID, channel, c.Request().Method,
			http.StatusBadRequest, endpoint, start, "invalid JSON body", "")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid JSON body"})
	}

	slog.InfoContext(ctx, "send dispatch starting", "endpoint", endpoint)

	result, err := send(ctx)
	providerName := contracts.ProviderNameFromResult(result)
	if providerName != "" {
		ctx = logger.WithProvider(ctx, providerName)
	}

	if err != nil {
		slog.ErrorContext(ctx, "send failed", "error", err, "endpoint", endpoint, "duration_ms", time.Since(start).Milliseconds())
		sh.recordLog(ctx, workspaceID, apiKeyID, channel, c.Request().Method,
			http.StatusInternalServerError, endpoint, start, err.Error(), providerName)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "send failed"})
	}

	slog.InfoContext(ctx, "send ok", "endpoint", endpoint, "duration_ms", time.Since(start).Milliseconds())
	sh.recordLog(ctx, workspaceID, apiKeyID, channel, c.Request().Method,
		http.StatusOK, endpoint, start, "", providerName)
	return c.JSON(http.StatusOK, result)
}

func (sh *SendHelper) recordLog(
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
	entry := domain.MessageRequestLog{
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
		slog.WarnContext(ctx, "failed to persist message request log", "error", err, "endpoint", endpoint)
	}
}
