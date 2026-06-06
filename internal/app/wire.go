package app

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/labstack/echo/v4"

	pkgconfig "github.com/weprodev/go-pkg/config"
	pkglogger "github.com/weprodev/go-pkg/logger"
	"github.com/weprodev/go-pkg/pgsql"

	"github.com/weprodev/go-pkg/crypto"

	gogate "github.com/weprodev/wpd-gogate"

	"github.com/weprodev/wpd-message-gateway/internal/core/service"
	"github.com/weprodev/wpd-message-gateway/internal/infrastructure/authgate"
	"github.com/weprodev/wpd-message-gateway/internal/infrastructure/provider/memory"
	"github.com/weprodev/wpd-message-gateway/internal/infrastructure/repository/postgres"
	"github.com/weprodev/wpd-message-gateway/internal/presentation"
	"github.com/weprodev/wpd-message-gateway/internal/presentation/handler"
)

// Application holds all wired dependencies produced by Wire().
type Application struct {
	Config         *Config
	AuthService    *service.AuthService
	GatewayService *service.GatewayService
	MemoryStore    *memory.Store
	PgClient       *pgsql.PgClient
	Echo           *echo.Echo
}

// Wire builds and connects all application dependencies.
// It returns an error if any required secret is missing or the DB is unreachable.
func Wire(cfg *Config, sysLogger *pkglogger.Logger) (*Application, error) {
	// Startup context: cap initialization time to avoid hanging indefinitely.
	startCtx, startCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer startCancel()

	// ── Database ────────────────────────────────────────────────────────────
	dbConfig := pkgconfig.ApplyDatabaseOverrides(pkgconfig.DatabaseConfig{})
	pgClient, err := pgsql.NewPgClient(pgsql.PgConfig{
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

	// ── RBAC Gate ───────────────────────────────────────────────────────────
	gateEngine := gogate.NewGate(pgClient.GetDB(startCtx), nil)
	if err := gateEngine.LoadPolicy(startCtx); err != nil {
		return nil, fmt.Errorf("load RBAC policies: %w", err)
	}
	authGate := authgate.NewGoGateAdapter(gateEngine)

	// ── Encryption service ──────────────────────────────────────────────────
	encKey, err := resolveSecret("MESSAGE_CONFIG_ENCRYPTION_KEY", cfg.Environment, 32)
	if err != nil {
		return nil, err
	}
	encService, err := crypto.NewAESService([]byte(encKey))
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
		Gate:         authGate,
	})

	// ── Auth service ──────────────────────────────────────────────────────
	authService := service.NewAuthService(
		userRepo,
		memory.NewEmailProvider(memory.GetStore()),
		jwtSecret,
		true,
		jwtTTL,
		24*time.Hour,
		cfg.Portal.BaseURL,
	)

	// ── Memory store & inbox writer ─────────────────────────────────────────
	memoryStore := memory.GetStore()
	inboxWriter := memory.NewInboxWriter(memoryStore)

	// ── Gateway service ─────────────────────────────────────────────────────
	gatewaySvc := service.NewGatewayService(intgRepo, tmplRepo, settingsRepo, inboxWriter)

	// ── Handlers ─────────────────────────────────────────────────────────────
	// Portal is always enabled — configuration, templates, and inbox require it.
	portalHandler := handler.NewPortalHandler(portalSvc, authService, gatewaySvc, logRepo)
	portalInboxHandler := handler.NewPortalInboxHandler(memoryStore)
	gatewayHandler := handler.NewGatewayHandler(gatewaySvc, logRepo)

	// ── Router ──────────────────────────────────────────────────────────────
	router := presentation.NewRouter(
		gatewayHandler, portalInboxHandler,
		apiKeyRepo, workspaceRepo, memberRepo,
		portalHandler, jwtSecret, sysLogger,
		gateEngine,
	)

	return &Application{
		Config:         cfg,
		AuthService:    authService,
		GatewayService: gatewaySvc,
		MemoryStore:    memoryStore,
		PgClient:       pgClient,
		Echo:           router.Setup(),
	}, nil
}

// resolveSecret returns the secret from the environment, failing if it is missing or too short.
func resolveSecret(envKey string, environment string, minLen int) (string, error) {
	val := os.Getenv(envKey)
	if len(val) >= minLen {
		return val, nil
	}
	return "", fmt.Errorf(
		"%s must be set to a %d+ character secret; "+
			"see docs/backend/usage.md or the .env file for configuration instructions",
		envKey, minLen,
	)
}

// validateSecret fails if the secret is missing or too short.
func validateSecret(label, val string, minLen int, environment string) error {
	if len(val) >= minLen {
		return nil
	}
	return fmt.Errorf("%s must be at least %d characters", label, minLen)
}
