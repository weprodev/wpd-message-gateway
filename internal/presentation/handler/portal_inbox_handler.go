// Package handler contains HTTP handlers for the web portal and gateway APIs.
// This file implements PortalInboxHandler which serves REST and SSE endpoints for the
// workspace-scoped memory-provider inbox (preview/test mode).
package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"

	"github.com/weprodev/wpd-message-gateway/internal/infrastructure/provider/memory"
	"github.com/weprodev/wpd-message-gateway/pkg/contracts"
)

// PortalInboxHandler serves REST + SSE for the workspace-scoped message inbox (memory provider preview).
// Routes require JWT + workspace membership + workspace API key (see middleware).
type PortalInboxHandler struct {
	store       *memory.Store
	mu          sync.RWMutex
	subscribers map[string]map[chan []byte]bool // workspaceID -> subscriber channels
}

// NewPortalInboxHandler creates a handler for /api/v1/workspaces/:wid/inbox/...
func NewPortalInboxHandler(store *memory.Store) *PortalInboxHandler {
	return &PortalInboxHandler{
		store:       store,
		subscribers: make(map[string]map[chan []byte]bool),
	}
}

// workspaceIDParam extracts the :wid route parameter.
func workspaceIDParam(c echo.Context) string {
	return c.Param("wid")
}

// notYetSupported returns a JSON 404 for inbox endpoints that are not yet workspace-scoped.
// These stubs exist to satisfy the router contract; they will be replaced once the
// respective memory store types support per-workspace isolation.
func notYetSupported(c echo.Context, msgType string) error {
	return c.JSON(http.StatusNotFound, map[string]string{"error": msgType + " not found"})
}

// HandleStats returns message counts for the workspace inbox.
func (h *PortalInboxHandler) HandleStats(c echo.Context) error {
	return c.JSON(http.StatusOK, h.store.StatsForWorkspace(workspaceIDParam(c)))
}

// HandleGetEmails returns stored emails for the workspace.
func (h *PortalInboxHandler) HandleGetEmails(c echo.Context) error {
	emails := h.store.EmailsForWorkspace(workspaceIDParam(c))
	return c.JSON(http.StatusOK, emails)
}

// HandleGetEmailByID returns a single email if it belongs to the workspace.
func (h *PortalInboxHandler) HandleGetEmailByID(c echo.Context) error {
	id := c.Param("id")
	email := h.store.EmailByIDForWorkspace(id, workspaceIDParam(c))
	if email == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "email not found"})
	}
	return c.JSON(http.StatusOK, email)
}

// HandleDeleteEmailByID deletes an email for this workspace.
func (h *PortalInboxHandler) HandleDeleteEmailByID(c echo.Context) error {
	id := c.Param("id")
	if !h.store.DeleteEmailByIDForWorkspace(id, workspaceIDParam(c)) {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "email not found"})
	}
	h.broadcast(workspaceIDParam(c), "email_deleted", id)
	return c.NoContent(http.StatusNoContent)
}

// HandleGetSMS returns an empty list.
// TODO: workspace-scope the SMS memory store, then return real results.
func (h *PortalInboxHandler) HandleGetSMS(c echo.Context) error {
	return c.JSON(http.StatusOK, []any{})
}

// HandleGetSMSByID is not yet workspace-scoped; always returns 404.
func (h *PortalInboxHandler) HandleGetSMSByID(c echo.Context) error {
	return notYetSupported(c, "sms")
}

// HandleDeleteSMSByID is not yet workspace-scoped; always returns 404.
func (h *PortalInboxHandler) HandleDeleteSMSByID(c echo.Context) error {
	return notYetSupported(c, "sms")
}

// HandleGetPush returns an empty list.
// TODO: workspace-scope the push memory store, then return real results.
func (h *PortalInboxHandler) HandleGetPush(c echo.Context) error {
	return c.JSON(http.StatusOK, []any{})
}

// HandleGetPushByID is not yet workspace-scoped; always returns 404.
func (h *PortalInboxHandler) HandleGetPushByID(c echo.Context) error {
	return notYetSupported(c, "push notification")
}

