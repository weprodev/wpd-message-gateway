package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/weprodev/go-pkg/pgsql"

	"github.com/weprodev/wpd-message-gateway/internal/core/port"
	"github.com/weprodev/wpd-message-gateway/pkg/contracts"
)

const (
	payloadInboxChannelEmail = "email"
	payloadInboxChannelSMS   = "sms"
	payloadInboxChannelPush  = "push"
	payloadInboxChannelChat  = "chat"

	payloadInboxDefaultPageSize = 50
	payloadInboxMaxPageSize     = 200
)

var (
	_ port.InboxReader = (*PayloadInboxReader)(nil)

	channelStatKeys = map[string]string{
		payloadInboxChannelEmail: "emails",
		payloadInboxChannelSMS:   "sms",
		payloadInboxChannelPush:  "push",
		payloadInboxChannelChat:  "chat",
	}
)

// PayloadInboxReader serves inbox reads from persisted message_request_payloads.
// An optional fallback reader covers the brief window before request logs are written.
type PayloadInboxReader struct {
	client   *pgsql.PgClient
	fallback port.InboxReader
}

func NewPayloadInboxReader(client *pgsql.PgClient, fallback port.InboxReader) *PayloadInboxReader {
	return &PayloadInboxReader{client: client, fallback: fallback}
}

type payloadInboxRow struct {
	logID          string
	workspaceID    string
	inboxMessageID sql.NullString
	createdAt      time.Time
	requestBody    string
}

func (row payloadInboxRow) messageID() string {
	if row.inboxMessageID.Valid && row.inboxMessageID.String != "" {
		return row.inboxMessageID.String
	}
	return row.logID
}

func clampPayloadInboxLimit(limit int) int {
	if limit <= 0 {
		return payloadInboxDefaultPageSize
	}
	if limit > payloadInboxMaxPageSize {
		return payloadInboxMaxPageSize
	}
	return limit
}

func unmarshalPayload[T any](requestBody string) (T, error) {
	var payload T
	if err := json.Unmarshal([]byte(requestBody), &payload); err != nil {
		return payload, err
	}
	return payload, nil
}

func listStoredMessages[T any](
	r *PayloadInboxReader,
	workspaceID, channel string,
	limit int,
	cursor string,
	parse func(payloadInboxRow) (T, error),
	idOf func(T) string,
) ([]T, string, bool) {
	rows, hasMore := r.queryPayloadRows(workspaceID, channel, limit, cursor)
	items := make([]T, 0, len(rows))
	for _, row := range rows {
		item, err := parse(row)
		if err != nil {
			slog.Warn("inbox payload unreadable", "error", err, "log_id", row.logID, "channel", channel)
			continue
		}
		items = append(items, item)
	}

	nextCursor := ""
	if hasMore && len(items) > 0 {
		nextCursor = idOf(items[len(items)-1])
	}
	return items, nextCursor, hasMore
}

func getStoredMessage[T any](primary, fallback func() (T, bool)) (T, bool) {
	if item, ok := primary(); ok {
		return item, true
	}
	return fallback()
}

func deleteStoredMessage(primary, fallback func() bool) bool {
	return primary() || fallback()
}

func (r *PayloadInboxReader) StatsForWorkspace(workspaceID string) map[string]int {
	ctx := context.Background()
	const query = `
		SELECT l.channel_type, COUNT(*)
		FROM message_request_logs l
		INNER JOIN message_request_payloads p ON p.log_id = l.id
		WHERE l.workspace_id = $1
		GROUP BY l.channel_type
	`

	rows, err := r.client.GetDB(ctx).QueryContext(ctx, query, workspaceID)
	if err != nil {
		slog.ErrorContext(ctx, "database error: failed to count inbox payloads", "error", err, "workspace_id", workspaceID)
		return emptyInboxStats()
	}
	defer rows.Close() //nolint:errcheck

	counts := emptyInboxStats()
	for rows.Next() {
		var channel string
		var count int
		if err := rows.Scan(&channel, &count); err != nil {
			continue
		}
		if key, ok := channelStatKeys[channel]; ok {
			counts[key] = count
		}
	}
	counts["total"] = counts["emails"] + counts["sms"] + counts["push"] + counts["chat"]
	return counts
}

