// Package logger provides the application-wide structured logger for wpd-message-gateway.
//
// It wraps github.com/weprodev/go-pkg/logger and installs it as the process-wide
// slog default via slog.SetDefault(). After calling New(), every slog.Info(),
// slog.Warn(), slog.Error(), and slog.ErrorContext() call in the codebase —
// across handlers, services, and infrastructure — routes through the configured
// handler without threading a *Logger pointer through every struct.
//
// # Format selection
//
//   - local / development / test  → text format, DEBUG level (human-readable)
//   - any other environment        → JSON format, INFO level (machine-parseable)
//
// # APM integration
//
// Pass ExtraHandlers to fan out to Sentry, OpenTelemetry, etc.:
//
//	log, _ := logger.New("production", slog.Handler(sentryHandler))
package logger

import (
	"context"
	"log/slog"

	pkglogger "github.com/weprodev/go-pkg/logger"
)

// contextKey is an unexported type for keys stored in context.Context.
type contextKey string

const (
	keyWorkspaceID contextKey = "workspace_id"
	keyAPIKeyID    contextKey = "api_key_id"
	keyChannel     contextKey = "channel"
	keyProvider    contextKey = "provider"
)

// New initialises the structured logger and installs it as the global slog default.
// Subsequent external handlers (e.g. Sentry) can be injected via extraHandlers.
//
// Call this once at process start — before loading config or wiring dependencies.
func New(env string, extraHandlers ...slog.Handler) (*pkglogger.Logger, error) {
	format := pkglogger.FormatText
	level := pkglogger.LevelDebug
	if !isLocalEnv(env) {
		format = pkglogger.FormatJSON
		level = pkglogger.LevelInfo
	}

	log, err := pkglogger.New(pkglogger.Config{
		Level:         level,
		Format:        format,
		OutputPath:    "stdout",
		ExtraHandlers: extraHandlers,
	})
	if err != nil {
		return nil, err
	}

	// All slog.X() calls throughout the process now use this logger.
	slog.SetDefault(log.Logger)

	return log, nil
}

// WithWorkspace returns a context enriched with workspace and api-key identifiers.
// Pass this context into handler calls so structured log records carry these fields.
func WithWorkspace(ctx context.Context, workspaceID, apiKeyID string) context.Context {
	ctx = context.WithValue(ctx, keyWorkspaceID, workspaceID)
	ctx = context.WithValue(ctx, keyAPIKeyID, apiKeyID)
	return ctx
}

// WithChannel enriches a context with the active message channel (email/sms/push/chat).
func WithChannel(ctx context.Context, channel string) context.Context {
	return context.WithValue(ctx, keyChannel, channel)
}

// WithProvider enriches a context with the resolved provider name.
func WithProvider(ctx context.Context, provider string) context.Context {
	return context.WithValue(ctx, keyProvider, provider)
}

// Attrs extracts all gateway-specific attributes from ctx into slog key-value pairs.
// Use this when calling slog.InfoContext / slog.ErrorContext to attach request context:
//
//	slog.ErrorContext(ctx, "send failed", append(logger.Attrs(ctx), "error", err)...)
func Attrs(ctx context.Context) []any {
	attrs := make([]any, 0, 8)
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
	return attrs
}

func isLocalEnv(env string) bool {
	return env == "" || env == "local" || env == "test" || env == "development"
}
