package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/sync/errgroup"

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
		client:  &http.Client{},
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

	senderUsernames, err := normalizeSenderUsernames(otp.PhoneNumber, otp.SenderUsername)
	if err != nil {
		return nil, fmt.Errorf("telegram: %w", err)
	}

	if len(otp.PhoneNumber) == 1 {
		return p.sendSingle(ctx, otp.PhoneNumber[0], senderUsernames[0], otp.CodeLength)
	}

	items := make([]contracts.SendResultItem, len(otp.PhoneNumber))
	g, ctx := errgroup.WithContext(ctx)

	for i, phone := range otp.PhoneNumber {
		i, phone := i, phone
		g.Go(func() error {
			item, err := p.sendToOne(ctx, phone, senderUsernames[i], otp.CodeLength)
			if err != nil {
				items[i] = contracts.SendResultItem{
					PhoneNumber: phone,
					Error:       err.Error(),
				}
				return err
			}
			items[i] = *item
			return nil
		})
	}

	err = g.Wait()

	statusCode := 200
	message := "OTP sent successfully"
	if err != nil {
		statusCode = 500
		message = "one or more OTPs failed"
	}

	return &contracts.SendResult{
		StatusCode: statusCode,
		Message:    message,
		Items:      items,
	}, nil
}

// sendSingle handles the single-phone case for backward compatibility.
func (p *Provider) sendSingle(ctx context.Context, phone, senderUsername string, codeLength int) (*contracts.SendResult, error) {
	item, err := p.sendToOne(ctx, phone, senderUsername, codeLength)
	if err != nil {
		return nil, err
	}
	return &contracts.SendResult{
		ID:         item.RequestID,
		StatusCode: item.StatusCode,
		Message:    "OTP sent successfully",
		Items:      []contracts.SendResultItem{*item},
	}, nil
}

// sendToOne sends a verification message to a single phone number.
func (p *Provider) sendToOne(ctx context.Context, phone, senderUsername string, codeLength int) (*contracts.SendResultItem, error) {
	body := map[string]interface{}{
		"phone_number": phone,
	}
	if senderUsername != "" {
		body["sender_username"] = senderUsername
	}
	if codeLength > 0 {
		body["code_length"] = codeLength
	}

	jsonData, err := json.Marshal(body)
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

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("telegram: %w", err)
	}

	var rawResp map[string]interface{}
	if err := json.Unmarshal(respBody, &rawResp); err != nil {
		return nil, fmt.Errorf("telegram: %w", err)
	}

	if ok, _ := rawResp["ok"].(bool); !ok {
		errMsg, _ := rawResp["error"].(string)
		return &contracts.SendResultItem{
			PhoneNumber: phone,
			StatusCode:  resp.StatusCode,
			Error:       errMsg,
		}, fmt.Errorf("telegram: %s", errMsg)
	}

	result, _ := rawResp["result"].(map[string]interface{})
	id, _ := result["request_id"].(string)

	return &contracts.SendResultItem{
		PhoneNumber: phone,
		RequestID:   id,
		StatusCode:  resp.StatusCode,
	}, nil
}

// normalizeSenderUsernames expands sender usernames to match the phone count.
func normalizeSenderUsernames(phoneNumbers, senderUsernames []string) ([]string, error) {
	switch len(senderUsernames) {
	case 0:
		return make([]string, len(phoneNumbers)), nil
	case 1:
		result := make([]string, len(phoneNumbers))
		for i := range result {
			result[i] = senderUsernames[0]
		}
		return result, nil
	default:
		if len(senderUsernames) != len(phoneNumbers) {
			return nil, fmt.Errorf(
				"sender_username count (%d) must be 0, 1, or match phone_number count (%d)",
				len(senderUsernames), len(phoneNumbers),
			)
		}
		return senderUsernames, nil
	}
}
