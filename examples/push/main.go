// Push example demonstrates sending a push notification using the message gateway.
//
// Usage:
//
//	cd examples/push
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
		DefaultPushProvider: "memory",
	})
	if err != nil {
		log.Fatalf("Failed to create gateway: %v", err)
	}

	deviceToken := "example-device-token"
	if len(os.Args) > 1 {
		deviceToken = os.Args[1]
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := gw.SendPush(ctx, &contracts.PushNotification{
		DeviceTokens: []string{deviceToken},
		Title:        "New Message",
		Body:         "You have a new notification from wpd-message-gateway!",
		Data: map[string]string{
			"action": "open_app",
			"screen": "notifications",
		},
	})
	if err != nil {
		log.Fatalf("Failed to send push notification: %v", err)
	}

	log.Printf("✅ Push notification sent successfully!")
	log.Printf("   Message ID: %s", result.ID)
	log.Printf("   Status: %s", result.Message)
}
