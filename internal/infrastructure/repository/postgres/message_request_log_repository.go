package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/weprodev/go-pkg/pgsql"

	"github.com/weprodev/wpd-message-gateway/internal/core/domain"
	"github.com/weprodev/wpd-message-gateway/internal/core/port"
)

type MessageRequestLogRepository struct {
	client *pgsql.PgClient
}

func NewMessageRequestLogRepository(client *pgsql.PgClient) port.MessageRequestLogRepository {
	return &MessageRequestLogRepository{client: client}
}

func (r *MessageRequestLogRepository) Create(ctx context.Context, log *domain.MessageRequestLog) error {
	query := `
		INSERT INTO message_request_logs (workspace_id, api_key_id, channel_type, http_method, status_code, endpoint, provider_name, request_id, duration_ms, error_message)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at
	`
	var apiKeyID interface{}
	if log.APIKeyID != "" {
		apiKeyID = log.APIKeyID
	}
	var reqID interface{}
	if log.RequestID != "" {
		reqID = log.RequestID
	}
	var dur interface{}
	if log.DurationMs > 0 {
		dur = log.DurationMs
	}
	var errMsg interface{}
	if log.ErrorMessage != "" {
		errMsg = log.ErrorMessage
	}
	var provider interface{}
	if log.ProviderName != "" {
		provider = log.ProviderName
	}

	err := r.client.GetDB(ctx).QueryRowContext(ctx, query,
		log.WorkspaceID, apiKeyID, log.ChannelType, log.HTTPMethod, log.StatusCode, log.Endpoint,
		provider, reqID, dur, errMsg,
	).Scan(&log.ID, &log.CreatedAt)
	if err != nil {
		slog.ErrorContext(ctx, "database error: failed to create message request log", "error", err, "endpoint", log.Endpoint)
		return err
	}

	if log.Payload != nil {
		if err := r.insertPayload(ctx, log); err != nil {
			// Metadata logs must persist even when optional payload storage fails.
			slog.WarnContext(ctx, "message request payload not stored", "error", err, "log_id", log.ID)
		}
	}

	return nil
}

func (r *MessageRequestLogRepository) insertPayload(ctx context.Context, log *domain.MessageRequestLog) error {
	log.Payload.LogID = log.ID
	payloadQuery := `
		INSERT INTO message_request_payloads (log_id, request_body, response_body)
		VALUES ($1, $2, $3)
		RETURNING created_at
	`
	err := r.client.GetDB(ctx).QueryRowContext(ctx, payloadQuery,
		log.Payload.LogID, log.Payload.RequestBody, log.Payload.ResponseBody,
	).Scan(&log.Payload.CreatedAt)
	if err != nil {
		return fmt.Errorf("insert message request payload: %w", err)
	}
	return nil
}

func (r *MessageRequestLogRepository) ListWithSource(ctx context.Context, q port.MessageLogQuery) ([]domain.MessageRequestLogWithSource, int, error) {
	db := r.client.GetDB(ctx)
	args := []any{q.WorkspaceID}
	where := "l.workspace_id = $1"
	argPos := 2
	if q.ChannelType != "" {
		where += fmt.Sprintf(" AND l.channel_type = $%d", argPos)
		args = append(args, q.ChannelType)
		argPos++
	}
	if q.From != nil {
		where += fmt.Sprintf(" AND l.created_at >= $%d", argPos)
		args = append(args, *q.From)
		argPos++
	}
	if q.To != nil {
		where += fmt.Sprintf(" AND l.created_at <= $%d", argPos)
		args = append(args, *q.To)
	}

	limit, offset := clampLimitOffset(q.Limit, q.Offset)
	listQuery := fmt.Sprintf(`
		SELECT l.id, l.workspace_id, l.api_key_id, l.channel_type, l.http_method, l.status_code, l.endpoint,
			l.provider_name, l.request_id, l.duration_ms, l.error_message, l.created_at,
			COALESCE(k.name, ''), COALESCE(k.client_id, ''),
			COUNT(*) OVER() AS total_count
		FROM message_request_logs l
		LEFT JOIN api_keys k ON k.id = l.api_key_id
		WHERE %s
		ORDER BY l.created_at DESC
		LIMIT %d OFFSET %d`, where, limit, offset)

	rows, err := db.QueryContext(ctx, listQuery, args...)
	if err != nil {
		slog.ErrorContext(ctx, "database error: failed to query message request logs", "error", err)
		return nil, 0, err
	}
	defer rows.Close() //nolint:errcheck

	var out []domain.MessageRequestLogWithSource
	total := 0
	for rows.Next() {
		var row domain.MessageRequestLogWithSource
		var apiKeyID sql.NullString
		var reqID sql.NullString
		var dur sql.NullInt64
		var errMsg sql.NullString
		var prov sql.NullString
		var totalCount int64
		if err := rows.Scan(
			&row.ID, &row.WorkspaceID, &apiKeyID, &row.ChannelType, &row.HTTPMethod, &row.StatusCode, &row.Endpoint,
			&prov, &reqID, &dur, &errMsg, &row.CreatedAt,
			&row.SourceName, &row.ClientID,
			&totalCount,
		); err != nil {
			slog.ErrorContext(ctx, "database error: failed to scan message request log", "error", err)
			return nil, 0, err
		}
		if total == 0 {
			total = int(totalCount)
		}
		if apiKeyID.Valid {
			row.APIKeyID = apiKeyID.String
		}
		if prov.Valid {
			row.ProviderName = prov.String
		}
		if reqID.Valid {
			row.RequestID = reqID.String
		}
		if dur.Valid {
			row.DurationMs = int(dur.Int64)
		}
		if errMsg.Valid {
			row.ErrorMessage = errMsg.String
		}
		out = append(out, row)
	}
	if err := rows.Err(); err != nil {
		slog.ErrorContext(ctx, "database error: rows iteration failed for message request logs", "error", err)
		return nil, 0, err
	}
	return out, total, nil
}
