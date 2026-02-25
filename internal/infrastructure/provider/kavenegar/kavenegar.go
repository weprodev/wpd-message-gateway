package kavenegar

import (
	"context"
	"errors"
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/kavenegar/kavenegar-go"
	"github.com/nyaruka/phonenumbers"

	"github.com/weprodev/wpd-message-gateway/pkg/contracts"
)

const (
	ProviderName   = "kavenegar"
	defaultTimeout = 40 * time.Second
)

type Config struct {
	APIKey    string
	FromPhone string
}

type Provider struct {
	client    *kavenegar.Kavenegar
	config    Config
	fromPhone string
}

func New(cfg Config) (*Provider, error) {

	if cfg.APIKey == "" {
		return nil, errors.New("kavenegar: API Key is required")

	}

	if cfg.FromPhone == "" {
		return nil, errors.New("kavenegar: Sender number is required")

	}

	return &Provider{
		client:    kavenegar.New(cfg.APIKey),
		config:    cfg,
		fromPhone: cfg.FromPhone,
	}, nil

}

func (p *Provider) Name() string {
	return ProviderName

}

func (p *Provider) Send(ctx context.Context, sms *contracts.SMS) (*contracts.SendResult, error) {

	if sms.From == "" {
		sms.From = p.fromPhone

	}

	if len(sms.To) == 0 {
		return nil, errors.New("at least one recipient is required")

	}

	messageLength := utf8.RuneCountInString(sms.Message)

	if messageLength > 160 {
		return nil, errors.New("message is too long, Limit is 160 characters")

	}

	seen := make(map[string]bool)

	for _, recipient := range sms.To {

		if seen[recipient] {
			return nil, fmt.Errorf("duplicate message detected: %s", recipient)
		}
		seen[recipient] = true

		num, err := phonenumbers.Parse(recipient, "")

		if err != nil {
			return nil, fmt.Errorf("failed to parse phone number %s: %w", recipient, err)

		}

		if !phonenumbers.IsValidNumber(num) {
			return nil, fmt.Errorf("the phone number %s is not valid", recipient)

		}
	}

	sendCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	res, err := p.client.Message.Send(sms.From, sms.To, sms.Message, nil)

	if err != nil {
		if errors.Is(sendCtx.Err(), context.DeadlineExceeded) {
			return nil, fmt.Errorf("request timed out after %v", defaultTimeout)
		}
		return nil, fmt.Errorf("kavenegar API error: %w", err)
	}

	if len(res) == 0 {
		return nil, errors.New("kavenegar returned an empty response entries list")

	}

	result := &contracts.SendResult{
		ID:         fmt.Sprintf("%d", res[0].MessageID),
		StatusCode: 200,
		Message:    "SMS sent successfully",
	}
	return result, nil

}
