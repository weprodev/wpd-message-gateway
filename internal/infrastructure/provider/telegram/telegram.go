package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/weprodev/wpd-message-gateway/pkg/contracts"
)

const (
	ProviderName   = "telegram"
	defaultBaseURL = "https://api.telegram.org"
	defaultTimeout = 30 * time.Second
)

// Config holds Telegram-specific configuration.
type Config struct {
	APIToken string
	BaseURL  string
}

// Provider implements port.PushSender for Telegram.
type Provider struct {
	config     Config
	baseURL    string
	httpClient *http.Client
}

// New creates a new Telegram push provider.
func New(cfg Config) (*Provider, error) {
	if cfg.APIToken == "" {
		return nil, errors.New("telegram: API token is required")
	}

	baseURL := strings.TrimRight(cfg.BaseURL, "/")
	if baseURL == "" {
		baseURL = defaultBaseURL
	}

	return &Provider{
		config:  cfg,
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
	}, nil
}

// Name returns the provider name.
func (p *Provider) Name() string {
	return ProviderName
}

// Telegram API errors.
var (
	ErrUserBlockedBot = errors.New("telegram: user blocked the bot")
	ErrChatNotFound   = errors.New("telegram: chat not found")
	ErrRateLimited    = errors.New("telegram: too many requests")
	ErrForbidden      = errors.New("telegram: forbidden")
)

// sendMessageRequest is the JSON body for POST /bot<token>/sendMessage.
type sendMessageRequest struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode,omitempty"`
}

// sendMessageResponse is the Telegram Bot API response envelope.
type sendMessageResponse struct {
	OK     bool   `json:"ok"`
	Result struct {
		MessageID int `json:"message_id"`
	} `json:"result"`
	ErrorCode   int    `json:"error_code"`
	Description string `json:"description"`
}

// Send sends a push notification via the Telegram Bot API.
// Each entry in DeviceTokens is treated as a Telegram chat_id.
// Individual failures are logged and skipped — the batch continues.
func (p *Provider) Send(ctx context.Context, notification *contracts.PushNotification) (*contracts.SendResult, error) {
	if len(notification.DeviceTokens) == 0 {
		return nil, errors.New("telegram: no chat_id provided (device_tokens is empty)")
	}

	if notification.Body == "" {
		return nil, errors.New("telegram: message body is required")
	}

	parseMode := notification.Data["parse_mode"]

	var (
		succeeded int
		failed    int
		lastID    string
	)

	for _, chatID := range notification.DeviceTokens {
		id, err := p.sendToChat(ctx, chatID, notification.Body, parseMode)
		if err != nil {
			failed++
			log.Printf("telegram: skipping chat %s: %v", chatID, err)
			continue
		}
		succeeded++
		lastID = id
	}

	if succeeded == 0 {
		return nil, fmt.Errorf("telegram: failed to deliver to all %d chats", failed)
	}

	return &contracts.SendResult{
		ID:         lastID,
		StatusCode: http.StatusOK,
		Message:    fmt.Sprintf("sent to %d/%d chats", succeeded, succeeded+failed),
		Meta: map[string]string{
			"succeeded": fmt.Sprintf("%d", succeeded),
			"failed":    fmt.Sprintf("%d", failed),
		},
	}, nil
}

func (p *Provider) sendToChat(ctx context.Context, chatID, text, parseMode string) (string, error) {
	body := sendMessageRequest{
		ChatID:    chatID,
		Text:      text,
		ParseMode: parseMode,
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/bot%s/sendMessage", p.baseURL, p.config.APIToken)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var apiResp sendMessageResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if !apiResp.OK {
		return "", classifyTelegramError(resp.StatusCode, apiResp.ErrorCode, apiResp.Description)
	}

	return fmt.Sprintf("%d", apiResp.Result.MessageID), nil
}

// classifyTelegramError maps Telegram API error codes to semantic errors.
func classifyTelegramError(httpStatus, apiErrorCode int, description string) error {
	switch {
	case httpStatus == http.StatusForbidden:
		return fmt.Errorf("%w: %s", ErrUserBlockedBot, description)
	case httpStatus == http.StatusNotFound:
		return fmt.Errorf("%w: %s", ErrChatNotFound, description)
	case httpStatus == http.StatusTooManyRequests:
		return fmt.Errorf("%w: %s", ErrRateLimited, description)
	case httpStatus == http.StatusBadRequest:
		return fmt.Errorf("telegram: bad request (code %d): %s", apiErrorCode, description)
	default:
		return fmt.Errorf("telegram: API error (HTTP %d, code %d): %s", httpStatus, apiErrorCode, description)
	}
}
