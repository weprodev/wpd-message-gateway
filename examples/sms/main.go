// SMS example demonstrates sending an SMS using the message gateway.
//
// Usage:
//
//	cd examples/sms
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
		DefaultSMSProvider: "kavenegar",
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

	result, err := gw.SendSMS(ctx, &contracts.SMS{
		From:    "2000660110",
		To:      []string{recipient},
		Message: "Hello from wpd-message-gateway via Kavenegar! Your verification code is 123456.",
	})
	if err != nil {
		log.Fatalf("Failed to send SMS: %v", err)
	}

	log.Printf("✅ SMS sent successfully!")
	log.Printf("   Message ID: %s", result.ID)
	log.Printf("   Status: %s", result.Message)
}
