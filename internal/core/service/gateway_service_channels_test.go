package service

import (
	"context"
	"testing"

	"github.com/weprodev/wpd-message-gateway/internal/core/domain"
	"github.com/weprodev/wpd-message-gateway/pkg/contracts"
)

func TestGatewayService_SendSMS_memoryOnly(t *testing.T) {
	t.Parallel()

	inbox := &stubInbox{}
	svc := NewGatewayService(&stubIntegrationRepo{}, nil, nil, inbox, nil)

	res, err := svc.SendSMS(context.Background(), "ws-1", contracts.SMS{
		To:      []string{"+15550001111"},
		Message: "hello",
	})
	if err != nil {
		t.Fatalf("SendSMS: %v", err)
	}
	if res.ID != "inbox-sms-1" {
		t.Fatalf("got ID %q", res.ID)
	}
	if res.Meta[contracts.MetaKeyDispatchMode] != string(domain.DispatchMemory) {
		t.Fatalf("dispatch_mode: %v", res.Meta[contracts.MetaKeyDispatchMode])
	}
	if res.Meta["channel"] != "sms" {
		t.Fatalf("channel: %v", res.Meta["channel"])
	}
}

func TestGatewayService_SendPush_memoryOnly(t *testing.T) {
	t.Parallel()

	inbox := &stubInbox{}
	svc := NewGatewayService(&stubIntegrationRepo{}, nil, nil, inbox, nil)

	res, err := svc.SendPush(context.Background(), "ws-1", contracts.PushNotification{
		DeviceTokens: []string{"device-token-1"},
		Title:        "Hi",
		Body:         "There",
	})
	if err != nil {
		t.Fatalf("SendPush: %v", err)
	}
	if res.ID != "inbox-push-1" {
		t.Fatalf("got ID %q", res.ID)
	}
	if res.Meta[contracts.MetaKeyDispatchMode] != string(domain.DispatchMemory) {
		t.Fatalf("dispatch_mode: %v", res.Meta[contracts.MetaKeyDispatchMode])
	}
	if res.Meta["channel"] != "push" {
		t.Fatalf("channel: %v", res.Meta["channel"])
	}
}

func TestGatewayService_SendChat_memoryOnly(t *testing.T) {
	t.Parallel()

	inbox := &stubInbox{}
	svc := NewGatewayService(&stubIntegrationRepo{}, nil, nil, inbox, nil)

	res, err := svc.SendChat(context.Background(), "ws-1", contracts.ChatMessage{
		To:      []string{"user-1"},
		Message: "ping",
	})
	if err != nil {
		t.Fatalf("SendChat: %v", err)
	}
	if res.ID != "inbox-chat-1" {
		t.Fatalf("got ID %q", res.ID)
	}
	if res.Meta[contracts.MetaKeyDispatchMode] != string(domain.DispatchMemory) {
		t.Fatalf("dispatch_mode: %v", res.Meta[contracts.MetaKeyDispatchMode])
	}
	if res.Meta["channel"] != "chat" {
		t.Fatalf("channel: %v", res.Meta["channel"])
	}
}

func TestGatewayService_SendSMS_memoryOnly_inboxNil(t *testing.T) {
	t.Parallel()

	svc := NewGatewayService(&stubIntegrationRepo{}, nil, nil, nil, nil)

	_, err := svc.SendSMS(context.Background(), "ws-1", contracts.SMS{To: []string{"+1"}, Message: "x"})
	if err == nil {
		t.Fatal("expected error when inbox is nil")
	}
}
