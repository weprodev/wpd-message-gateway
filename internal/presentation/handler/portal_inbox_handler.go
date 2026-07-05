// Package handler contains HTTP handlers for the web portal and gateway APIs.
// This file implements PortalInboxHandler which serves REST and SSE endpoints for the
// workspace-scoped memory-provider inbox (preview/test mode).
package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/labstack/echo/v4"

	"github.com/weprodev/wpd-message-gateway/internal/core/port"
	"github.com/weprodev/wpd-message-gateway/pkg/contracts"
)

// PortalInboxHandler serves REST + SSE for the workspace-scoped message inbox (memory provider preview).
// Routes require JWT + workspace membership (see middleware).
type PortalInboxHandler struct {
	reader      port.InboxReader
	writer      port.InboxWriter
	mu          sync.RWMutex
	subscribers map[string]map[chan []byte]bool // workspaceID -> subscriber channels
}

func NewPortalInboxHandler(reader port.InboxReader, writer port.InboxWriter) *PortalInboxHandler {
	return &PortalInboxHandler{
		reader:      reader,
		writer:      writer,
		subscribers: make(map[string]map[chan []byte]bool),
	}
}

func workspaceIDParam(c echo.Context) string {
	return c.Param("wid")
}

func (h *PortalInboxHandler) HandleStats(c echo.Context) error {
	return c.JSON(http.StatusOK, h.reader.StatsForWorkspace(workspaceIDParam(c)))
}

func (h *PortalInboxHandler) HandleGetEmails(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	cursor := c.QueryParam("cursor")
	page := h.reader.ListEmailsForWorkspace(workspaceIDParam(c), limit, cursor)
	return c.JSON(http.StatusOK, page)
}

func (h *PortalInboxHandler) HandleGetEmailByID(c echo.Context) error {
	id := c.Param("id")
	email, ok := h.reader.EmailByIDForWorkspace(id, workspaceIDParam(c))
	if !ok {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "email not found"})
	}
	return c.JSON(http.StatusOK, email)
}

func (h *PortalInboxHandler) HandleDeleteEmailByID(c echo.Context) error {
	id := c.Param("id")
	if !h.reader.DeleteEmailByIDForWorkspace(id, workspaceIDParam(c)) {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "email not found"})
	}
	h.broadcast(workspaceIDParam(c), "email_deleted", id)
	return c.NoContent(http.StatusNoContent)
}

func (h *PortalInboxHandler) HandleGetSMS(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	cursor := c.QueryParam("cursor")
	page := h.reader.ListSMSForWorkspace(workspaceIDParam(c), limit, cursor)
	return c.JSON(http.StatusOK, page)
}

func (h *PortalInboxHandler) HandleGetSMSByID(c echo.Context) error {
	id := c.Param("id")
	sms, ok := h.reader.SMSByIDForWorkspace(id, workspaceIDParam(c))
	if !ok {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "sms not found"})
	}
	return c.JSON(http.StatusOK, sms)
}

func (h *PortalInboxHandler) HandleDeleteSMSByID(c echo.Context) error {
	id := c.Param("id")
	if !h.reader.DeleteSMSByIDForWorkspace(id, workspaceIDParam(c)) {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "sms not found"})
	}
	h.broadcast(workspaceIDParam(c), "sms_deleted", id)
	return c.NoContent(http.StatusNoContent)
}

func (h *PortalInboxHandler) HandleGetPush(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	cursor := c.QueryParam("cursor")
	page := h.reader.ListPushForWorkspace(workspaceIDParam(c), limit, cursor)
	return c.JSON(http.StatusOK, page)
}

func (h *PortalInboxHandler) HandleGetPushByID(c echo.Context) error {
	id := c.Param("id")
	push, ok := h.reader.PushByIDForWorkspace(id, workspaceIDParam(c))
	if !ok {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "push not found"})
	}
	return c.JSON(http.StatusOK, push)
}

func (h *PortalInboxHandler) HandleDeletePushByID(c echo.Context) error {
	id := c.Param("id")
	if !h.reader.DeletePushByIDForWorkspace(id, workspaceIDParam(c)) {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "push not found"})
	}
	h.broadcast(workspaceIDParam(c), "push_deleted", id)
	return c.NoContent(http.StatusNoContent)
}

func (h *PortalInboxHandler) HandleGetChat(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	cursor := c.QueryParam("cursor")
	page := h.reader.ListChatForWorkspace(workspaceIDParam(c), limit, cursor)
	return c.JSON(http.StatusOK, page)
}

func (h *PortalInboxHandler) HandleGetChatByID(c echo.Context) error {
	id := c.Param("id")
	chat, ok := h.reader.ChatByIDForWorkspace(id, workspaceIDParam(c))
	if !ok {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "chat not found"})
	}
	return c.JSON(http.StatusOK, chat)
}