func emptyInboxStats() map[string]int {
	return map[string]int{"emails": 0, "sms": 0, "push": 0, "chat": 0, "total": 0}
}

func (r *PayloadInboxReader) EmailsForWorkspace(workspaceID string) []port.StoredEmail {
	return r.ListEmailsForWorkspace(workspaceID, 0, "").Items
}

func (r *PayloadInboxReader) ListEmailsForWorkspace(workspaceID string, limit int, cursor string) port.InboxEmailPage {
	items, nextCursor, hasMore := listStoredMessages(
		r, workspaceID, payloadInboxChannelEmail, limit, cursor, parseStoredEmail,
		func(item port.StoredEmail) string { return item.ID },
	)
	return port.InboxEmailPage{Items: items, NextCursor: nextCursor, HasMore: hasMore}
}

func (r *PayloadInboxReader) EmailByIDForWorkspace(id, workspaceID string) (port.StoredEmail, bool) {
	return getStoredMessage(
		func() (port.StoredEmail, bool) { return r.storedEmailByMessageID(workspaceID, id) },
		func() (port.StoredEmail, bool) {
			if r.fallback == nil {
				return port.StoredEmail{}, false
			}
			return r.fallback.EmailByIDForWorkspace(id, workspaceID)
		},
	)
}

func (r *PayloadInboxReader) DeleteEmailByIDForWorkspace(id, workspaceID string) bool {
	return deleteStoredMessage(
		func() bool { return r.deletePayloadByMessageID(workspaceID, payloadInboxChannelEmail, id) },
		func() bool {
			if r.fallback == nil {
				return false
			}
			return r.fallback.DeleteEmailByIDForWorkspace(id, workspaceID)
		},
	)
}

func (r *PayloadInboxReader) SMSForWorkspace(workspaceID string) []port.StoredSMS {
	return r.ListSMSForWorkspace(workspaceID, 0, "").Items
}

func (r *PayloadInboxReader) ListSMSForWorkspace(workspaceID string, limit int, cursor string) port.InboxSMSPage {
	items, nextCursor, hasMore := listStoredMessages(
		r, workspaceID, payloadInboxChannelSMS, limit, cursor, parseStoredSMS,
		func(item port.StoredSMS) string { return item.ID },
	)
	return port.InboxSMSPage{Items: items, NextCursor: nextCursor, HasMore: hasMore}
}

func (r *PayloadInboxReader) SMSByIDForWorkspace(id, workspaceID string) (port.StoredSMS, bool) {
	return getStoredMessage(
		func() (port.StoredSMS, bool) { return r.storedSMSByMessageID(workspaceID, id) },
		func() (port.StoredSMS, bool) {
			if r.fallback == nil {
				return port.StoredSMS{}, false
			}
			return r.fallback.SMSByIDForWorkspace(id, workspaceID)
		},
	)
}

func (r *PayloadInboxReader) DeleteSMSByIDForWorkspace(id, workspaceID string) bool {
	return deleteStoredMessage(
		func() bool { return r.deletePayloadByMessageID(workspaceID, payloadInboxChannelSMS, id) },
		func() bool {
			if r.fallback == nil {
				return false
			}
			return r.fallback.DeleteSMSByIDForWorkspace(id, workspaceID)
		},
	)
}

func (r *PayloadInboxReader) PushForWorkspace(workspaceID string) []port.StoredPush {
	return r.ListPushForWorkspace(workspaceID, 0, "").Items
}

