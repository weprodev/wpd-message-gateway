package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"

	"github.com/weprodev/wpd-message-gateway/internal/infrastructure/provider/memory"
	"github.com/weprodev/wpd-message-gateway/pkg/contracts"
)

// DevBoxHandler provides REST API endpoints for the development inbox.
type DevBoxHandler struct {
	store       *memory.Store
	mailpitCfg  memory.MailpitConfig
	mu          sync.RWMutex // Protects subscribers map
	subscribers map[chan []byte]bool
}

// NewDevBoxHandler creates a new devbox handler.
func NewDevBoxHandler(store *memory.Store, mailpitCfg memory.MailpitConfig) *DevBoxHandler {
	return &DevBoxHandler{
		store:       store,
		mailpitCfg:  mailpitCfg,
		subscribers: make(map[chan []byte]bool),
	}
}

// Store returns the underlying memory store.
func (h *DevBoxHandler) Store() *memory.Store {
	return h.store
}

// --- Stats ---

// HandleStats returns message counts by type.
func (h *DevBoxHandler) HandleStats(w http.ResponseWriter, r *http.Request) {
	stats := h.store.Stats()
	respondJSON(w, http.StatusOK, stats)
}

// --- Emails ---

// HandleGetEmails returns all stored emails.
func (h *DevBoxHandler) HandleGetEmails(w http.ResponseWriter, r *http.Request) {
	emails := h.store.Emails()
	respondJSON(w, http.StatusOK, emails)
}

// HandleGetEmailByID returns a single email by ID.
func (h *DevBoxHandler) HandleGetEmailByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	email := h.store.EmailByID(id)
	if email == nil {
		respondError(w, http.StatusNotFound, "email not found")
		return
	}
	respondJSON(w, http.StatusOK, email)
}

// HandleDeleteEmailByID deletes a single email by ID.
func (h *DevBoxHandler) HandleDeleteEmailByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if !h.store.DeleteEmailByID(id) {
		respondError(w, http.StatusNotFound, "email not found")
		return
	}
	h.broadcast("email_deleted", id)
	w.WriteHeader(http.StatusNoContent)
}

// --- SMS ---

// HandleGetSMS returns all stored SMS messages.
func (h *DevBoxHandler) HandleGetSMS(w http.ResponseWriter, r *http.Request) {
	sms := h.store.AllSMS()
	respondJSON(w, http.StatusOK, sms)
}

// HandleGetSMSByID returns a single SMS by ID.
func (h *DevBoxHandler) HandleGetSMSByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	sms := h.store.SMSByID(id)
	if sms == nil {
		respondError(w, http.StatusNotFound, "sms not found")
		return
	}
	respondJSON(w, http.StatusOK, sms)
}

// HandleDeleteSMSByID deletes a single SMS by ID.
func (h *DevBoxHandler) HandleDeleteSMSByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if !h.store.DeleteSMSByID(id) {
		respondError(w, http.StatusNotFound, "sms not found")
		return
	}
	h.broadcast("sms_deleted", id)
	w.WriteHeader(http.StatusNoContent)
}

// --- Push Notifications ---

// HandleGetPush returns all stored push notifications.
func (h *DevBoxHandler) HandleGetPush(w http.ResponseWriter, r *http.Request) {
	pushes := h.store.Pushes()
	respondJSON(w, http.StatusOK, pushes)
}

// HandleGetPushByID returns a single push notification by ID.
func (h *DevBoxHandler) HandleGetPushByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	push := h.store.PushByID(id)
	if push == nil {
		respondError(w, http.StatusNotFound, "push notification not found")
		return
	}
	respondJSON(w, http.StatusOK, push)
}

// HandleDeletePushByID deletes a single push notification by ID.
func (h *DevBoxHandler) HandleDeletePushByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if !h.store.DeletePushByID(id) {
		respondError(w, http.StatusNotFound, "push notification not found")
		return
	}
	h.broadcast("push_deleted", id)
	w.WriteHeader(http.StatusNoContent)
}

// --- Chat Messages ---

// HandleGetChat returns all stored chat messages.
func (h *DevBoxHandler) HandleGetChat(w http.ResponseWriter, r *http.Request) {
	chats := h.store.Chats()
	respondJSON(w, http.StatusOK, chats)
}

// HandleGetChatByID returns a single chat message by ID.
func (h *DevBoxHandler) HandleGetChatByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	chat := h.store.ChatByID(id)
	if chat == nil {
		respondError(w, http.StatusNotFound, "chat message not found")
		return
	}
	respondJSON(w, http.StatusOK, chat)
}

// HandleDeleteChatByID deletes a single chat message by ID.
func (h *DevBoxHandler) HandleDeleteChatByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if !h.store.DeleteChatByID(id) {
		respondError(w, http.StatusNotFound, "chat message not found")
		return
	}
	h.broadcast("chat_deleted", id)
	w.WriteHeader(http.StatusNoContent)
}

// --- Clear All ---