func (h *PortalInboxHandler) HandleDeleteChatByID(c echo.Context) error {
	id := c.Param("id")
	if !h.reader.DeleteChatByIDForWorkspace(id, workspaceIDParam(c)) {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "chat not found"})
	}
	h.broadcast(workspaceIDParam(c), "chat_deleted", id)
	return c.NoContent(http.StatusNoContent)
}

func (h *PortalInboxHandler) HandleClearAll(c echo.Context) error {
	h.reader.ClearWorkspace(workspaceIDParam(c))
	h.broadcast(workspaceIDParam(c), "messages_cleared", nil)
	return c.NoContent(http.StatusNoContent)
}

func (h *PortalInboxHandler) HandleIngestEmail(c echo.Context) error {
	w := workspaceIDParam(c)
	var email contracts.Email
	if err := json.NewDecoder(c.Request().Body).Decode(&email); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid email payload"})
	}

	id, err := h.writer.WriteEmail(c.Request().Context(), w, email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to store email"})
	}

	return c.JSON(http.StatusCreated, map[string]string{"id": id})
}

func (h *PortalInboxHandler) HandleIngestSMS(c echo.Context) error {
	var sms contracts.SMS
	if err := json.NewDecoder(c.Request().Body).Decode(&sms); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid sms payload"})
	}

	id, err := h.writer.WriteSMS(c.Request().Context(), workspaceIDParam(c), sms)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to store sms"})
	}

	return c.JSON(http.StatusCreated, map[string]string{"id": id})
}

func (h *PortalInboxHandler) HandleIngestPush(c echo.Context) error {
	var push contracts.PushNotification
	if err := json.NewDecoder(c.Request().Body).Decode(&push); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid push payload"})
	}

	id, err := h.writer.WritePush(c.Request().Context(), workspaceIDParam(c), push)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to store push"})
	}

	return c.JSON(http.StatusCreated, map[string]string{"id": id})
}

func (h *PortalInboxHandler) HandleIngestChat(c echo.Context) error {
	var chat contracts.ChatMessage
	if err := json.NewDecoder(c.Request().Body).Decode(&chat); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid chat payload"})
	}

	id, err := h.writer.WriteChat(c.Request().Context(), workspaceIDParam(c), chat)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to store chat"})
	}

	return c.JSON(http.StatusCreated, map[string]string{"id": id})
}

// PublishInboxEvent notifies SSE subscribers of an inbox change.
func (h *PortalInboxHandler) PublishInboxEvent(workspaceID, eventType string, data any) {
	h.broadcast(workspaceID, eventType, data)
}

// HandleSSE streams Server-Sent Events for the workspace inbox (one connection per workspace).
func (h *PortalInboxHandler) HandleSSE(c echo.Context) error {
	w := workspaceIDParam(c)
	wr := c.Response().Writer
	r := c.Request()

	wr.Header().Set("Content-Type", "text/event-stream")
	wr.Header().Set("Cache-Control", "no-cache")
	wr.Header().Set("Connection", "keep-alive")
	wr.Header().Set("Access-Control-Allow-Origin", "*")

	events := make(chan []byte, 10)
	h.addSubscriber(w, events)
	defer h.removeSubscriber(w, events)

	flusher, ok := wr.(http.Flusher)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "SSE not supported")
	}

	_, _ = fmt.Fprintf(wr, "event: connected\ndata: {\"status\":\"connected\",\"workspace_id\":\"%s\"}\n\n", w)
	flusher.Flush()

	for {
		select {
		case <-r.Context().Done():
			return nil
		case event, ok := <-events:
			if !ok {
				return nil
			}
			_, _ = fmt.Fprintf(wr, "event: message\ndata: %s\n\n", event)
			flusher.Flush()
		}
	}
}

func (h *PortalInboxHandler) addSubscriber(workspaceID string, ch chan []byte) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.subscribers[workspaceID] == nil {
		h.subscribers[workspaceID] = make(map[chan []byte]bool)
	}
	h.subscribers[workspaceID][ch] = true
}

func (h *PortalInboxHandler) removeSubscriber(workspaceID string, ch chan []byte) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if m := h.subscribers[workspaceID]; m != nil {
		delete(m, ch)
	}
	close(ch)
}

func (h *PortalInboxHandler) broadcast(workspaceID, eventType string, data interface{}) {
	event := map[string]interface{}{
		"type":         eventType,
		"data":         data,
		"workspace_id": workspaceID,
	}

	eventJSON, err := json.Marshal(event)
	if err != nil {
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()
	m := h.subscribers[workspaceID]
	if m == nil {
		return
	}
	for ch := range m {
		select {
		case ch <- eventJSON:
		default:
		}
	}
}
