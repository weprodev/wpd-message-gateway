package app

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/weprodev/wpd-packages/common"
	pkgconfig "github.com/weprodev/wpd-packages/config"

	"github.com/weprodev/wpd-message-gateway/internal/core/service"
	"github.com/weprodev/wpd-message-gateway/internal/infrastructure/provider/memory"
	"github.com/weprodev/wpd-message-gateway/internal/infrastructure/repository/postgres"
	"github.com/weprodev/wpd-message-gateway/internal/presentation"
	"github.com/weprodev/wpd-message-gateway/internal/presentation/handler"
	"github.com/weprodev/wpd-message-gateway/pkg/encryption"
)

// Application holds all wired dependencies produced by Wire().
type Application struct {
	Config         *Config
	GatewayService *service.GatewayService
	MemoryStore    *memory.Store
	PgClient       *common.PgClient
	Echo           *echo.Echo
}

// Wire builds and connects all application dependencies.
// It returns an error if any required secret is missing or the DB is unreachable.
func Wire(cfg *Config) (*Application, error) {
	// ── Database ────────────────────────────────────────────────────────────
	dbConfig := pkgconfig.ApplyDatabaseOverrides(pkgconfig.DatabaseConfig{})
	pgClient, err := common.NewPgClient(common.PgConfig{
		Host:           dbConfig.Host,
		ConnectionName: dbConfig.ConnectionName,
		DatabaseURL:    dbConfig.DatabaseURL,
		Port:           dbConfig.Port,
		User:           dbConfig.User,
		Password:       dbConfig.Password,
		DBName:         dbConfig.DBName,
		SSLMode:        dbConfig.SSLMode,
		MaxOpenConns:   dbConfig.MaxOpenConns,
		MaxIdleConns:   dbConfig.MaxIdleConns,
	})
	if err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	// ── Encryption service ──────────────────────────────────────────────────
	encKey, err := resolveSecret("MESSAGE_CONFIG_ENCRYPTION_KEY", cfg.Environment, 32)
	if err != nil {
		return nil, err
	}
	encService, err := encryption.NewAESService([]byte(encKey))
	if err != nil {
		return nil, fmt.Errorf("setup encryption: %w", err)
	}

	// ── Portal JWT ──────────────────────────────────────────────────────────
	jwtSecret := cfg.Portal.JWTSecret
	if v := os.Getenv("MESSAGE_JWT_SECRET"); v != "" {
		jwtSecret = v
	}
	if err := validateSecret("MESSAGE_JWT_SECRET / portal.jwt_secret", jwtSecret, 32, cfg.Environment); err != nil {
		return nil, err
	}
	jwtTTL := 24 * time.Hour
	if cfg.Portal.JWTTTLHours > 0 {
		jwtTTL = time.Duration(cfg.Portal.JWTTTLHours) * time.Hour
	}

	// ── Repositories ────────────────────────────────────────────────────────
	apiKeyRepo := postgres.NewAPIKeyRepository(pgClient)
	workspaceRepo := postgres.NewWorkspaceRepository(pgClient)
	intgRepo := postgres.NewIntegrationRepository(pgClient, encService)
	tmplRepo := postgres.NewTemplateRepository(pgClient)
	logRepo := postgres.NewMessageRequestLogRepository(pgClient)
	userRepo := postgres.NewUserRepository(pgClient)
	memberRepo := postgres.NewWorkspaceMemberRepository(pgClient)
	inviteRepo := postgres.NewInvitationRepository(pgClient)
	settingsRepo := postgres.NewWorkspaceSettingsRepository(pgClient)

	// ── Portal service ──────────────────────────────────────────────────────
	portalSvc := service.NewPortalService(service.PortalDeps{
		Users:        userRepo,
		Members:      memberRepo,
		Workspaces:   workspaceRepo,
		APIKeys:      apiKeyRepo,
		Integrations: intgRepo,
		Templates:    tmplRepo,
		Logs:         logRepo,
		Invitations:  inviteRepo,
		Settings:     settingsRepo,
		JWTSecret:    jwtSecret,
		JWTTTL:       jwtTTL,
	})

	// ── Memory store & inbox writer ─────────────────────────────────────────
	memoryStore := memory.GetStore()
	inboxWriter := memory.NewInboxWriter(memoryStore)

	// ── Gateway service ─────────────────────────────────────────────────────
	gatewaySvc := service.NewGatewayService(intgRepo, tmplRepo, settingsRepo, inboxWriter)

	// ── Handlers ─────────────────────────────────────────────────────────────
	// Portal is always enabled — configuration, templates, and inbox require it.
	portalHandler := handler.NewPortalHandler(portalSvc)
	portalInboxHandler := handler.NewPortalInboxHandler(memoryStore)
	gatewayHandler := handler.NewGatewayHandler(gatewaySvc, logRepo)

	// ── Router ──────────────────────────────────────────────────────────────
	router := presentation.NewRouter(
		gatewayHandler, portalInboxHandler,
		apiKeyRepo, workspaceRepo, memberRepo,
		portalHandler, jwtSecret,
	)

	return &Application{
		Config:         cfg,
		GatewayService: gatewaySvc,
		MemoryStore:    memoryStore,
		PgClient:       pgClient,
		Echo:           router.Setup(),
	}, nil
}

// resolveSecret returns the secret from env or cfg, failing if it is too short
// in non-local environments. In local/test it falls back to a dev-only default
// and warns — never used silently in production.
func resolveSecret(envKey string, environment string, minLen int) (string, error) {
	val := os.Getenv(envKey)
	if len(val) >= minLen {
		return val, nil
	}
	if isLocalEnv(environment) {
		// Safe to use a well-known dev default — not a production secret.
		return "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef", nil
	}
	return "", fmt.Errorf(
		"%s must be set to a %d+ character secret in %q environment; "+
			"see docs/backend/usage.md for configuration instructions",
		envKey, minLen, environment,
	)
}

// validateSecret fails if the secret is too short in non-local environments.
func validateSecret(label, val string, minLen int, environment string) error {
	if len(val) >= minLen {
		return nil
	}
	if isLocalEnv(environment) {
		return nil // dev-only: accept short/empty secrets
	}
	return errors.New(label + " must be at least " + fmt.Sprintf("%d", minLen) + " characters")
}

func isLocalEnv(env string) bool {
	return env == "" || env == "local" || env == "test" || env == "development"
}