// HandleClearAll removes all stored messages.
func (h *DevBoxHandler) HandleClearAll(w http.ResponseWriter, r *http.Request) {
	h.store.Clear()
	h.broadcast("messages_cleared", nil)
	w.WriteHeader(http.StatusNoContent)
}

// --- Internal Ingest Endpoints ---

// HandleIngestEmail receives an email from external sources.
func (h *DevBoxHandler) HandleIngestEmail(w http.ResponseWriter, r *http.Request) {
	var email contracts.Email
	if err := json.NewDecoder(r.Body).Decode(&email); err != nil {
		respondError(w, http.StatusBadRequest, "invalid email payload: "+err.Error())
		return
	}

	emailProvider := memory.NewEmailProvider(h.store, h.mailpitCfg)
	result, err := emailProvider.Send(r.Context(), &email)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to store email: "+err.Error())
		return
	}

	h.broadcast("email_received", map[string]string{"id": result.ID})
	respondJSON(w, http.StatusCreated, map[string]string{"id": result.ID})
}

// HandleIngestSMS receives an SMS from external sources.
func (h *DevBoxHandler) HandleIngestSMS(w http.ResponseWriter, r *http.Request) {
	var sms contracts.SMS
	if err := json.NewDecoder(r.Body).Decode(&sms); err != nil {
		respondError(w, http.StatusBadRequest, "invalid sms payload: "+err.Error())
		return
	}

	smsProvider := memory.NewSMSProvider(h.store)
	result, err := smsProvider.Send(r.Context(), &sms)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to store sms: "+err.Error())
		return
	}

	h.broadcast("sms_received", map[string]string{"id": result.ID})
	respondJSON(w, http.StatusCreated, map[string]string{"id": result.ID})
}

// HandleIngestPush receives a push notification from external sources.
func (h *DevBoxHandler) HandleIngestPush(w http.ResponseWriter, r *http.Request) {
	var push contracts.PushNotification
	if err := json.NewDecoder(r.Body).Decode(&push); err != nil {
		respondError(w, http.StatusBadRequest, "invalid push payload: "+err.Error())
		return
	}

	pushProvider := memory.NewPushProvider(h.store)
	result, err := pushProvider.Send(r.Context(), &push)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to store push: "+err.Error())
		return
	}

	h.broadcast("push_received", map[string]string{"id": result.ID})
	respondJSON(w, http.StatusCreated, map[string]string{"id": result.ID})
}

// HandleIngestChat receives a chat message from external sources.
func (h *DevBoxHandler) HandleIngestChat(w http.ResponseWriter, r *http.Request) {
	var chat contracts.ChatMessage
	if err := json.NewDecoder(r.Body).Decode(&chat); err != nil {
		respondError(w, http.StatusBadRequest, "invalid chat payload: "+err.Error())
		return
	}

	chatProvider := memory.NewChatProvider(h.store)
	result, err := chatProvider.Send(r.Context(), &chat)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to store chat: "+err.Error())
		return
	}

	h.broadcast("chat_received", map[string]string{"id": result.ID})
	respondJSON(w, http.StatusCreated, map[string]string{"id": result.ID})
}

// --- SSE (Server-Sent Events) ---

// HandleSSE handles Server-Sent Events connections.
func (h *DevBoxHandler) HandleSSE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	events := make(chan []byte, 10)
	h.addSubscriber(events)
	defer h.removeSubscriber(events)

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "SSE not supported", http.StatusInternalServerError)
		return
	}

	_, _ = fmt.Fprintf(w, "event: connected\ndata: {\"status\":\"connected\"}\n\n")
	flusher.Flush()

	for {
		select {
		case <-r.Context().Done():
			return
		case event, ok := <-events:
			if !ok {
				return
			}
			_, _ = fmt.Fprintf(w, "event: message\ndata: %s\n\n", event)
			flusher.Flush()
		}
	}
}

func (h *DevBoxHandler) addSubscriber(ch chan []byte) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.subscribers[ch] = true
}

func (h *DevBoxHandler) removeSubscriber(ch chan []byte) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.subscribers, ch)
	close(ch)
}

func (h *DevBoxHandler) broadcast(eventType string, data interface{}) {
	event := map[string]interface{}{
		"type": eventType,
		"data": data,
	}

	eventJSON, err := json.Marshal(event)
	if err != nil {
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	for ch := range h.subscribers {
		select {
		case ch <- eventJSON:
		default:
		}
	}
}

// BroadcastNewEmail broadcasts a new email event.
func (h *DevBoxHandler) BroadcastNewEmail(id string) {
	h.broadcast("email_received", map[string]string{"id": id})
}

// BroadcastNewSMS broadcasts a new SMS event.
func (h *DevBoxHandler) BroadcastNewSMS(id string) {
	h.broadcast("sms_received", map[string]string{"id": id})
}

// BroadcastNewPush broadcasts a new push notification event.
func (h *DevBoxHandler) BroadcastNewPush(id string) {
	h.broadcast("push_received", map[string]string{"id": id})
}

// BroadcastNewChat broadcasts a new chat message event.
func (h *DevBoxHandler) BroadcastNewChat(id string) {
	h.broadcast("chat_received", map[string]string{"id": id})
}