func (r *PayloadInboxReader) ListPushForWorkspace(workspaceID string, limit int, cursor string) port.InboxPushPage {
	items, nextCursor, hasMore := listStoredMessages(
		r, workspaceID, payloadInboxChannelPush, limit, cursor, parseStoredPush,
		func(item port.StoredPush) string { return item.ID },
	)
	return port.InboxPushPage{Items: items, NextCursor: nextCursor, HasMore: hasMore}
}

func (r *PayloadInboxReader) PushByIDForWorkspace(id, workspaceID string) (port.StoredPush, bool) {
	return getStoredMessage(
		func() (port.StoredPush, bool) { return r.storedPushByMessageID(workspaceID, id) },
		func() (port.StoredPush, bool) {
			if r.fallback == nil {
				return port.StoredPush{}, false
			}
			return r.fallback.PushByIDForWorkspace(id, workspaceID)
		},
	)
}

func (r *PayloadInboxReader) DeletePushByIDForWorkspace(id, workspaceID string) bool {
	return deleteStoredMessage(
		func() bool { return r.deletePayloadByMessageID(workspaceID, payloadInboxChannelPush, id) },
		func() bool {
			if r.fallback == nil {
				return false
			}
			return r.fallback.DeletePushByIDForWorkspace(id, workspaceID)
		},
	)
}

func (r *PayloadInboxReader) ChatForWorkspace(workspaceID string) []port.StoredChat {
	return r.ListChatForWorkspace(workspaceID, 0, "").Items
}

func (r *PayloadInboxReader) ListChatForWorkspace(workspaceID string, limit int, cursor string) port.InboxChatPage {
	items, nextCursor, hasMore := listStoredMessages(
		r, workspaceID, payloadInboxChannelChat, limit, cursor, parseStoredChat,
		func(item port.StoredChat) string { return item.ID },
	)
	return port.InboxChatPage{Items: items, NextCursor: nextCursor, HasMore: hasMore}
}

func (r *PayloadInboxReader) ChatByIDForWorkspace(id, workspaceID string) (port.StoredChat, bool) {
	return getStoredMessage(
		func() (port.StoredChat, bool) { return r.storedChatByMessageID(workspaceID, id) },
		func() (port.StoredChat, bool) {
			if r.fallback == nil {
				return port.StoredChat{}, false
			}
			return r.fallback.ChatByIDForWorkspace(id, workspaceID)
		},
	)
}

func (r *PayloadInboxReader) DeleteChatByIDForWorkspace(id, workspaceID string) bool {
	return deleteStoredMessage(
		func() bool { return r.deletePayloadByMessageID(workspaceID, payloadInboxChannelChat, id) },
		func() bool {
			if r.fallback == nil {
				return false
			}
			return r.fallback.DeleteChatByIDForWorkspace(id, workspaceID)
		},
	)
}

func (r *PayloadInboxReader) ClearWorkspace(workspaceID string) {
	ctx := context.Background()
	const query = `
		DELETE FROM message_request_payloads
		WHERE log_id IN (
			SELECT id FROM message_request_logs WHERE workspace_id = $1
		)
	`
	if _, err := r.client.GetDB(ctx).ExecContext(ctx, query, workspaceID); err != nil {
		slog.ErrorContext(ctx, "database error: failed to clear inbox payloads", "error", err, "workspace_id", workspaceID)
	}
	if r.fallback != nil {
		r.fallback.ClearWorkspace(workspaceID)
	}
}

func (r *PayloadInboxReader) storedEmailByMessageID(workspaceID, messageID string) (port.StoredEmail, bool) {
	return storedMessageByID(r, workspaceID, payloadInboxChannelEmail, messageID, parseStoredEmail)
}

func (r *PayloadInboxReader) storedSMSByMessageID(workspaceID, messageID string) (port.StoredSMS, bool) {
	return storedMessageByID(r, workspaceID, payloadInboxChannelSMS, messageID, parseStoredSMS)
}

