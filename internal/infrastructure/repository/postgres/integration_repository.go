package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	"github.com/weprodev/go-pkg/pgsql"

	"github.com/weprodev/wpd-message-gateway/internal/core/domain"
	"github.com/weprodev/wpd-message-gateway/internal/core/port"
)

type IntegrationRepository struct {
	client *pgsql.PgClient
	enc    port.EncryptionService
}

func NewIntegrationRepository(client *pgsql.PgClient, enc port.EncryptionService) port.IntegrationRepository {
	return &IntegrationRepository{
		client: client,
		enc:    enc,
	}
}

func (r *IntegrationRepository) getProviderID(ctx context.Context, tx *sql.Tx, name, channelType string) (string, error) {
	var id string
	query := `SELECT id FROM providers WHERE name = $1 AND channel_type = $2`
	var err error
	if tx != nil {
		err = tx.QueryRowContext(ctx, query, name, channelType).Scan(&id)
	} else {
		err = r.client.GetDB(ctx).QueryRowContext(ctx, query, name, channelType).Scan(&id)
	}
	if errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("provider %s for channel %s not found: %w", name, channelType, port.ErrNotFound)
	}
	return id, err
}

func (r *IntegrationRepository) Create(ctx context.Context, integration *domain.Integration) error {
	encryptedConfig, err := r.enc.Encrypt(integration.Config)
	if err != nil {
		return err
	}

	providerID, err := r.getProviderID(ctx, nil, integration.ProviderName, integration.ChannelType)
	if err != nil {
		slog.ErrorContext(ctx, "database error: failed to resolve provider id for creation", "error", err, "workspace_id", integration.WorkspaceID, "provider_name", integration.ProviderName, "channel_type", integration.ChannelType)
		return err
	}

	query := `
		INSERT INTO integrations (workspace_id, channel_type, provider_id, encrypted_config, status, is_default)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	err = r.client.GetDB(ctx).QueryRowContext(ctx, query,
		integration.WorkspaceID, integration.ChannelType, providerID,
		encryptedConfig, integration.Status, integration.IsDefault,
	).Scan(&integration.ID, &integration.CreatedAt, &integration.UpdatedAt)

	if err != nil {
		slog.ErrorContext(ctx, "database error: failed to create integration", "error", err, "workspace_id", integration.WorkspaceID, "provider_name", integration.ProviderName, "channel_type", integration.ChannelType)
	}

	return err
}

func (r *IntegrationRepository) GetActiveByWorkspaceAndChannel(ctx context.Context, workspaceID, channelType string) (*domain.Integration, error) {
	query := `
		SELECT i.id, i.workspace_id, i.channel_type, p.name AS provider_name, i.encrypted_config, i.status, i.is_default, i.created_at, i.updated_at
		FROM integrations i
		JOIN providers p ON i.provider_id = p.id
		WHERE i.workspace_id = $1 AND i.channel_type = $2 AND i.status = 'connected'
		ORDER BY i.is_default DESC, i.created_at DESC
		LIMIT 1
	`
	var intg domain.Integration
	var encryptedConfig []byte

	err := r.client.GetDB(ctx).QueryRowContext(ctx, query, workspaceID, channelType).
		Scan(&intg.ID, &intg.WorkspaceID, &intg.ChannelType, &intg.ProviderName,
			&encryptedConfig, &intg.Status, &intg.IsDefault, &intg.CreatedAt, &intg.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("integration workspace=%s channel=%s: %w", workspaceID, channelType, port.ErrNotFound)
	}
	if err != nil {
		slog.ErrorContext(ctx, "database error: failed to get active integration", "error", err, "workspace_id", workspaceID, "channel_type", channelType)
		return nil, err
	}

	decryptedConfig, err := r.enc.Decrypt(encryptedConfig)
	if err != nil {
		return nil, err
	}
	intg.Config = decryptedConfig

	return &intg, nil
}

func (r *IntegrationRepository) ListByWorkspace(ctx context.Context, workspaceID string) ([]domain.Integration, error) {
	rows, err := r.client.GetDB(ctx).QueryContext(ctx, `
		SELECT i.id, i.workspace_id, i.channel_type, p.name AS provider_name, i.encrypted_config, i.status, i.is_default, i.created_at, i.updated_at
		FROM integrations i
		JOIN providers p ON i.provider_id = p.id
		WHERE i.workspace_id = $1
		ORDER BY i.channel_type ASC, p.name ASC
	`, workspaceID)
	if err != nil {
		slog.ErrorContext(ctx, "database error: failed to list integrations for workspace", "error", err, "workspace_id", workspaceID)
		return nil, err
	}
	defer rows.Close() //nolint:errcheck

	var out []domain.Integration
	for rows.Next() {
		var intg domain.Integration
		var enc []byte
		if err := rows.Scan(&intg.ID, &intg.WorkspaceID, &intg.ChannelType, &intg.ProviderName,
			&enc, &intg.Status, &intg.IsDefault, &intg.CreatedAt, &intg.UpdatedAt); err != nil {
			slog.ErrorContext(ctx, "database error: failed to scan integration in list", "error", err, "workspace_id", workspaceID)
			return nil, err
		}
		dec, err := r.enc.Decrypt(enc)
		if err != nil {
			return nil, err
		}
		intg.Config = dec
		out = append(out, intg)
	}
	if err := rows.Err(); err != nil {
		slog.ErrorContext(ctx, "database error: rows iteration failed for integrations", "error", err, "workspace_id", workspaceID)
		return nil, err
	}
	return out, nil
}

func (r *IntegrationRepository) GetByID(ctx context.Context, id string) (*domain.Integration, error) {
	query := `
		SELECT i.id, i.workspace_id, i.channel_type, p.name AS provider_name, i.encrypted_config, i.status, i.is_default, i.created_at, i.updated_at
		FROM integrations i
		JOIN providers p ON i.provider_id = p.id
		WHERE i.id = $1
	`
	var intg domain.Integration
	var enc []byte
	err := r.client.GetDB(ctx).QueryRowContext(ctx, query, id).
		Scan(&intg.ID, &intg.WorkspaceID, &intg.ChannelType, &intg.ProviderName,
			&enc, &intg.Status, &intg.IsDefault, &intg.CreatedAt, &intg.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("integration %s: %w", id, port.ErrNotFound)
	}
	if err != nil {
		slog.ErrorContext(ctx, "database error: failed to get integration by id", "error", err, "id", id)
		return nil, err
	}
	dec, err := r.enc.Decrypt(enc)
	if err != nil {
		return nil, err
	}
	intg.Config = dec
	return &intg, nil
}

func (r *IntegrationRepository) Delete(ctx context.Context, id string) error {
	res, err := r.client.GetDB(ctx).ExecContext(ctx, `DELETE FROM integrations WHERE id = $1`, id)
	if err != nil {
		slog.ErrorContext(ctx, "database error: failed to delete integration", "error", err, "id", id)
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("integration %s: %w", id, port.ErrNotFound)
	}
	return nil
}

func (r *IntegrationRepository) Upsert(ctx context.Context, integration *domain.Integration) error {
	encryptedConfig, err := r.enc.Encrypt(integration.Config)
	if err != nil {
		return err
	}

	err = r.client.RunInTransaction(ctx, func(txCtx context.Context) error {
		db := r.client.GetDB(txCtx)
		if integration.IsDefault {
			if _, err := db.ExecContext(txCtx, `
				UPDATE integrations SET is_default = FALSE
				WHERE workspace_id = $1 AND channel_type = $2
			`, integration.WorkspaceID, integration.ChannelType); err != nil {
				return err
			}
		}

		var providerID string
		err := db.QueryRowContext(txCtx, `SELECT id FROM providers WHERE name = $1 AND channel_type = $2`,
			integration.ProviderName, integration.ChannelType).Scan(&providerID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf("provider %s for channel %s not found: %w", integration.ProviderName, integration.ChannelType, port.ErrNotFound)
			}
			return err
		}

		query := `
			INSERT INTO integrations (workspace_id, channel_type, provider_id, encrypted_config, status, is_default)
			VALUES ($1, $2, $3, $4, $5, $6)
			ON CONFLICT (workspace_id, provider_id)
			DO UPDATE SET
				encrypted_config = EXCLUDED.encrypted_config,
				status = EXCLUDED.status,
				is_default = EXCLUDED.is_default,
				updated_at = NOW()
			RETURNING id, created_at, updated_at
		`
		return db.QueryRowContext(txCtx, query,
			integration.WorkspaceID, integration.ChannelType, providerID,
			encryptedConfig, integration.Status, integration.IsDefault,
		).Scan(&integration.ID, &integration.CreatedAt, &integration.UpdatedAt)
	})

	if err != nil {
		slog.ErrorContext(ctx, "database error: failed to upsert integration", "error", err, "workspace_id", integration.WorkspaceID, "provider_name", integration.ProviderName, "channel_type", integration.ChannelType)
	}
	return err
}

func (r *IntegrationRepository) GetProviderFields(ctx context.Context, providerName string) ([]domain.ProviderConfigField, error) {
	// Look up the provider first
	var providerID string
	err := r.client.GetDB(ctx).QueryRowContext(ctx, `SELECT id FROM providers WHERE name = $1`, providerName).Scan(&providerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.WarnContext(ctx, "provider not found", "provider_name", providerName)
			return nil, fmt.Errorf("provider %s: %w", providerName, port.ErrNotFound)
		}
		slog.ErrorContext(ctx, "database error: failed to lookup provider", "error", err, "provider_name", providerName)
		return nil, err
	}

	rows, err := r.client.GetDB(ctx).QueryContext(ctx, `
		SELECT id, provider_id, key, label, COALESCE(description, ''), field_type, required, COALESCE(default_value, ''), options, sort_order, created_at, updated_at
		FROM provider_config_fields
		WHERE provider_id = $1
		ORDER BY sort_order ASC, key ASC
	`, providerID)
	if err != nil {
		slog.ErrorContext(ctx, "database error: failed to query provider config fields", "error", err, "provider_name", providerName)
		return nil, err
	}
	defer rows.Close() //nolint:errcheck

	var fields []domain.ProviderConfigField
	for rows.Next() {
		var f domain.ProviderConfigField
		var options []byte
		var desc string
		var defVal string
		if err := rows.Scan(&f.ID, &f.ProviderID, &f.Key, &f.Label, &desc, &f.FieldType, &f.Required, &defVal, &options, &f.SortOrder, &f.CreatedAt, &f.UpdatedAt); err != nil {
			slog.ErrorContext(ctx, "database error: failed to scan provider config field", "error", err, "provider_name", providerName)
			return nil, err
		}
		f.Description = desc
		f.DefaultValue = defVal
		if len(options) > 0 {
			f.Options = json.RawMessage(options)
		}
		fields = append(fields, f)
	}
	if err := rows.Err(); err != nil {
		slog.ErrorContext(ctx, "database error: rows iteration failed for provider config fields", "error", err, "provider_name", providerName)
		return nil, err
	}
	return fields, nil
}
