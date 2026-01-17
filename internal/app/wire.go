package app

import (
	"fmt"
	"log"
	"strings"

	"github.com/weprodev/wpd-message-gateway/internal/core/service"
	"github.com/weprodev/wpd-message-gateway/internal/infrastructure/provider/memory"
	"github.com/weprodev/wpd-message-gateway/internal/presentation"
	"github.com/weprodev/wpd-message-gateway/internal/presentation/handler"
)

// Application holds all wired dependencies.
type Application struct {
	Config         *Config
	GatewayService *service.GatewayService
	MemoryStore    *memory.Store
	Router         *presentation.Router
}

// Wire creates and wires all application dependencies.
func Wire(cfg *Config) (*Application, error) {
	// Create memory store (shared across all memory providers)
	memoryStore := memory.NewStore()

	// Create registry
	registry := service.NewRegistry()

	// Create provider factory
	factory := NewProviderFactory(cfg, memoryStore)

	// Initialize default providers
	if err := initializeDefaultProviders(cfg, factory, registry); err != nil {
		return nil, fmt.Errorf("failed to initialize providers: %w", err)
	}

	// Create gateway service
	gatewaySvc := service.NewGatewayService(cfg, registry)

	// Create handlers
	gatewayHandler := handler.NewGatewayHandler(gatewaySvc, memoryStore)

	var devboxHandler *handler.DevBoxHandler
	if cfg.DevBox.Enabled || cfg.Providers.Defaults.Email == "memory" {
		mailpitCfg := memory.MailpitConfig{Enabled: cfg.Mailpit.Enabled}
		devboxHandler = handler.NewDevBoxHandler(memoryStore, mailpitCfg)
	}

	// Create router
	router := presentation.NewRouter(gatewayHandler, devboxHandler)

	return &Application{
		Config:         cfg,
		GatewayService: gatewaySvc,
		MemoryStore:    memoryStore,
		Router:         router,
	}, nil
}

// MustWire creates application or panics.
func MustWire(cfg *Config) *Application {
	app, err := Wire(cfg)
	if err != nil {
		panic(fmt.Sprintf("failed to wire application: %v", err))
	}
	return app
}

// initializeDefaultProviders initializes providers set as defaults.
func initializeDefaultProviders(cfg *Config, factory *ProviderFactory, registry *service.Registry) error {
	// Initialize default email provider
	if name := cfg.DefaultEmailProvider(); name != "" {
		provider, err := factory.CreateEmailProvider(name)
		if err != nil && !isUnknownProviderError(err) {
			return fmt.Errorf("failed to initialize email provider %s: %w", name, err)
		}
		if provider != nil {
			registry.RegisterEmailProvider(name, provider)
			log.Printf("Registered email provider: %s", name)
		}
	}

	// Initialize default SMS provider
	if name := cfg.DefaultSMSProvider(); name != "" {
		provider, err := factory.CreateSMSProvider(name)
		if err != nil && !isUnknownProviderError(err) {
			return fmt.Errorf("failed to initialize SMS provider %s: %w", name, err)
		}
		if provider != nil {
			registry.RegisterSMSProvider(name, provider)
			log.Printf("Registered SMS provider: %s", name)
		}
	}

	// Initialize default push provider
	if name := cfg.DefaultPushProvider(); name != "" {
		provider, err := factory.CreatePushProvider(name)
		if err != nil && !isUnknownProviderError(err) {
			return fmt.Errorf("failed to initialize push provider %s: %w", name, err)
		}
		if provider != nil {
			registry.RegisterPushProvider(name, provider)
			log.Printf("Registered push provider: %s", name)
		}
	}

	// Initialize default chat provider
	if name := cfg.DefaultChatProvider(); name != "" {
		provider, err := factory.CreateChatProvider(name)
		if err != nil && !isUnknownProviderError(err) {
			return fmt.Errorf("failed to initialize chat provider %s: %w", name, err)
		}
		if provider != nil {
			registry.RegisterChatProvider(name, provider)
			log.Printf("Registered chat provider: %s", name)
		}
	}

	return nil
}

// isUnknownProviderError checks if error indicates an unknown provider.
func isUnknownProviderError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return strings.Contains(errStr, "no configuration found") || strings.Contains(errStr, "unknown")
}
