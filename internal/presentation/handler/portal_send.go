package handler

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/weprodev/wpd-message-gateway/pkg/contracts"
)

// SendTest dispatches a test message through the gateway for the workspace (portal JWT auth).
func (h *PortalHandler) SendTest(c echo.Context) error {
	if h.gateway == nil || h.helper == nil {
		return echo.NewHTTPError(http.StatusServiceUnavailable, "send test unavailable")
	}

	wid := c.Param("wid")
	channel := c.Param("channel")
	endpoint := "/api/v1/workspaces/" + wid + "/send-test/" + channel

	switch channel {
	case "email":
		var req contracts.Email
		return h.helper.DispatchAndLog(c, channel, wid, "", endpoint, &req, func(sendCtx context.Context) (*contracts.SendResult, error) {
			return h.gateway.SendEmail(sendCtx, wid, &req)
		})
	case "sms":
		var req contracts.SMS
		return h.helper.DispatchAndLog(c, channel, wid, "", endpoint, &req, func(sendCtx context.Context) (*contracts.SendResult, error) {
			return h.gateway.SendSMS(sendCtx, wid, &req)
		})
	case "push":
		var req contracts.PushNotification
		return h.helper.DispatchAndLog(c, channel, wid, "", endpoint, &req, func(sendCtx context.Context) (*contracts.SendResult, error) {
			return h.gateway.SendPush(sendCtx, wid, &req)
		})
	case "chat":
		var req contracts.ChatMessage
		return h.helper.DispatchAndLog(c, channel, wid, "", endpoint, &req, func(sendCtx context.Context) (*contracts.SendResult, error) {
			return h.gateway.SendChat(sendCtx, wid, &req)
		})
	default:
		return echo.NewHTTPError(http.StatusBadRequest, "unsupported channel")
	}
}
