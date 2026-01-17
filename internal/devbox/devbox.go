// Package devbox provides REST API handlers for the development inbox.
//
// The devbox API allows viewing and managing intercepted messages
// (Email, SMS, Push, Chat) during local development and E2E testing.
//
// Usage:
//
//	store := memory.New()
//	handler := devbox.NewHandler(store)
//
//	// Mount on your router
//	r.Mount("/api/v1", handler.Routes())
package devbox

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/weprodev/wpd-message-gateway/config"
	"github.com/weprodev/wpd-message-gateway/providers/memory"
)

// Handler provides REST API endpoints for the devbox.
type Handler struct {
	store      *memory.Provider
	mailpitCfg config.MailpitConfig
	// SSE subscribers for real-time updates
	subscribers map[chan []byte]bool
}

// NewHandler creates a new devbox API handler.
func NewHandler(store *memory.Provider, mailpitCfg config.MailpitConfig) *Handler {
	return &Handler{
		store:       store,
		mailpitCfg:  mailpitCfg,
		subscribers: make(map[chan []byte]bool),
	}
}

// Routes returns a chi router with all devbox API endpoints.
// This can be mounted on any parent router.
func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()

	// Stats
	r.Get("/stats", h.handleStats)

	// Emails
	r.Get("/emails", h.handleGetEmails)
	r.Get("/emails/{id}", h.handleGetEmailByID)
	r.Delete("/emails/{id}", h.handleDeleteEmailByID)

	// SMS
	r.Get("/sms", h.handleGetSMS)
	r.Get("/sms/{id}", h.handleGetSMSByID)
	r.Delete("/sms/{id}", h.handleDeleteSMSByID)

	// Push
	r.Get("/push", h.handleGetPush)
	r.Get("/push/{id}", h.handleGetPushByID)
	r.Delete("/push/{id}", h.handleDeletePushByID)

	// Chat
	r.Get("/chat", h.handleGetChat)
	r.Get("/chat/{id}", h.handleGetChatByID)
	r.Delete("/chat/{id}", h.handleDeleteChatByID)

	// Clear all
	r.Delete("/messages", h.handleClearAll)

	// SSE
	r.Get("/events", h.handleSSE)

	// Internal ingest endpoints
	r.Route("/internal", func(r chi.Router) {
		r.Post("/email", h.handleIngestEmail)
		r.Post("/sms", h.handleIngestSMS)
		r.Post("/push", h.handleIngestPush)
		r.Post("/chat", h.handleIngestChat)
	})

	return r
}

// RoutesWithMiddleware returns a chi router with middleware pre-configured.
// Use this for standalone deployment or when you want default middleware.
func (h *Handler) RoutesWithMiddleware() chi.Router {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Mount("/", h.Routes())

	return r
}

// Store returns the underlying memory store for testing purposes.
func (h *Handler) Store() *memory.Provider {
	return h.store
}
