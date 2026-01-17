package devbox

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/weprodev/wpd-message-gateway/contracts"
)

// --- Response Helpers ---

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}

// --- Stats ---

// handleStats returns message counts by type.
// GET /api/v1/stats
func (h *Handler) handleStats(w http.ResponseWriter, r *http.Request) {
	stats := h.store.Stats()
	respondJSON(w, http.StatusOK, stats)
}

// --- Emails ---

// handleGetEmails returns all stored emails.
// GET /api/v1/emails
func (h *Handler) handleGetEmails(w http.ResponseWriter, r *http.Request) {
	emails := h.store.Emails()
	respondJSON(w, http.StatusOK, emails)
}

// handleGetEmailByID returns a single email by ID.
// GET /api/v1/emails/{id}
func (h *Handler) handleGetEmailByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	email := h.store.EmailByID(id)
	if email == nil {
		respondError(w, http.StatusNotFound, "email not found")
		return
	}
	respondJSON(w, http.StatusOK, email)
}

// handleDeleteEmailByID deletes a single email by ID.
// DELETE /api/v1/emails/{id}
func (h *Handler) handleDeleteEmailByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if !h.store.DeleteEmailByID(id) {
		respondError(w, http.StatusNotFound, "email not found")
		return
	}
	h.broadcast("email_deleted", id)
	w.WriteHeader(http.StatusNoContent)
}

// --- SMS ---

// handleGetSMS returns all stored SMS messages.
// GET /api/v1/sms
func (h *Handler) handleGetSMS(w http.ResponseWriter, r *http.Request) {
	sms := h.store.SMS()
	respondJSON(w, http.StatusOK, sms)
}

// handleGetSMSByID returns a single SMS by ID.
// GET /api/v1/sms/{id}
func (h *Handler) handleGetSMSByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	sms := h.store.SMSByID(id)
	if sms == nil {
		respondError(w, http.StatusNotFound, "sms not found")
		return
	}
	respondJSON(w, http.StatusOK, sms)
}

// handleDeleteSMSByID deletes a single SMS by ID.
// DELETE /api/v1/sms/{id}
func (h *Handler) handleDeleteSMSByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if !h.store.DeleteSMSByID(id) {
		respondError(w, http.StatusNotFound, "sms not found")
		return
	}
	h.broadcast("sms_deleted", id)
	w.WriteHeader(http.StatusNoContent)
}

// --- Push Notifications ---

// handleGetPush returns all stored push notifications.
// GET /api/v1/push
func (h *Handler) handleGetPush(w http.ResponseWriter, r *http.Request) {
	pushes := h.store.Pushes()
	respondJSON(w, http.StatusOK, pushes)
}

// handleGetPushByID returns a single push notification by ID.
// GET /api/v1/push/{id}
func (h *Handler) handleGetPushByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	push := h.store.PushByID(id)
	if push == nil {
		respondError(w, http.StatusNotFound, "push notification not found")
		return
	}
	respondJSON(w, http.StatusOK, push)
}

// handleDeletePushByID deletes a single push notification by ID.
// DELETE /api/v1/push/{id}
func (h *Handler) handleDeletePushByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if !h.store.DeletePushByID(id) {
		respondError(w, http.StatusNotFound, "push notification not found")
		return
	}
	h.broadcast("push_deleted", id)
	w.WriteHeader(http.StatusNoContent)
}

// --- Chat Messages ---

// handleGetChat returns all stored chat messages.
// GET /api/v1/chat
func (h *Handler) handleGetChat(w http.ResponseWriter, r *http.Request) {
	chats := h.store.Chats()
	respondJSON(w, http.StatusOK, chats)
}

// handleGetChatByID returns a single chat message by ID.
// GET /api/v1/chat/{id}
func (h *Handler) handleGetChatByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	chat := h.store.ChatByID(id)
	if chat == nil {
		respondError(w, http.StatusNotFound, "chat message not found")
		return
	}
	respondJSON(w, http.StatusOK, chat)
}

// handleDeleteChatByID deletes a single chat message by ID.
// DELETE /api/v1/chat/{id}
func (h *Handler) handleDeleteChatByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if !h.store.DeleteChatByID(id) {
		respondError(w, http.StatusNotFound, "chat message not found")
		return
	}
	h.broadcast("chat_deleted", id)
	w.WriteHeader(http.StatusNoContent)
}

// --- Clear All ---

// handleClearAll removes all stored messages.
// DELETE /api/v1/messages
func (h *Handler) handleClearAll(w http.ResponseWriter, r *http.Request) {
	h.store.Clear()
	h.broadcast("messages_cleared", nil)
	w.WriteHeader(http.StatusNoContent)
}

// --- Internal Ingest Endpoints ---
// These endpoints receive messages from the devbox HTTP provider.

// handleIngestEmail receives an email from external sources.
// POST /api/v1/internal/email
func (h *Handler) handleIngestEmail(w http.ResponseWriter, r *http.Request) {
	var email contracts.Email
	if err := json.NewDecoder(r.Body).Decode(&email); err != nil {
		respondError(w, http.StatusBadRequest, "invalid email payload: "+err.Error())
		return
	}

	// Use the memory provider's email provider to store (also forwards to Mailpit if enabled)
	result, err := h.store.EmailProvider(h.mailpitCfg).Send(r.Context(), &email)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to store email: "+err.Error())
		return
	}

	h.broadcast("email_received", map[string]string{"id": result.ID})
	respondJSON(w, http.StatusCreated, map[string]string{"id": result.ID})
}

// handleIngestSMS receives an SMS from external sources.
// POST /api/v1/internal/sms
func (h *Handler) handleIngestSMS(w http.ResponseWriter, r *http.Request) {
	var sms contracts.SMS
	if err := json.NewDecoder(r.Body).Decode(&sms); err != nil {
		respondError(w, http.StatusBadRequest, "invalid sms payload: "+err.Error())
		return
	}

	result, err := h.store.SMSProvider().Send(r.Context(), &sms)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to store sms: "+err.Error())
		return
	}

	h.broadcast("sms_received", map[string]string{"id": result.ID})
	respondJSON(w, http.StatusCreated, map[string]string{"id": result.ID})
}

// handleIngestPush receives a push notification from external sources.
// POST /api/v1/internal/push
func (h *Handler) handleIngestPush(w http.ResponseWriter, r *http.Request) {
	var push contracts.PushNotification
	if err := json.NewDecoder(r.Body).Decode(&push); err != nil {
		respondError(w, http.StatusBadRequest, "invalid push payload: "+err.Error())
		return
	}

	result, err := h.store.PushProvider().Send(r.Context(), &push)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to store push: "+err.Error())
		return
	}

	h.broadcast("push_received", map[string]string{"id": result.ID})
	respondJSON(w, http.StatusCreated, map[string]string{"id": result.ID})
}

// handleIngestChat receives a chat message from external sources.
// POST /api/v1/internal/chat
func (h *Handler) handleIngestChat(w http.ResponseWriter, r *http.Request) {
	var chat contracts.ChatMessage
	if err := json.NewDecoder(r.Body).Decode(&chat); err != nil {
		respondError(w, http.StatusBadRequest, "invalid chat payload: "+err.Error())
		return
	}

	result, err := h.store.ChatProvider().Send(r.Context(), &chat)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to store chat: "+err.Error())
		return
	}

	h.broadcast("chat_received", map[string]string{"id": result.ID})
	respondJSON(w, http.StatusCreated, map[string]string{"id": result.ID})
}
