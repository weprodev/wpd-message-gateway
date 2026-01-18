// Chat example demonstrates sending a chat message using the message gateway.
//
// Usage:
//
//	cd examples/chat
//	go run main.go
package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/weprodev/wpd-message-gateway/pkg/contracts"
	"github.com/weprodev/wpd-message-gateway/pkg/gateway"
)

func main() {
	gw, err := gateway.New(gateway.Config{
		DefaultChatProvider: "memory",
	})
	if err != nil {
		log.Fatalf("Failed to create gateway: %v", err)
	}

	recipient := "+1234567890"
	if len(os.Args) > 1 {
		recipient = os.Args[1]
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := gw.SendChat(ctx, &contracts.ChatMessage{
		To:      []string{recipient},
		Message: "Hello from wpd-message-gateway! This is a chat message.",
	})
	if err != nil {
		log.Fatalf("Failed to send chat message: %v", err)
	}

	log.Printf("✅ Chat message sent successfully!")
	log.Printf("   Message ID: %s", result.ID)
	log.Printf("   Status: %s", result.Message)
}
