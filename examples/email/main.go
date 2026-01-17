// Email example demonstrates sending an email using go-message-provider.
//
// Usage:
//  1. make sure you are in the examples/email directory (cd examples/email)
//  2. Copy .env.mailgun.example to .env (or another provider's example)
//  3. Fill in your credentials
//  4. Run: go run main.go
package main

import (
	"bufio"
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/weprodev/wpd-message-gateway/config"
	"github.com/weprodev/wpd-message-gateway/contracts"
	"github.com/weprodev/wpd-message-gateway/manager"
)

func main() {
	// Load .env file from current directory (optional, falls back to system env)
	if err := loadEnvFile(".env"); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Load configuration
	cfg, err := config.LoadConfig("")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create manager
	mgr, err := manager.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create manager: %v", err)
	}

	// Show available providers
	providers := mgr.AvailableEmailProviders()
	if len(providers) == 0 {
		log.Fatal("No email providers configured. Copy .env.mailgun.example to .env and add your credentials.")
	}
	log.Printf("Available providers: %v", providers)

	// Get recipient from args or use default
	recipient := "test@example.com"
	if len(os.Args) > 1 {
		recipient = os.Args[1]
	}

	// Send email
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := mgr.SendEmail(ctx, &contracts.Email{
		To:        []string{recipient},
		Subject:   "Test Email from Go Message Provider",
		HTML:      "<h1>Hello!</h1><p>This email was sent using <strong>go-message-provider</strong>.</p>",
		PlainText: "Hello! This email was sent using go-message-provider.",
	})
	if err != nil {
		log.Fatalf("Failed to send email: %v", err)
	}

	log.Printf("✅ Email sent successfully!")
	log.Printf("Default Email Provider: %s", cfg.DefaultEmailProvider())
	log.Printf("   Message ID: %s", result.ID)
}

// loadEnvFile loads environment variables from a file (simple .env loader).
func loadEnvFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if idx := strings.Index(line, "="); idx != -1 {
			key := strings.TrimSpace(line[:idx])
			value := strings.TrimSpace(line[idx+1:])
			// Remove quotes
			if len(value) >= 2 && ((value[0] == '"' && value[len(value)-1] == '"') ||
				(value[0] == '\'' && value[len(value)-1] == '\'')) {
				value = value[1 : len(value)-1]
			}
			if os.Getenv(key) == "" {
				_ = os.Setenv(key, value)
			}
		}
	}
	return scanner.Err()
}
