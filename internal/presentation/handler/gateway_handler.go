package handler

import (
	"context"

	"github.com/labstack/echo/v4"

	"github.com/weprodev/wpd-message-gateway/internal/core/service"
	"github.com/weprodev/wpd-message-gateway/internal/presentation/middleware"
	"github.com/weprodev/wpd-message-gateway/pkg/contracts"
)

// GatewayHandler handles the send-message API endpoints (/v1/*).
// Each handler parses the request, delegates to GatewayService, records the
// request log, and returns a uniform JSON response.
type GatewayHandler struct {
	service *service.GatewayService
	helper  *SendHelper
}

// NewGatewayHandler creates a GatewayHandler.
func NewGatewayHandler(svc *service.GatewayService) *GatewayHandler {
	return &GatewayHandler{
		service: svc,
		helper:  NewSendHelper(svc),
	}
}

// HandleSendEmail handles POST /v1/email.
func (h *GatewayHandler) HandleSendEmail(c echo.Context) error {
	var req contracts.Email
	ctx := c.Request().Context()
	workspaceID := middleware.GetWorkspaceID(ctx)
	apiKeyID := middleware.GetAPIKeyID(ctx)
	return h.helper.DispatchAndLog(c, "email", workspaceID, apiKeyID, c.Path(), &req, func(sendCtx context.Context) (*contracts.SendResult, error) {
		return h.service.SendEmail(sendCtx, workspaceID, req)
	})
}

// HandleSendSMS handles POST /v1/sms.
func (h *GatewayHandler) HandleSendSMS(c echo.Context) error {
	var req contracts.SMS
	ctx := c.Request().Context()
	workspaceID := middleware.GetWorkspaceID(ctx)
	apiKeyID := middleware.GetAPIKeyID(ctx)
	return h.helper.DispatchAndLog(c, "sms", workspaceID, apiKeyID, c.Path(), &req, func(sendCtx context.Context) (*contracts.SendResult, error) {
		return h.service.SendSMS(sendCtx, workspaceID, req)
	})
}

// HandleSendPush handles POST /v1/push.
func (h *GatewayHandler) HandleSendPush(c echo.Context) error {
	var req contracts.PushNotification
	ctx := c.Request().Context()
	workspaceID := middleware.GetWorkspaceID(ctx)
	apiKeyID := middleware.GetAPIKeyID(ctx)
	return h.helper.DispatchAndLog(c, "push", workspaceID, apiKeyID, c.Path(), &req, func(sendCtx context.Context) (*contracts.SendResult, error) {
		return h.service.SendPush(sendCtx, workspaceID, req)
	})
}

// HandleSendChat handles POST /v1/chat.
func (h *GatewayHandler) HandleSendChat(c echo.Context) error {
	var req contracts.ChatMessage
	ctx := c.Request().Context()
	workspaceID := middleware.GetWorkspaceID(ctx)
	apiKeyID := middleware.GetAPIKeyID(ctx)
	return h.helper.DispatchAndLog(c, "chat", workspaceID, apiKeyID, c.Path(), &req, func(sendCtx context.Context) (*contracts.SendResult, error) {
		return h.service.SendChat(sendCtx, workspaceID, req)
	})
}
