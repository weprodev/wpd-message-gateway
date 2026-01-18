package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/weprodev/wpd-message-gateway/internal/core/service"
	"github.com/weprodev/wpd-message-gateway/pkg/contracts"
)

// GatewayHandler handles message sending API endpoints.
type GatewayHandler struct {
	service *service.GatewayService
}

// NewGatewayHandler creates a new gateway handler.
func NewGatewayHandler(svc *service.GatewayService) *GatewayHandler {
	return &GatewayHandler{
		service: svc,
	}
}

// HandleSendEmail handles POST /v1/email
func (h *GatewayHandler) HandleSendEmail(w http.ResponseWriter, r *http.Request) {
	var req contracts.Email
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	result, err := h.service.SendEmail(r.Context(), &req)
	if err != nil {
		log.Printf("Send email error: %v", err)
		http.Error(w, fmt.Sprintf("Failed to send: %v", err), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, result)
}

// HandleSendSMS handles POST /v1/sms
func (h *GatewayHandler) HandleSendSMS(w http.ResponseWriter, r *http.Request) {
	var req contracts.SMS
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	result, err := h.service.SendSMS(r.Context(), &req)
	if err != nil {
		log.Printf("Send SMS error: %v", err)
		http.Error(w, fmt.Sprintf("Failed to send: %v", err), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, result)
}

// HandleSendPush handles POST /v1/push
func (h *GatewayHandler) HandleSendPush(w http.ResponseWriter, r *http.Request) {
	var req contracts.PushNotification
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	result, err := h.service.SendPush(r.Context(), &req)
	if err != nil {
		log.Printf("Send push error: %v", err)
		http.Error(w, fmt.Sprintf("Failed to send: %v", err), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, result)
}

// HandleSendChat handles POST /v1/chat
func (h *GatewayHandler) HandleSendChat(w http.ResponseWriter, r *http.Request) {
	var req contracts.ChatMessage
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	result, err := h.service.SendChat(r.Context(), &req)
	if err != nil {
		log.Printf("Send chat error: %v", err)
		http.Error(w, fmt.Sprintf("Failed to send: %v", err), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, result)
}

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
