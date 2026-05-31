package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/weprodev/wpd-message-gateway/internal/app"
	"github.com/weprodev/wpd-message-gateway/internal/infrastructure/logger"
)

func main() {
	// Initialise the structured logger first so every subsequent slog call is
	// formatted and levelled correctly. Uses APP_ENVIRONMENT to pick text (local)
	// vs JSON (production) format.
	env := os.Getenv("APP_ENVIRONMENT")
	sysLogger, err := logger.New(env)
	if err != nil {
		// Logger not yet up — fall back to bare stderr before exiting.
		fmt.Fprintf(os.Stderr, "FATAL: init logger: %v\n", err)
		os.Exit(1)
	}

	configPath := os.Getenv("CONFIG_PATH")
	cfg, err := app.LoadConfig(configPath)
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	if err := app.ValidateConfig(cfg); err != nil {
		slog.Error("configuration error",
			"error", err,
			"hint", "each message type requires a default provider (e.g. memory, mailgun); "+
				"copy configs/local.example.yml to configs/local.yml and configure providers",
		)
		os.Exit(1)
	}

	application, err := app.Wire(cfg, sysLogger)
	if err != nil {
		slog.Error("failed to initialise application", "error", err)
		os.Exit(1)
	}

	logConfiguration(cfg)

	port := resolvePort(cfg)
	slog.Info("gateway server starting", "port", port, "environment", cfg.Environment)

	if err := application.Echo.Start(":" + port); err != nil {
		slog.Error("server stopped", "error", err)
		os.Exit(1)
	}
}

func logConfiguration(cfg *app.Config) {
	slog.Info("loaded configuration",
		"email_provider", cfg.DefaultEmailProvider(),
		"sms_provider", cfg.DefaultSMSProvider(),
		"push_provider", cfg.DefaultPushProvider(),
		"chat_provider", cfg.DefaultChatProvider(),
		"otp_provider", cfg.DefaultOTPProvider(),
	)
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
