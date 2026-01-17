package memory

import (
	"context"
	"testing"

	"github.com/weprodev/wpd-message-gateway/config"
	"github.com/weprodev/wpd-message-gateway/contracts"
)

func TestEmailProvider_SendAndRetrieve(t *testing.T) {
	store := New()
	provider := store.EmailProvider(config.MailpitConfig{})

	email := &contracts.Email{
		To:      []string{"test@example.com"},
		Subject: "Test Subject",
		HTML:    "<p>Test body</p>",
	}

	result, err := provider.Send(context.Background(), email)
	if err != nil {
		t.Fatalf("Send failed: %v", err)
	}

	if result.ID == "" {
		t.Error("Expected non-empty ID")
	}

	if result.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", result.StatusCode)
	}

	// Verify stored
	emails := store.Emails()
	if len(emails) != 1 {
		t.Fatalf("Expected 1 email, got %d", len(emails))
	}

	stored := emails[0]
	if stored.ID != result.ID {
		t.Errorf("ID mismatch: %s != %s", stored.ID, result.ID)
	}
	if stored.Email.Subject != "Test Subject" {
		t.Errorf("Subject mismatch: %s", stored.Email.Subject)
	}
	if stored.CreatedAt.IsZero() {
		t.Error("CreatedAt should not be zero")
	}
}

func TestEmailByID(t *testing.T) {
	store := New()
	provider := store.EmailProvider(config.MailpitConfig{})

	result, _ := provider.Send(context.Background(), &contracts.Email{
		Subject: "Find me",
	})

	found := store.EmailByID(result.ID)
	if found == nil {
		t.Fatal("Expected to find email by ID")
	}
	if found.Email.Subject != "Find me" {
		t.Errorf("Wrong email found")
	}

	notFound := store.EmailByID("non-existent-id")
	if notFound != nil {
		t.Error("Expected nil for non-existent ID")
	}
}

func TestDeleteEmailByID(t *testing.T) {
	store := New()
	provider := store.EmailProvider(config.MailpitConfig{})

	result, _ := provider.Send(context.Background(), &contracts.Email{
		Subject: "Delete me",
	})

	if store.Count() != 1 {
		t.Fatalf("Expected 1 message, got %d", store.Count())
	}

	deleted := store.DeleteEmailByID(result.ID)
	if !deleted {
		t.Error("Expected delete to return true")
	}

	if store.Count() != 0 {
		t.Errorf("Expected 0 messages after delete, got %d", store.Count())
	}

	deletedAgain := store.DeleteEmailByID(result.ID)
	if deletedAgain {
		t.Error("Expected delete to return false for already deleted ID")
	}
}

func TestSMSProvider_SendAndRetrieve(t *testing.T) {
	store := New()
	provider := store.SMSProvider()

	sms := &contracts.SMS{
		From:    "+1234567890",
		To:      []string{"+0987654321"},
		Message: "Test SMS",
	}

	result, err := provider.Send(context.Background(), sms)
	if err != nil {
		t.Fatalf("Send failed: %v", err)
	}

	if result.ID == "" {
		t.Error("Expected non-empty ID")
	}

	allSMS := store.SMS()
	if len(allSMS) != 1 {
		t.Fatalf("Expected 1 SMS, got %d", len(allSMS))
	}

	found := store.SMSByID(result.ID)
	if found == nil || found.SMS.Message != "Test SMS" {
		t.Error("SMS retrieval failed")
	}
}

func TestDeleteSMSByID(t *testing.T) {
	store := New()
	provider := store.SMSProvider()

	result, _ := provider.Send(context.Background(), &contracts.SMS{Message: "Delete me"})

	if !store.DeleteSMSByID(result.ID) {
		t.Error("Expected delete to return true")
	}
	if len(store.SMS()) != 0 {
		t.Error("Expected 0 SMS after delete")
	}
}

