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

func wid(c echo.Context) string {
	return c.Param("wid")
}

// HandleStats returns message counts for the workspace inbox.
func (h *PortalInboxHandler) HandleStats(c echo.Context) error {
	return c.JSON(http.StatusOK, h.store.StatsForWorkspace(wid(c)))
}

// HandleGetEmails returns stored emails for the workspace.
func (h *PortalInboxHandler) HandleGetEmails(c echo.Context) error {
	emails := h.store.EmailsForWorkspace(wid(c))
	return c.JSON(http.StatusOK, emails)
}

// HandleGetEmailByID returns a single email if it belongs to the workspace.
func (h *PortalInboxHandler) HandleGetEmailByID(c echo.Context) error {
	id := c.Param("id")
	email := h.store.EmailByIDForWorkspace(id, wid(c))
	if email == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "email not found"})
	}
	return c.JSON(http.StatusOK, email)
}

// HandleDeleteEmailByID deletes an email for this workspace.
func (h *PortalInboxHandler) HandleDeleteEmailByID(c echo.Context) error {
	id := c.Param("id")
	if !h.store.DeleteEmailByIDForWorkspace(id, wid(c)) {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "email not found"})
	}
	h.broadcast(wid(c), "email_deleted", id)
	return c.NoContent(http.StatusNoContent)
}

// HandleGetSMS returns SMS for this workspace (empty until SMS rows are workspace-tagged).
func (h *PortalInboxHandler) HandleGetSMS(c echo.Context) error {
	return c.JSON(http.StatusOK, []any{})
}

// HandleGetSMSByID is not available per workspace until SMS store is tagged.
func (h *PortalInboxHandler) HandleGetSMSByID(c echo.Context) error {
	return c.JSON(http.StatusNotFound, map[string]string{"error": "sms not found"})
}

func (h *PortalInboxHandler) HandleDeleteSMSByID(c echo.Context) error {
	return c.JSON(http.StatusNotFound, map[string]string{"error": "sms not found"})
}

func (h *PortalInboxHandler) HandleGetPush(c echo.Context) error {
	return c.JSON(http.StatusOK, []any{})
}

func (h *PortalInboxHandler) HandleGetPushByID(c echo.Context) error {
	return c.JSON(http.StatusNotFound, map[string]string{"error": "push notification not found"})
}

func (h *PortalInboxHandler) HandleDeletePushByID(c echo.Context) error {
	return c.JSON(http.StatusNotFound, map[string]string{"error": "push notification not found"})
}

func (h *PortalInboxHandler) HandleGetChat(c echo.Context) error {
	return c.JSON(http.StatusOK, []any{})
}

func (h *PortalInboxHandler) HandleGetChatByID(c echo.Context) error {
	return c.JSON(http.StatusNotFound, map[string]string{"error": "chat message not found"})
}

func (h *PortalInboxHandler) HandleDeleteChatByID(c echo.Context) error {
	return c.JSON(http.StatusNotFound, map[string]string{"error": "chat message not found"})
}

// HandleClearAll removes in-memory messages for this workspace (emails).
func (h *PortalInboxHandler) HandleClearAll(c echo.Context) error {
	h.store.ClearWorkspace(wid(c))
	h.broadcast(wid(c), "messages_cleared", nil)
	return c.NoContent(http.StatusNoContent)
}

// HandleIngestEmail receives an email for a workspace (internal automation).
func (h *PortalInboxHandler) HandleIngestEmail(c echo.Context) error {
	w := wid(c)
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

// HandleIngestSMS stores SMS scoped to workspace when SMS memory provider supports workspace IDs.
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

	h.broadcast(wid(c), "sms_received", map[string]string{"id": result.ID})
	return c.JSON(http.StatusCreated, map[string]string{"id": result.ID})
}

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

	h.broadcast(wid(c), "push_received", map[string]string{"id": result.ID})
	return c.JSON(http.StatusCreated, map[string]string{"id": result.ID})
}

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

	h.broadcast(wid(c), "chat_received", map[string]string{"id": result.ID})
	return c.JSON(http.StatusCreated, map[string]string{"id": result.ID})
}

// HandleSSE streams events for one workspace only.
func (h *PortalInboxHandler) HandleSSE(c echo.Context) error {
	w := wid(c)
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
