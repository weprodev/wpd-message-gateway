// Email example demonstrates sending an email using the message gateway.
//
// Usage:
//  1. cd examples/email
//  2. go run main.go
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
	// Create gateway with memory provider (for testing)
	gw, err := gateway.New(gateway.Config{
		DefaultEmailProvider: "memory",
		MailpitEnabled:       false,
	})
	if err != nil {
		log.Fatalf("Failed to create gateway: %v", err)
	}

	// Get recipient from args or use default
	recipient := "test@example.com"
	if len(os.Args) > 1 {
		recipient = os.Args[1]
	}

	// Send email
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := gw.SendEmail(ctx, &contracts.Email{
		To:        []string{recipient},
		Subject:   "Test Email from Go Message Gateway",
		HTML:      "<h1>Hello!</h1><p>This email was sent using <strong>wpd-message-gateway</strong>.</p>",
		PlainText: "Hello! This email was sent using wpd-message-gateway.",
	})
	if err != nil {
		log.Fatalf("Failed to send email: %v", err)
	}

	log.Printf("✅ Email sent successfully!")
	log.Printf("   Message ID: %s", result.ID)
	log.Printf("   Status: %s", result.Message)
}
