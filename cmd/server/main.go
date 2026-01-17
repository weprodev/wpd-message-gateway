package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/weprodev/wpd-message-gateway/config"
	"github.com/weprodev/wpd-message-gateway/contracts"
	"github.com/weprodev/wpd-message-gateway/internal/devbox"
	"github.com/weprodev/wpd-message-gateway/manager"
	"github.com/weprodev/wpd-message-gateway/providers/memory"
)

type Server struct {
	mgr *manager.Manager
}

func main() {
	// Load configuration (defaults to configs/local.yml or CONFIG_PATH)
	configPath := os.Getenv("CONFIG_PATH")
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Validate required configuration
	if err := validateConfig(cfg); err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	log.Printf("Loaded Configuration:")
	log.Printf("- Email Provider: %s (Default)", cfg.DefaultEmailProvider())
	log.Printf("- SMS Provider:   %s (Default)", cfg.DefaultSMSProvider())
	log.Printf("- Push Provider:  %s (Default)", cfg.DefaultPushProvider())
	log.Printf("- Chat Provider:  %s (Default)", cfg.DefaultChatProvider())

	mgr, err := manager.New(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize manager: %v", err)
	}

	server := &Server{mgr: mgr}

	r := chi.NewRouter()
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
		r.Post("/email", server.handleSendEmail)
		r.Post("/sms", server.handleSendSMS)
		r.Post("/push", server.handleSendPush)
		r.Post("/chat", server.handleSendChat)
		r.Get("/inbox", server.handleGetInbox)
		r.Delete("/inbox", server.handleClearInbox)
	})

	// DevBox API - for viewing intercepted messages
	if store := mgr.GetMemoryStore(); store != nil {
		devboxHandler := devbox.NewHandler(store, cfg.Mailpit)
		r.Mount("/api/v1", devboxHandler.Routes())
		log.Println("DevBox API enabled at /api/v1")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "10101"
	}

	log.Printf("Gateway server listening on :%s", port)
	log.Printf("Default Email Provider: %s", cfg.DefaultEmailProvider())

	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

// POST /v1/email
func (s *Server) handleSendEmail(w http.ResponseWriter, r *http.Request) {
	var req contracts.Email
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	result, err := s.mgr.SendEmail(r.Context(), &req)
	if err != nil {
		log.Printf("Send error: %v", err)
		http.Error(w, fmt.Sprintf("Failed to send: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(result)
}

// POST /v1/sms
func (s *Server) handleSendSMS(w http.ResponseWriter, r *http.Request) {
	var req contracts.SMS
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	result, err := s.mgr.SendSMS(r.Context(), &req)
	if err != nil {
		log.Printf("Send error: %v", err)
		http.Error(w, fmt.Sprintf("Failed to send: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(result)
}

// POST /v1/push
func (s *Server) handleSendPush(w http.ResponseWriter, r *http.Request) {
	var req contracts.PushNotification
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	result, err := s.mgr.SendPush(r.Context(), &req)
	if err != nil {
		log.Printf("Send error: %v", err)
		http.Error(w, fmt.Sprintf("Failed to send: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(result)
}

// POST /v1/chat
func (s *Server) handleSendChat(w http.ResponseWriter, r *http.Request) {
	var req contracts.ChatMessage
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	result, err := s.mgr.SendChat(r.Context(), &req)
	if err != nil {
		log.Printf("Send error: %v", err)
		http.Error(w, fmt.Sprintf("Failed to send: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(result)
}

// GET /v1/inbox
func (s *Server) handleGetInbox(w http.ResponseWriter, r *http.Request) {
	provider, err := s.mgr.EmailProvider("memory")
	if err != nil {
		http.Error(w, "Memory provider not active", http.StatusNotFound)
		return
	}

	memProvider, ok := provider.(*memory.EmailProvider)
	if !ok {
		http.Error(w, "Invalid provider type", http.StatusInternalServerError)
		return
	}

	messages := memProvider.Store().Emails()
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(messages)
}

// DELETE /v1/inbox
func (s *Server) handleClearInbox(w http.ResponseWriter, r *http.Request) {
	provider, err := s.mgr.EmailProvider("memory")
	if err != nil {
		http.Error(w, "Memory provider not active", http.StatusNotFound)
		return
	}

	memProvider, ok := provider.(*memory.EmailProvider)
	if !ok {
		http.Error(w, "Invalid provider type", http.StatusInternalServerError)
		return
	}

	memProvider.Store().Clear()
	w.WriteHeader(http.StatusNoContent)
}
