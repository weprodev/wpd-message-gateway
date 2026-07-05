package service

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"

	"github.com/weprodev/wpd-message-gateway/internal/core/port"
	"github.com/weprodev/wpd-message-gateway/pkg/contracts"
	"github.com/weprodev/wpd-message-gateway/pkg/registry"
)

func applyEmailSenderDefaults(email contracts.Email, cfg registry.EmailConfig) contracts.Email {
	enriched := email
	if enriched.From == "" && cfg.FromEmail != "" {
		enriched.From = cfg.FromEmail
	}
	if enriched.FromName == "" && cfg.FromName != "" {
		enriched.FromName = cfg.FromName
	}
	return enriched
}

func (s *GatewayService) enrichEmailFromIntegration(ctx context.Context, workspaceID string, email contracts.Email) contracts.Email {
	if email.From != "" && email.FromName != "" {
		return email
	}
	if s.integrations == nil {
		return email
	}

	intg, err := s.integrations.GetActiveByWorkspaceAndChannel(ctx, workspaceID, channelEmail)
	if err != nil || intg == nil {
		if err != nil && !errors.Is(err, port.ErrNotFound) {
			slog.WarnContext(ctx, "failed to load email integration for sender enrichment", "error", err)
		}
		return email
	}

	var cfg registry.EmailConfig
	if err := json.Unmarshal(intg.Config, &cfg); err != nil {
		slog.WarnContext(ctx, "failed to parse email integration config for sender enrichment", "error", err)
		return email
	}

	return applyEmailSenderDefaults(email, cfg)
}