func (r *PayloadInboxReader) storedPushByMessageID(workspaceID, messageID string) (port.StoredPush, bool) {
	return storedMessageByID(r, workspaceID, payloadInboxChannelPush, messageID, parseStoredPush)
}

func (r *PayloadInboxReader) storedChatByMessageID(workspaceID, messageID string) (port.StoredChat, bool) {
	return storedMessageByID(r, workspaceID, payloadInboxChannelChat, messageID, parseStoredChat)
}

func storedMessageByID[T any](
	r *PayloadInboxReader,
	workspaceID, channel, messageID string,
	parse func(payloadInboxRow) (T, error),
) (T, bool) {
	var zero T
	row, ok := r.payloadRowByMessageID(workspaceID, channel, messageID)
	if !ok {
		return zero, false
	}
	item, err := parse(row)
	if err != nil {
		return zero, false
	}
	return item, true
}

func (r *PayloadInboxReader) queryPayloadRows(workspaceID, channel string, limit int, cursor string) ([]payloadInboxRow, bool) {
	ctx := context.Background()
	limit = clampPayloadInboxLimit(limit)
	fetchLimit := limit + 1

	cursorCreatedAt, cursorLogID, hasCursor := r.resolvePayloadCursor(ctx, workspaceID, channel, cursor)

	const baseSelect = `
		SELECT l.id, l.workspace_id, l.inbox_message_id, l.created_at, p.request_body
		FROM message_request_logs l
		INNER JOIN message_request_payloads p ON p.log_id = l.id
		WHERE l.workspace_id = $1 AND l.channel_type = $2
	`

	var (
		rows *sql.Rows
		err  error
	)
	if hasCursor {
		query := baseSelect + `
			AND (l.created_at, l.id) < ($3, $4)
			ORDER BY l.created_at DESC, l.id DESC
			LIMIT $5
		`
		rows, err = r.client.GetDB(ctx).QueryContext(ctx, query, workspaceID, channel, cursorCreatedAt, cursorLogID, fetchLimit)
	} else {
		query := baseSelect + `
			ORDER BY l.created_at DESC, l.id DESC
			LIMIT $3
		`
		rows, err = r.client.GetDB(ctx).QueryContext(ctx, query, workspaceID, channel, fetchLimit)
	}
	if err != nil {
		slog.ErrorContext(ctx, "database error: failed to list inbox payloads", "error", err, "workspace_id", workspaceID, "channel", channel)
		return nil, false
	}
	defer rows.Close() //nolint:errcheck

	var items []payloadInboxRow
	for rows.Next() {
		var row payloadInboxRow
		if err := rows.Scan(&row.logID, &row.workspaceID, &row.inboxMessageID, &row.createdAt, &row.requestBody); err != nil {
			continue
		}
		items = append(items, row)
	}
	if err := rows.Err(); err != nil {
		slog.ErrorContext(ctx, "database error: inbox payload list iteration failed", "error", err)
		return nil, false
	}

	hasMore := len(items) > limit
	if hasMore {
		items = items[:limit]
	}
	return items, hasMore
}

func (r *PayloadInboxReader) resolvePayloadCursor(ctx context.Context, workspaceID, channel, cursor string) (time.Time, string, bool) {
	if cursor == "" {
		return time.Time{}, "", false
	}

	const query = `
		SELECT l.created_at, l.id
		FROM message_request_logs l
		INNER JOIN message_request_payloads p ON p.log_id = l.id
		WHERE l.workspace_id = $1 AND l.channel_type = $2
		  AND (l.inbox_message_id = $3 OR l.id::text = $3)
	`
	var createdAt time.Time
	var logID string
	err := r.client.GetDB(ctx).QueryRowContext(ctx, query, workspaceID, channel, cursor).Scan(&createdAt, &logID)
	if errors.Is(err, sql.ErrNoRows) {
		return time.Time{}, "", false
	}
	if err != nil {
		slog.ErrorContext(ctx, "database error: failed to resolve inbox payload cursor", "error", err, "cursor", cursor)
		return time.Time{}, "", false
	}
	return createdAt, logID, true
}

