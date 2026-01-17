package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/weprodev/wpd-message-gateway/internal/app"
)

func main() {
	configPath := os.Getenv("CONFIG_PATH")
	cfg, err := app.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err := app.ValidateConfig(cfg); err != nil {
		log.Fatalf("Configuration error: %v\n\n"+
			"Each message type requires a valid default provider (e.g. 'memory', 'mailgun', 'twilio').\n"+
			"Please configure these in configs/local.yml or via environment variables.\n\n"+
			"💡 Tip: Copy configs/local.example.yml to configs/local.yml and configure your providers", err)
	}

	application, err := app.Wire(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	logConfiguration(cfg)

	router := application.Router.Setup()

	port := resolvePort(cfg)
	log.Printf("Gateway server listening on :%s", port)

	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func logConfiguration(cfg *app.Config) {
	log.Printf("Loaded Configuration:")
	log.Printf("- Email Provider: %s (Default)", cfg.DefaultEmailProvider())
	log.Printf("- SMS Provider:   %s (Default)", cfg.DefaultSMSProvider())
	log.Printf("- Push Provider:  %s (Default)", cfg.DefaultPushProvider())
	log.Printf("- Chat Provider:  %s (Default)", cfg.DefaultChatProvider())

	if cfg.Mailpit.Enabled {
		log.Printf("- Mailpit:        enabled")
	}
}

func resolvePort(cfg *app.Config) string {
	if port := os.Getenv("PORT"); port != "" {
		return port
	}
	if cfg.Server.Port != 0 {
		return fmt.Sprintf("%d", cfg.Server.Port)
	}
	return "10101"
}