// HandleDeletePushByID is not yet workspace-scoped; always returns 404.
func (h *PortalInboxHandler) HandleDeletePushByID(c echo.Context) error {
	return notYetSupported(c, "push notification")
}

// HandleGetChat returns an empty list.
// TODO: workspace-scope the chat memory store, then return real results.
func (h *PortalInboxHandler) HandleGetChat(c echo.Context) error {
	return c.JSON(http.StatusOK, []any{})
}

// HandleGetChatByID is not yet workspace-scoped; always returns 404.
func (h *PortalInboxHandler) HandleGetChatByID(c echo.Context) error {
	return notYetSupported(c, "chat message")
}

// HandleDeleteChatByID is not yet workspace-scoped; always returns 404.
func (h *PortalInboxHandler) HandleDeleteChatByID(c echo.Context) error {
	return notYetSupported(c, "chat message")
}

// HandleClearAll removes all in-memory messages for this workspace.
func (h *PortalInboxHandler) HandleClearAll(c echo.Context) error {
	h.store.ClearWorkspace(workspaceIDParam(c))
	h.broadcast(workspaceIDParam(c), "messages_cleared", nil)
	return c.NoContent(http.StatusNoContent)
}

// HandleIngestEmail receives an email for a workspace (internal automation).
func (h *PortalInboxHandler) HandleIngestEmail(c echo.Context) error {
	w := workspaceIDParam(c)
	var email contracts.Email
	if err := json.NewDecoder(c.Request().Body).Decode(&email); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid email payload: " + err.Error()})
	}

	emailProvider := memory.NewEmailProviderForWorkspace(h.store, w)
	result, err := emailProvider.Send(c.Request().Context(), &email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to store email: " + err.Error()})
	}

	h.broadcast(w, "email_received", map[string]string{"id": result.ID})
	return c.JSON(http.StatusCreated, map[string]string{"id": result.ID})
}

// HandleIngestSMS stores SMS for a workspace.
// NOTE: The SMS memory store is not yet workspace-scoped; messages are stored globally.
func (h *PortalInboxHandler) HandleIngestSMS(c echo.Context) error {
	var sms contracts.SMS
	if err := json.NewDecoder(c.Request().Body).Decode(&sms); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid sms payload: " + err.Error()})
	}

	smsProvider := memory.NewSMSProvider(h.store)
	result, err := smsProvider.Send(c.Request().Context(), &sms)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to store sms: " + err.Error()})
	}

	h.broadcast(workspaceIDParam(c), "sms_received", map[string]string{"id": result.ID})
	return c.JSON(http.StatusCreated, map[string]string{"id": result.ID})
}

// HandleIngestPush stores a push notification for a workspace.
// NOTE: The push memory store is not yet workspace-scoped.
func (h *PortalInboxHandler) HandleIngestPush(c echo.Context) error {
	var push contracts.PushNotification
	if err := json.NewDecoder(c.Request().Body).Decode(&push); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid push payload: " + err.Error()})
	}

	pushProvider := memory.NewPushProvider(h.store)
	result, err := pushProvider.Send(c.Request().Context(), &push)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to store push: " + err.Error()})
	}

	h.broadcast(workspaceIDParam(c), "push_received", map[string]string{"id": result.ID})
	return c.JSON(http.StatusCreated, map[string]string{"id": result.ID})
}

// HandleIngestChat stores a chat message for a workspace.
// NOTE: The chat memory store is not yet workspace-scoped.
func (h *PortalInboxHandler) HandleIngestChat(c echo.Context) error {
	var chat contracts.ChatMessage
	if err := json.NewDecoder(c.Request().Body).Decode(&chat); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid chat payload: " + err.Error()})
	}

	chatProvider := memory.NewChatProvider(h.store)
	result, err := chatProvider.Send(c.Request().Context(), &chat)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to store chat: " + err.Error()})
	}

	h.broadcast(workspaceIDParam(c), "chat_received", map[string]string{"id": result.ID})
	return c.JSON(http.StatusCreated, map[string]string{"id": result.ID})
}

// HandleSSE streams Server-Sent Events for the workspace inbox (one connection per workspace).
// Clients should reconnect on disconnect. Requires JWT + workspace API key (see middleware).
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
