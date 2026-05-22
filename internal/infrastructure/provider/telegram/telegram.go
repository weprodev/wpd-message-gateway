package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/nyaruka/phonenumbers"

	"github.com/weprodev/wpd-message-gateway/internal/core/port"
	"github.com/weprodev/wpd-message-gateway/pkg/contracts"
)

const (
	ProviderName   = "telegram"
	defaultBaseURL = "https://gatewayapi.telegram.org/"
)

// Config holds Telegram-specific configuration.
type Config struct {
	APIKey  string
	BaseURL string
}

// Provider implements port.OTPSender for Telegram Gateway API.
type Provider struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

// Compile-time interface verification.
var _ port.OTPSender = (*Provider)(nil)

// New creates a new Telegram OTP provider.
func New(cfg Config) (*Provider, error) {
	baseURL := cfg.BaseURL
	if baseURL == "" {
		baseURL = defaultBaseURL
	}

	return &Provider{
		apiKey:  cfg.APIKey,
		baseURL: baseURL,
		client: &http.Client{},
	}, nil
}

// Name returns the provider name.
func (p *Provider) Name() string {
	return ProviderName
}

// Send sends an OTP via Telegram Gateway API.
func (p *Provider) Send(ctx context.Context, otp *contracts.OTP) (*contracts.SendResult, error) {
	for i, phone := range otp.PhoneNumber {
		num, err := phonenumbers.Parse(phone, "")
		if err != nil {
			return nil, fmt.Errorf("Failed to Parse Phone Number: %s", err)
		}
		if !phonenumbers.IsValidNumber(num) {
			return nil, fmt.Errorf("Invalid Phone Number")
		}
		otp.PhoneNumber[i] = phonenumbers.Format(num, phonenumbers.E164)
	}

	jsonData, err := json.Marshal(otp)
	if err != nil {
		return nil, fmt.Errorf("telegram: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.baseURL+"sendVerificationMessage", bytes.NewReader(jsonData))
	if err != nil {
		return nil, fmt.Errorf("telegram: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.apiKey)

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("telegram: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("telegram: %w", err)
	}

	var rawResp map[string]interface{}
	if err := json.Unmarshal(body, &rawResp); err != nil {
		return nil, fmt.Errorf("telegram: %w", err)
	}

	if ok, _ := rawResp["ok"].(bool); !ok {
		errMsg, _ := rawResp["error"].(string)
		return nil, fmt.Errorf("telegram: %s", errMsg)
	}

	result, _ := rawResp["result"].(map[string]interface{})
	id, _ := result["request_id"].(string)

	return &contracts.SendResult{
		ID:         id,
		StatusCode: resp.StatusCode,
		Message:    "OTP sent successfully",
	}, nil
}
