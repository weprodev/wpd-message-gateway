package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"testing"
)

func TestContextHandler(t *testing.T) {
	var buf bytes.Buffer
	baseHandler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	handler := NewContextHandler(baseHandler)
	testLogger := slog.New(handler)

	ctx := context.Background()
	ctx = WithWorkspace(ctx, "ws-123", "key-456")
	ctx = WithChannel(ctx, "sms")
	ctx = WithProvider(ctx, "twilio")
	ctx = WithRequestID(ctx, "req-789")

	testLogger.InfoContext(ctx, "test log message", "extra", "value")

	var result map[string]any
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("failed to parse log JSON: %v", err)
	}

	expected := map[string]string{
		"msg":          "test log message",
		"workspace_id": "ws-123",
		"api_key_id":   "key-456",
		"channel":      "sms",
		"provider":     "twilio",
		"request_id":   "req-789",
		"extra":        "value",
	}

	for k, v := range expected {
		val, ok := result[k]
		if !ok {
			t.Errorf("missing expected key: %s", k)
			continue
		}
		if valStr, ok := val.(string); !ok || valStr != v {
			t.Errorf("expected key %s to be %s, got %v", k, v, val)
		}
	}
}