func TestPushProvider_SendAndRetrieve(t *testing.T) {
	store := New()
	provider := store.PushProvider()

	push := &contracts.PushNotification{
		DeviceTokens: []string{"token123"},
		Title:        "Test Push",
		Body:         "Test body",
	}

	result, err := provider.Send(context.Background(), push)
	if err != nil {
		t.Fatalf("Send failed: %v", err)
	}

	if result.ID == "" {
		t.Error("Expected non-empty ID")
	}

	allPushes := store.Pushes()
	if len(allPushes) != 1 {
		t.Fatalf("Expected 1 push, got %d", len(allPushes))
	}

	found := store.PushByID(result.ID)
	if found == nil || found.Push.Title != "Test Push" {
		t.Error("Push retrieval failed")
	}
}

func TestDeletePushByID(t *testing.T) {
	store := New()
	provider := store.PushProvider()

	result, _ := provider.Send(context.Background(), &contracts.PushNotification{Title: "Delete me"})

	if !store.DeletePushByID(result.ID) {
		t.Error("Expected delete to return true")
	}
	if len(store.Pushes()) != 0 {
		t.Error("Expected 0 pushes after delete")
	}
}

func TestChatProvider_SendAndRetrieve(t *testing.T) {
	store := New()
	provider := store.ChatProvider()

	chat := &contracts.ChatMessage{
		From:    "user1",
		To:      []string{"user2"},
		Message: "Test chat",
	}

	result, err := provider.Send(context.Background(), chat)
	if err != nil {
		t.Fatalf("Send failed: %v", err)
	}

	if result.ID == "" {
		t.Error("Expected non-empty ID")
	}

	allChats := store.Chats()
	if len(allChats) != 1 {
		t.Fatalf("Expected 1 chat, got %d", len(allChats))
	}

	found := store.ChatByID(result.ID)
	if found == nil || found.Chat.Message != "Test chat" {
		t.Error("Chat retrieval failed")
	}
}

func TestDeleteChatByID(t *testing.T) {
	store := New()
	provider := store.ChatProvider()

	result, _ := provider.Send(context.Background(), &contracts.ChatMessage{Message: "Delete me"})

	if !store.DeleteChatByID(result.ID) {
		t.Error("Expected delete to return true")
	}
	if len(store.Chats()) != 0 {
		t.Error("Expected 0 chats after delete")
	}
}

func TestStats(t *testing.T) {
	store := New()
	mailpitCfg := config.MailpitConfig{}

	_, _ = store.EmailProvider(mailpitCfg).Send(context.Background(), &contracts.Email{Subject: "e1"})
	_, _ = store.EmailProvider(mailpitCfg).Send(context.Background(), &contracts.Email{Subject: "e2"})
	_, _ = store.SMSProvider().Send(context.Background(), &contracts.SMS{Message: "s1"})
	_, _ = store.PushProvider().Send(context.Background(), &contracts.PushNotification{Title: "p1"})
	_, _ = store.ChatProvider().Send(context.Background(), &contracts.ChatMessage{Message: "c1"})

	stats := store.Stats()
	if stats["emails"] != 2 {
		t.Errorf("Expected 2 emails, got %d", stats["emails"])
	}
	if stats["sms"] != 1 {
		t.Errorf("Expected 1 sms, got %d", stats["sms"])
	}
	if stats["push"] != 1 {
		t.Errorf("Expected 1 push, got %d", stats["push"])
	}
	if stats["chat"] != 1 {
		t.Errorf("Expected 1 chat, got %d", stats["chat"])
	}
	if stats["total"] != 5 {
		t.Errorf("Expected 5 total, got %d", stats["total"])
	}
}

func TestClear(t *testing.T) {
	store := New()
	mailpitCfg := config.MailpitConfig{}

	_, _ = store.EmailProvider(mailpitCfg).Send(context.Background(), &contracts.Email{Subject: "e1"})
	_, _ = store.SMSProvider().Send(context.Background(), &contracts.SMS{Message: "s1"})
	_, _ = store.PushProvider().Send(context.Background(), &contracts.PushNotification{Title: "p1"})
	_, _ = store.ChatProvider().Send(context.Background(), &contracts.ChatMessage{Message: "c1"})

	if store.Count() != 4 {
		t.Fatalf("Expected 4 messages, got %d", store.Count())
	}

	store.Clear()

	if store.Count() != 0 {
		t.Errorf("Expected 0 messages after clear, got %d", store.Count())
	}
}