func (r *PayloadInboxReader) payloadRowByMessageID(workspaceID, channel, messageID string) (payloadInboxRow, bool) {
	ctx := context.Background()
	const query = `
		SELECT l.id, l.workspace_id, l.inbox_message_id, l.created_at, p.request_body
		FROM message_request_logs l
		INNER JOIN message_request_payloads p ON p.log_id = l.id
		WHERE l.workspace_id = $1 AND l.channel_type = $2
		  AND (l.inbox_message_id = $3 OR l.id::text = $3)
	`

	var row payloadInboxRow
	err := r.client.GetDB(ctx).QueryRowContext(ctx, query, workspaceID, channel, messageID).
		Scan(&row.logID, &row.workspaceID, &row.inboxMessageID, &row.createdAt, &row.requestBody)
	if errors.Is(err, sql.ErrNoRows) {
		return payloadInboxRow{}, false
	}
	if err != nil {
		slog.ErrorContext(ctx, "database error: failed to get inbox payload", "error", err, "message_id", messageID)
		return payloadInboxRow{}, false
	}
	return row, true
}

func (r *PayloadInboxReader) deletePayloadByMessageID(workspaceID, channel, messageID string) bool {
	ctx := context.Background()
	const query = `
		DELETE FROM message_request_payloads
		WHERE log_id IN (
			SELECT l.id
			FROM message_request_logs l
			WHERE l.workspace_id = $1 AND l.channel_type = $2
			  AND (l.inbox_message_id = $3 OR l.id::text = $3)
		)
	`
	res, err := r.client.GetDB(ctx).ExecContext(ctx, query, workspaceID, channel, messageID)
	if err != nil {
		slog.ErrorContext(ctx, "database error: failed to delete inbox payload", "error", err, "message_id", messageID)
		return false
	}
	n, err := res.RowsAffected()
	return err == nil && n > 0
}

func parseStoredEmail(row payloadInboxRow) (port.StoredEmail, error) {
	email, err := unmarshalPayload[contracts.Email](row.requestBody)
	if err != nil {
		return port.StoredEmail{}, fmt.Errorf("unmarshal email payload: %w", err)
	}
	return port.StoredEmail{
		ID: row.messageID(), WorkspaceID: row.workspaceID, CreatedAt: row.createdAt, Email: email,
	}, nil
}

func parseStoredSMS(row payloadInboxRow) (port.StoredSMS, error) {
	sms, err := unmarshalPayload[contracts.SMS](row.requestBody)
	if err != nil {
		return port.StoredSMS{}, fmt.Errorf("unmarshal sms payload: %w", err)
	}
	return port.StoredSMS{
		ID: row.messageID(), WorkspaceID: row.workspaceID, CreatedAt: row.createdAt, SMS: sms,
	}, nil
}

func parseStoredPush(row payloadInboxRow) (port.StoredPush, error) {
	push, err := unmarshalPayload[contracts.PushNotification](row.requestBody)
	if err != nil {
		return port.StoredPush{}, fmt.Errorf("unmarshal push payload: %w", err)
	}
	return port.StoredPush{
		ID: row.messageID(), WorkspaceID: row.workspaceID, CreatedAt: row.createdAt, Push: push,
	}, nil
}

func parseStoredChat(row payloadInboxRow) (port.StoredChat, error) {
	chat, err := unmarshalPayload[contracts.ChatMessage](row.requestBody)
	if err != nil {
		return port.StoredChat{}, fmt.Errorf("unmarshal chat payload: %w", err)
	}
	return port.StoredChat{
		ID: row.messageID(), WorkspaceID: row.workspaceID, CreatedAt: row.createdAt, Chat: chat,
	}, nil
}
