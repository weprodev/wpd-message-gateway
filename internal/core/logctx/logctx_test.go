package logctx

import (
	"context"
	"testing"
)

func TestWithProvider_addsProviderToAttrs(t *testing.T) {
	t.Parallel()

	ctx := WithProvider(context.Background(), "mailgun")
	attrs := Attrs(ctx)

	if len(attrs) != 2 || attrs[0] != "provider" || attrs[1] != "mailgun" {
		t.Fatalf("Attrs: got %#v", attrs)
	}
}

func TestWithWorkspace_addsWorkspaceFields(t *testing.T) {
	t.Parallel()

	ctx := WithWorkspace(context.Background(), "ws-1", "key-1")
	attrs := Attrs(ctx)

	if len(attrs) != 4 {
		t.Fatalf("expected 4 attr entries, got %#v", attrs)
	}
}

func TestWithRequestID_roundTrip(t *testing.T) {
	t.Parallel()

	ctx := WithRequestID(context.Background(), "req-123")
	if GetRequestID(ctx) != "req-123" {
		t.Fatalf("GetRequestID: got %q", GetRequestID(ctx))
	}
	attrs := Attrs(ctx)
	if len(attrs) != 2 || attrs[0] != "request_id" || attrs[1] != "req-123" {
		t.Fatalf("Attrs: got %#v", attrs)
	}
}
