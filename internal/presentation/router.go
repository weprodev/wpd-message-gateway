package presentation

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/weprodev/wpd-message-gateway/internal/presentation/handler"
)

// Router holds all HTTP handlers and provides route configuration.
type Router struct {
	gatewayHandler *handler.GatewayHandler
	devboxHandler  *handler.DevBoxHandler
}

// NewRouter creates a new router with the given handlers.
func NewRouter(gateway *handler.GatewayHandler, devbox *handler.DevBoxHandler) *Router {
	return &Router{
		gatewayHandler: gateway,
		devboxHandler:  devbox,
	}
}

// Setup creates and configures the chi router with all routes.
func (rt *Router) Setup() chi.Router {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Gateway API - for sending messages
	r.Route("/v1", func(r chi.Router) {
		r.Post("/email", rt.gatewayHandler.HandleSendEmail)
		r.Post("/sms", rt.gatewayHandler.HandleSendSMS)
		r.Post("/push", rt.gatewayHandler.HandleSendPush)
		r.Post("/chat", rt.gatewayHandler.HandleSendChat)
		r.Get("/inbox", rt.gatewayHandler.HandleGetInbox)
		r.Delete("/inbox", rt.gatewayHandler.HandleClearInbox)
	})

	// DevBox API - for viewing intercepted messages
	if rt.devboxHandler != nil {
		r.Mount("/api/v1", rt.devboxRoutes())
	}

	return r
}

// devboxRoutes returns a chi router with all devbox API endpoints.
func (rt *Router) devboxRoutes() chi.Router {
	r := chi.NewRouter()

	// Stats
	r.Get("/stats", rt.devboxHandler.HandleStats)

	// Emails
	r.Get("/emails", rt.devboxHandler.HandleGetEmails)
	r.Get("/emails/{id}", rt.devboxHandler.HandleGetEmailByID)
	r.Delete("/emails/{id}", rt.devboxHandler.HandleDeleteEmailByID)

	// SMS
	r.Get("/sms", rt.devboxHandler.HandleGetSMS)
	r.Get("/sms/{id}", rt.devboxHandler.HandleGetSMSByID)
	r.Delete("/sms/{id}", rt.devboxHandler.HandleDeleteSMSByID)

	// Push
	r.Get("/push", rt.devboxHandler.HandleGetPush)
	r.Get("/push/{id}", rt.devboxHandler.HandleGetPushByID)
	r.Delete("/push/{id}", rt.devboxHandler.HandleDeletePushByID)

	// Chat
	r.Get("/chat", rt.devboxHandler.HandleGetChat)
	r.Get("/chat/{id}", rt.devboxHandler.HandleGetChatByID)
	r.Delete("/chat/{id}", rt.devboxHandler.HandleDeleteChatByID)

	// Clear all
	r.Delete("/messages", rt.devboxHandler.HandleClearAll)

	// SSE
	r.Get("/events", rt.devboxHandler.HandleSSE)

	// Internal ingest endpoints
	r.Route("/internal", func(r chi.Router) {
		r.Post("/email", rt.devboxHandler.HandleIngestEmail)
		r.Post("/sms", rt.devboxHandler.HandleIngestSMS)
		r.Post("/push", rt.devboxHandler.HandleIngestPush)
		r.Post("/chat", rt.devboxHandler.HandleIngestChat)
	})

	return r
}
