package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/weprodev/wpd-message-gateway/internal/core/domain"
	"github.com/weprodev/wpd-message-gateway/internal/core/port"
	"github.com/weprodev/wpd-message-gateway/internal/core/service"
	"github.com/weprodev/wpd-message-gateway/internal/infrastructure/logger"
	"github.com/weprodev/wpd-message-gateway/internal/presentation/middleware"
	"github.com/weprodev/wpd-message-gateway/pkg/contracts"
)

// GatewayHandler handles the send-message API endpoints (/v1/*).
// Each handler parses the request, delegates to GatewayService, records the
// request log, and returns a uniform JSON response.
type GatewayHandler struct {
	service *service.GatewayService
	logs    port.MessageRequestLogRepository
}

// NewGatewayHandler creates a GatewayHandler.
func NewGatewayHandler(svc *service.GatewayService, logs port.MessageRequestLogRepository) *GatewayHandler {
	return &GatewayHandler{service: svc, logs: logs}
}

// HandleSendEmail handles POST /v1/email.
func (h *GatewayHandler) HandleSendEmail(c echo.Context) error {
	var req contracts.Email
	return h.handleSend(c, "email", &req, func(ctx context.Context, wsID string) (*contracts.SendResult, error) {
		return h.service.SendEmail(ctx, wsID, &req)
	})
}

// HandleSendSMS handles POST /v1/sms.
func (h *GatewayHandler) HandleSendSMS(c echo.Context) error {
	var req contracts.SMS
	return h.handleSend(c, "sms", &req, func(ctx context.Context, wsID string) (*contracts.SendResult, error) {
		return h.service.SendSMS(ctx, wsID, &req)
	})
}

// HandleSendPush handles POST /v1/push.
func (h *GatewayHandler) HandleSendPush(c echo.Context) error {
	var req contracts.PushNotification
	return h.handleSend(c, "push", &req, func(ctx context.Context, wsID string) (*contracts.SendResult, error) {
		return h.service.SendPush(ctx, wsID, &req)
	})
}

// HandleSendChat handles POST /v1/chat.
func (h *GatewayHandler) HandleSendChat(c echo.Context) error {
	var req contracts.ChatMessage
	return h.handleSend(c, "chat", &req, func(ctx context.Context, wsID string) (*contracts.SendResult, error) {
		return h.service.SendChat(ctx, wsID, &req)
	})
}

// HandleSendOTP handles POST /v1/otp.
func (h *GatewayHandler) HandleSendOTP(c echo.Context) error {
	var req contracts.OTP
	return h.handleSend(c, "otp", &req, func(ctx context.Context, wsID string) (*contracts.SendResult, error) {
		return h.service.SendOTP(ctx, wsID, &req)
	})
}

// HandleRevokeOTP handles POST /v1/otp/revoke.
func (h *GatewayHandler) HandleRevokeOTP(c echo.Context) error {
	start := time.Now()
	ctx := c.Request().Context()
	workspaceID := middleware.GetWorkspaceID(ctx)
	apiKeyID := middleware.GetAPIKeyID(ctx)

	ctx = logger.WithWorkspace(ctx, workspaceID, apiKeyID)
	ctx = logger.WithChannel(ctx, "otp")

	if workspaceID == "" {
		h.recordLog(ctx, workspaceID, apiKeyID, "otp", c.Request().Method,
			http.StatusUnauthorized, c.Path(), start, "missing workspace context", "")
		return c.JSON(http.StatusUnauthorized, errorBody("missing workspace context"))
	}

	var req struct {
		RequestID string `json:"request_id"`
	}
	if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
		h.recordLog(ctx, workspaceID, apiKeyID, "otp", c.Request().Method,
			http.StatusBadRequest, c.Path(), start, "invalid JSON body", "")
		return c.JSON(http.StatusBadRequest, errorBody("invalid JSON body"))
	}

	if req.RequestID == "" {
		h.recordLog(ctx, workspaceID, apiKeyID, "otp", c.Request().Method,
			http.StatusBadRequest, c.Path(), start, "missing request_id", "")
		return c.JSON(http.StatusBadRequest, errorBody("missing request_id"))
	}

	result, err := h.service.RevokeOTP(ctx, workspaceID, req.RequestID)
	if err != nil {
		slog.ErrorContext(ctx, "revoke OTP failed",
			append(logger.Attrs(ctx), "error", err, "request_id", req.RequestID)...)
		h.recordLog(ctx, workspaceID, apiKeyID, "otp", c.Request().Method,
			http.StatusInternalServerError, c.Path(), start, err.Error(), "")
		return c.JSON(http.StatusInternalServerError, errorBody(err.Error()))
	}

	slog.InfoContext(ctx, "revoke OTP ok",
		append(logger.Attrs(ctx), "duration_ms", time.Since(start).Milliseconds(), "request_id", req.RequestID, "status_code", result.StatusCode)...)
	h.recordLog(ctx, workspaceID, apiKeyID, "otp", c.Request().Method,
		http.StatusOK, c.Path(), start, "", "")
	return c.JSON(http.StatusOK, result)
}

