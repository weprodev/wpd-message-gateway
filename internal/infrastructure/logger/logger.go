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

	"github.com/weprodev/wpd-message-gateway/internal/core/logctx"
)

// ContextHandler extracts request context fields and appends them to slog.Record.
type ContextHandler struct {
	next slog.Handler
}

// NewContextHandler wraps another slog.Handler.
func NewContextHandler(next slog.Handler) *ContextHandler {
	return &ContextHandler{next: next}
}

func (h *ContextHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.next.Enabled(ctx, level)
}

func (h *ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	if ctx == nil {
		return h.next.Handle(ctx, r)
	}

	attrs := logctx.Attrs(ctx)
	for i := 0; i < len(attrs); i += 2 {
		key, ok := attrs[i].(string)
		if ok {
			r.AddAttrs(slog.Any(key, attrs[i+1]))
		}
	}

	return h.next.Handle(ctx, r)
}

func (h *ContextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &ContextHandler{next: h.next.WithAttrs(attrs)}
}

func (h *ContextHandler) WithGroup(name string) slog.Handler {
	return &ContextHandler{next: h.next.WithGroup(name)}
}

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

	extractors := []pkglogger.ContextExtractor{
		func(ctx context.Context) []any {
			return logctx.Attrs(ctx)
		},
	}

	log, err := pkglogger.New(pkglogger.Config{
		Level:             level,
		Format:            format,
		OutputPath:        "stdout",
		ExtraHandlers:     extraHandlers,
		ContextExtractors: extractors,
	})
	if err != nil {
		return nil, err
	}

	// ContextExtractors append correlation fields; avoid also wrapping ContextHandler (duplicate attrs).
	slog.SetDefault(log.Logger)

	return log, nil
}

// WithRequestID returns a context enriched with request_id (correlation ID).
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return logctx.WithRequestID(ctx, requestID)
}

// GetRequestID retrieves the request_id from context.
func GetRequestID(ctx context.Context) string {
	return logctx.GetRequestID(ctx)
}

// WithWorkspace returns a context enriched with workspace and api-key identifiers.
func WithWorkspace(ctx context.Context, workspaceID, apiKeyID string) context.Context {
	return logctx.WithWorkspace(ctx, workspaceID, apiKeyID)
}

// WithChannel enriches a context with the active message channel (email/sms/push/chat).
func WithChannel(ctx context.Context, channel string) context.Context {
	return logctx.WithChannel(ctx, channel)
}

// WithProvider enriches a context with the resolved provider name.
func WithProvider(ctx context.Context, provider string) context.Context {
	return logctx.WithProvider(ctx, provider)
}

// Attrs extracts gateway correlation fields from ctx for go-pkg ContextExtractors.
func Attrs(ctx context.Context) []any {
	return logctx.Attrs(ctx)
}

func isLocalEnv(env string) bool {
	return env == "" || env == "local" || env == "test" || env == "development"
}
