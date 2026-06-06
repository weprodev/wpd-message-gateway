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

	"github.com/weprodev/wpd-message-gateway/pkg/contracts"
)

const (
	ProviderName   = "telegram"
	defaultBaseURL = "https://api.telegram.org"
)

// Config holds Telegram-specific configuration.
type Config struct {
	APIToken 	  string
	BaseURL  	  string
	KnownCommands []string
}

// Provider implements port.PushSender for Telegram.
type Provider struct {
	config     Config
	baseURL    string
	httpClient *http.Client
	offset     int64
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
		httpClient: &http.Client{},
	}, nil
}

// Name returns the provider name.
func (p *Provider) Name() string {
	return ProviderName
}


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

// getUpdatesResponse is the response from GET /bot<token>/getUpdates.
type getUpdatesResponse struct {
	OK     bool     `json:"ok"`
	Result []Update `json:"result"`
}

// Update represents a Telegram update from getUpdates.
type Update struct {
	UpdateID int64    `json:"update_id"`
	Message  *Message `json:"message,omitempty"`
}

// Message represents a Telegram message.
type Message struct {
	MessageID int    `json:"message_id"`
	Chat      Chat   `json:"chat"`
	Text      string `json:"text,omitempty"`
}

// Chat represents a Telegram chat.
type Chat struct {
	ID int64 `json:"id"`
}

// sendToChat sends a plain-text message to a Telegram chat.
func (p *Provider) sendToChat(ctx context.Context, chatID, text, parseMode string) (string, error) {
	body := sendMessageRequest{
		ChatID:    chatID,
		Text:      text,
		ParseMode: parseMode,
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/bot%s/sendMessage", p.baseURL, p.config.APIToken)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	var apiResp sendMessageResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	if !apiResp.OK {
		switch {
		case resp.StatusCode == http.StatusForbidden:
			return "", fmt.Errorf("telegram: user blocked the bot: %s", apiResp.Description)
		case resp.StatusCode == http.StatusNotFound:
			return "", fmt.Errorf("telegram: chat not found: %s", apiResp.Description)
		case resp.StatusCode == http.StatusTooManyRequests:
			return "", fmt.Errorf("telegram: too many requests: %s", apiResp.Description)
		case resp.StatusCode == http.StatusBadRequest:
			return "", fmt.Errorf("telegram: bad request (code %d): %s", apiResp.ErrorCode, apiResp.Description)
		default:
			return "", fmt.Errorf("telegram: API error (HTTP %d, code %d): %s", resp.StatusCode, apiResp.ErrorCode, apiResp.Description)
		}
	}

	return fmt.Sprintf("%d", apiResp.Result.MessageID), nil
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

// UnsupportedMessageResponse sends an unsupported-message reply if the
// message is not a recognised command.
func (p *Provider) UnsupportedMessageResponse(ctx context.Context, msg *Message) bool {
	if msg == nil || msg.Text == "" {
		return false
	}

	if strings.HasPrefix(msg.Text, "/") {
		return false
	}

	for _, cmd := range p.config.KnownCommands {
		if msg.Text == cmd {
			return false
		}
	}

	chatID := fmt.Sprintf("%d", msg.Chat.ID)
	if _, err := p.sendToChat(ctx, chatID, "⚠️ Unsupported message! Please use the menu or available commands", ""); err != nil {
		log.Printf("telegram: failed to reply to chat %s: %v", chatID, err)
		return false
	}
	return true
}

// StartPolling begins long-polling the Telegram Bot API for incoming updates.
// It blocks until ctx is cancelled. Run it in a goroutine:
//
//	go provider.StartPolling(ctx)
//
// For every message that is not a recognised command, the bot replies with an
// unsupported-message warning.
func (p *Provider) StartPolling(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		updates, err := p.fetchUpdates(ctx)
		if err != nil {
			log.Printf("telegram: polling error: %v", err)
			continue
		}

		for _, upd := range updates {
			p.UnsupportedMessageResponse(ctx, upd.Message)
			p.offset = upd.UpdateID + 1
		}
	}
}

// fetchUpdates calls getUpdates with the current offset and a long-polling timeout.
func (p *Provider) fetchUpdates(ctx context.Context) ([]Update, error) {
	url := fmt.Sprintf(
		"%s/bot%s/getUpdates?offset=%d",
		p.baseURL, p.config.APIToken, p.offset,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create getUpdates request: %w", err)
	}

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("getUpdates HTTP error: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read getUpdates response: %w", err)
	}

	var apiResp getUpdatesResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, fmt.Errorf("decode getUpdates response: %w", err)
	}

	if !apiResp.OK {
		return nil, errors.New("telegram: getUpdates returned not-ok")
	}

	return apiResp.Result, nil
}

