package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

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

func (r *IntegrationRepository) Create(ctx context.Context, integration *domain.Integration) error {
	encryptedConfig, err := r.enc.Encrypt(integration.Config)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO integrations (workspace_id, channel_type, provider_name, encrypted_config, status, is_default)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	err = r.client.GetDB(ctx).QueryRowContext(ctx, query,
		integration.WorkspaceID, integration.ChannelType, integration.ProviderName,
		encryptedConfig, integration.Status, integration.IsDefault,
	).Scan(&integration.ID, &integration.CreatedAt, &integration.UpdatedAt)

	return err
}

func (r *IntegrationRepository) GetActiveByWorkspaceAndChannel(ctx context.Context, workspaceID, channelType string) (*domain.Integration, error) {
	query := `
		SELECT id, workspace_id, channel_type, provider_name, encrypted_config, status, is_default, created_at, updated_at
		FROM integrations
		WHERE workspace_id = $1 AND channel_type = $2 AND status = 'connected'
		ORDER BY is_default DESC, created_at DESC
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
		SELECT id, workspace_id, channel_type, provider_name, encrypted_config, status, is_default, created_at, updated_at
		FROM integrations
		WHERE workspace_id = $1
		ORDER BY channel_type ASC, provider_name ASC
	`, workspaceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close() //nolint:errcheck

	var out []domain.Integration
	for rows.Next() {
		var intg domain.Integration
		var enc []byte
		if err := rows.Scan(&intg.ID, &intg.WorkspaceID, &intg.ChannelType, &intg.ProviderName,
			&enc, &intg.Status, &intg.IsDefault, &intg.CreatedAt, &intg.UpdatedAt); err != nil {
			return nil, err
		}
		dec, err := r.enc.Decrypt(enc)
		if err != nil {
			return nil, err
		}
		intg.Config = dec
		out = append(out, intg)
	}
	return out, rows.Err()
}

func (r *IntegrationRepository) GetByID(ctx context.Context, id string) (*domain.Integration, error) {
	query := `
		SELECT id, workspace_id, channel_type, provider_name, encrypted_config, status, is_default, created_at, updated_at
		FROM integrations WHERE id = $1
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

	return r.client.RunInTransaction(ctx, func(txCtx context.Context) error {
		db := r.client.GetDB(txCtx)
		if integration.IsDefault {
			if _, err := db.ExecContext(txCtx, `
				UPDATE integrations SET is_default = FALSE
				WHERE workspace_id = $1 AND channel_type = $2
			`, integration.WorkspaceID, integration.ChannelType); err != nil {
				return err
			}
		}
		query := `
			INSERT INTO integrations (workspace_id, channel_type, provider_name, encrypted_config, status, is_default)
			VALUES ($1, $2, $3, $4, $5, $6)
			ON CONFLICT (workspace_id, channel_type, provider_name)
			DO UPDATE SET
				encrypted_config = EXCLUDED.encrypted_config,
				status = EXCLUDED.status,
				is_default = EXCLUDED.is_default,
				updated_at = NOW()
			RETURNING id, created_at, updated_at
		`
		return db.QueryRowContext(txCtx, query,
			integration.WorkspaceID, integration.ChannelType, integration.ProviderName,
			encryptedConfig, integration.Status, integration.IsDefault,
		).Scan(&integration.ID, &integration.CreatedAt, &integration.UpdatedAt)
	})
}
