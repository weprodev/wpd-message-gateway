package devbox

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// SSE manages Server-Sent Events subscribers.
// This allows the frontend to receive real-time updates when messages change.

// SSE event types
const (
	EventEmailReceived   = "email_received"
	EventEmailDeleted    = "email_deleted"
	EventSMSReceived     = "sms_received"
	EventSMSDeleted      = "sms_deleted"
	EventPushReceived    = "push_received"
	EventPushDeleted     = "push_deleted"
	EventChatReceived    = "chat_received"
	EventChatDeleted     = "chat_deleted"
	EventMessagesCleared = "messages_cleared"
)

// SSEEvent represents an event to be sent to subscribers.
type SSEEvent struct {
	Type string      `json:"type"`
	Data interface{} `json:"data,omitempty"`
}

// handleSSE handles Server-Sent Events connections.
// GET /api/v1/events
func (h *Handler) handleSSE(w http.ResponseWriter, r *http.Request) {
	// Set SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Create subscriber channel
	events := make(chan []byte, 10)
	h.addSubscriber(events)

	// Remove subscriber on disconnect
	defer h.removeSubscriber(events)

	// Flush support
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "SSE not supported", http.StatusInternalServerError)
		return
	}

	// Send initial connected event
	_, _ = fmt.Fprintf(w, "event: connected\ndata: {\"status\":\"connected\"}\n\n")
	flusher.Flush()

	// Listen for events or client disconnect
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

// addSubscriber adds a new SSE subscriber.
func (h *Handler) addSubscriber(ch chan []byte) {
	h.subscribers[ch] = true
}

// removeSubscriber removes an SSE subscriber.
func (h *Handler) removeSubscriber(ch chan []byte) {
	delete(h.subscribers, ch)
	close(ch)
}

// broadcast sends an event to all SSE subscribers.
func (h *Handler) broadcast(eventType string, data interface{}) {
	event := SSEEvent{
		Type: eventType,
		Data: data,
	}

	eventJSON, err := json.Marshal(event)
	if err != nil {
		return
	}

	for ch := range h.subscribers {
		select {
		case ch <- eventJSON:
		default:
			// Channel full, skip (non-blocking)
		}
	}
}

// BroadcastNewEmail broadcasts a new email event to all subscribers.
func (h *Handler) BroadcastNewEmail(id string) {
	h.broadcast(EventEmailReceived, map[string]string{"id": id})
}

// BroadcastNewSMS broadcasts a new SMS event to all subscribers.
func (h *Handler) BroadcastNewSMS(id string) {
	h.broadcast(EventSMSReceived, map[string]string{"id": id})
}

// BroadcastNewPush broadcasts a new push notification event to all subscribers.
func (h *Handler) BroadcastNewPush(id string) {
	h.broadcast(EventPushReceived, map[string]string{"id": id})
}

// BroadcastNewChat broadcasts a new chat message event to all subscribers.
func (h *Handler) BroadcastNewChat(id string) {
	h.broadcast(EventChatReceived, map[string]string{"id": id})
}