// HandleCheckOTPStatus handles GET /v1/otp/status/:requestID.
func (h *GatewayHandler) HandleCheckOTPStatus(c echo.Context) error {
	start := time.Now()
	ctx := c.Request().Context()
	workspaceID := middleware.GetWorkspaceID(ctx)
	apiKeyID := middleware.GetAPIKeyID(ctx)
	requestID := c.Param("requestID")

	ctx = logger.WithWorkspace(ctx, workspaceID, apiKeyID)
	ctx = logger.WithChannel(ctx, "otp")

	if workspaceID == "" {
		h.recordLog(ctx, workspaceID, apiKeyID, "otp", c.Request().Method,
			http.StatusUnauthorized, c.Path(), start, "missing workspace context", "")
		return c.JSON(http.StatusUnauthorized, errorBody("missing workspace context"))
	}

	if requestID == "" {
		h.recordLog(ctx, workspaceID, apiKeyID, "otp", c.Request().Method,
			http.StatusBadRequest, c.Path(), start, "missing request_id", "")
		return c.JSON(http.StatusBadRequest, errorBody("missing request_id"))
	}

	status, err := h.service.CheckOTPStatus(ctx, workspaceID, requestID)
	if err != nil {
		slog.ErrorContext(ctx, "check OTP status failed",
			append(logger.Attrs(ctx), "error", err, "request_id", requestID)...)
		h.recordLog(ctx, workspaceID, apiKeyID, "otp", c.Request().Method,
			http.StatusInternalServerError, c.Path(), start, err.Error(), "")
		return c.JSON(http.StatusInternalServerError, errorBody(err.Error()))
	}

	slog.InfoContext(ctx, "check OTP status ok",
		append(logger.Attrs(ctx), "duration_ms", time.Since(start).Milliseconds(), "request_id", requestID, "delivery_status", deliveryStatusString(status.DeliveryStatus))...)
	h.recordLog(ctx, workspaceID, apiKeyID, "otp", c.Request().Method,
		http.StatusOK, c.Path(), start, "", "")
	return c.JSON(http.StatusOK, status)
}

// handleSend is the single DRY dispatcher for all channels.
// It decodes the JSON body into dst, calls send, records the audit log, and
// returns the appropriate HTTP response.
func (h *GatewayHandler) handleSend(
	c echo.Context,
	channel string,
	dst any,
	send func(ctx context.Context, workspaceID string) (*contracts.SendResult, error),
) error {
	start := time.Now()
	ctx := c.Request().Context()
	workspaceID := middleware.GetWorkspaceID(ctx)
	apiKeyID := middleware.GetAPIKeyID(ctx)

	// Enrich context with gateway-domain attributes so every slog call in this
	// request automatically carries workspace_id, api_key_id, and channel.
	ctx = logger.WithWorkspace(ctx, workspaceID, apiKeyID)
	ctx = logger.WithChannel(ctx, channel)

	if workspaceID == "" {
		h.recordLog(ctx, workspaceID, apiKeyID, channel, c.Request().Method,
			http.StatusUnauthorized, c.Path(), start, "missing workspace context", "")
		return c.JSON(http.StatusUnauthorized, errorBody("missing workspace context"))
	}

	if err := json.NewDecoder(c.Request().Body).Decode(dst); err != nil {
		h.recordLog(ctx, workspaceID, apiKeyID, channel, c.Request().Method,
			http.StatusBadRequest, c.Path(), start, "invalid JSON body", "")
		return c.JSON(http.StatusBadRequest, errorBody("invalid JSON body"))
	}

	result, err := send(ctx, workspaceID)
	if err != nil {
		slog.ErrorContext(ctx, "send failed",
			append(logger.Attrs(ctx), "error", err)...)
		h.recordLog(ctx, workspaceID, apiKeyID, channel, c.Request().Method,
			http.StatusInternalServerError, c.Path(), start, err.Error(), providerFromResult(result))
		return c.JSON(http.StatusInternalServerError, errorBody(err.Error()))
	}

	slog.InfoContext(ctx, "send ok",
		append(logger.Attrs(ctx), "duration_ms", time.Since(start).Milliseconds())...)
	h.recordLog(ctx, workspaceID, apiKeyID, channel, c.Request().Method,
		http.StatusOK, c.Path(), start, "", providerFromResult(result))
	return c.JSON(http.StatusOK, result)
}

// recordLog persists a MessageRequestLog entry. Failures are non-fatal and logged.
func (h *GatewayHandler) recordLog(
	ctx context.Context,
	workspaceID, apiKeyID, channelType, method string,
	statusCode int,
	endpoint string,
	start time.Time,
	errMsg, providerName string,
) {
	if h.logs == nil || workspaceID == "" {
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
		DurationMs:   int(time.Since(start).Milliseconds()),
		ErrorMessage: errMsg,
	}
	if err := h.logs.Create(ctx, entry); err != nil {
		slog.ErrorContext(ctx, "message_request_log insert failed",
			append(logger.Attrs(ctx), "error", err)...)
	}
}

func errorBody(msg string) map[string]string { return map[string]string{"error": msg} }

func deliveryStatusString(ds *contracts.DeliveryStatus) string {
	if ds == nil {
		return "unknown"
	}
	return ds.Status
}

func providerFromResult(r *contracts.SendResult) string {
	if r == nil || r.Meta == nil {
		return ""
	}
	return r.Meta["provider_name"]
}
