// Package logctx provides request-scoped correlation fields for structured logging.
// Core and presentation layers use this package — not infrastructure adapters.
package logctx

import "context"

type contextKey string

const (
	keyWorkspaceID contextKey = "workspace_id"
	keyAPIKeyID    contextKey = "api_key_id"
	keyChannel     contextKey = "channel"
	keyProvider    contextKey = "provider"
	keyRequestID   contextKey = "request_id"
)

func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, keyRequestID, requestID)
}

func GetRequestID(ctx context.Context) string {
	v, _ := ctx.Value(keyRequestID).(string)
	return v
}

func WithWorkspace(ctx context.Context, workspaceID, apiKeyID string) context.Context {
	ctx = context.WithValue(ctx, keyWorkspaceID, workspaceID)
	ctx = context.WithValue(ctx, keyAPIKeyID, apiKeyID)
	return ctx
}

func WithChannel(ctx context.Context, channel string) context.Context {
	return context.WithValue(ctx, keyChannel, channel)
}

func WithProvider(ctx context.Context, provider string) context.Context {
	return context.WithValue(ctx, keyProvider, provider)
}

// Attrs extracts gateway correlation fields from ctx for structured log handlers.
func Attrs(ctx context.Context) []any {
	if ctx == nil {
		return nil
	}
	attrs := make([]any, 0, 10)
	if v, _ := ctx.Value(keyWorkspaceID).(string); v != "" {
		attrs = append(attrs, "workspace_id", v)
	}
	if v, _ := ctx.Value(keyAPIKeyID).(string); v != "" {
		attrs = append(attrs, "api_key_id", v)
	}
	if v, _ := ctx.Value(keyChannel).(string); v != "" {
		attrs = append(attrs, "channel", v)
	}
	if v, _ := ctx.Value(keyProvider).(string); v != "" {
		attrs = append(attrs, "provider", v)
	}
	if v, _ := ctx.Value(keyRequestID).(string); v != "" {
		attrs = append(attrs, "request_id", v)
	}
	return attrs
}
