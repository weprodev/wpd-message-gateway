package presentation

import (
	"log/slog"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"

	"github.com/weprodev/wpd-message-gateway/internal/core/domain"
	"github.com/weprodev/wpd-message-gateway/internal/core/port"
	applogger "github.com/weprodev/wpd-message-gateway/internal/infrastructure/logger"
	"github.com/weprodev/wpd-message-gateway/internal/presentation/handler"
	customMiddleware "github.com/weprodev/wpd-message-gateway/internal/presentation/middleware"

	pkgapi "github.com/weprodev/go-pkg/api"
	pkglogger "github.com/weprodev/go-pkg/logger"
	pkgvalidator "github.com/weprodev/go-pkg/validator"
	gogate "github.com/weprodev/wpd-gogate"
)

// Router holds HTTP handlers and Echo configuration.
type Router struct {
	gatewayHandler           *handler.GatewayHandler
	portalInboxHandler       *handler.PortalInboxHandler
	portalAuthHandler        *handler.PortalAuthHandler
	portalWorkspaceHandler   *handler.PortalWorkspaceHandler
	portalIntegrationHandler *handler.PortalIntegrationHandler
	portalTemplateHandler    *handler.PortalTemplateHandler
	portalHandler            *handler.PortalHandler
	jwtSecret                string
	apiKeyRepo               port.APIKeyRepository
	workspaceRepo            port.WorkspaceRepository
	memberRepo               port.WorkspaceMemberRepository
	logger                   *pkglogger.Logger
	gate                     *gogate.Gate
}

// NewRouter creates a new router bundle.
func NewRouter(
	gateway *handler.GatewayHandler,
	portalInbox *handler.PortalInboxHandler,
	apiKeyRepo port.APIKeyRepository,
	workspaceRepo port.WorkspaceRepository,
	memberRepo port.WorkspaceMemberRepository,
	portalAuth *handler.PortalAuthHandler,
	portalWorkspace *handler.PortalWorkspaceHandler,
	portalIntegration *handler.PortalIntegrationHandler,
	portalTemplate *handler.PortalTemplateHandler,
	portal *handler.PortalHandler,
	jwtSecret string,
	logger *pkglogger.Logger,
	gate *gogate.Gate,
) *Router {
	return &Router{
		gatewayHandler:           gateway,
		portalInboxHandler:       portalInbox,
		portalAuthHandler:        portalAuth,
		portalWorkspaceHandler:   portalWorkspace,
		portalIntegrationHandler: portalIntegration,
		portalTemplateHandler:    portalTemplate,
		portalHandler:            portal,
		jwtSecret:                jwtSecret,
		apiKeyRepo:               apiKeyRepo,
		workspaceRepo:            workspaceRepo,
		memberRepo:               memberRepo,
		logger:                   logger,
		gate:                     gate,
	}
}

// Setup configures Echo with all routes.
func (rt *Router) Setup() *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.HTTPErrorHandler = pkgapi.NewEchoErrorHandler(rt.logger.Logger).EchoHandler
	e.Validator = pkgvalidator.NewValidator()

	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.RequestID())
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			reqID := c.Response().Header().Get(echo.HeaderXRequestID)
			if reqID == "" {
				reqID = c.Request().Header.Get(echo.HeaderXRequestID)
			}
			if reqID != "" {
				ctx := c.Request().Context()
				ctx = applogger.WithRequestID(ctx, reqID)
				c.SetRequest(c.Request().WithContext(ctx))
			}
			return next(c)
		}
	})
	e.Use(echomiddleware.RequestLoggerWithConfig(echomiddleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogMethod:   true,
		LogLatency:  true,
		LogRemoteIP: true,
		LogValuesFunc: func(c echo.Context, v echomiddleware.RequestLoggerValues) error {
			slog.InfoContext(c.Request().Context(), "request summary",
				slog.String("method", v.Method),
				slog.String("uri", v.URI),
				slog.Int("status", v.Status),
				slog.Duration("latency", v.Latency),
				slog.String("remote_ip", v.RemoteIP),
			)
			return nil
		},
	}))
	e.Use(echomiddleware.CORSWithConfig(echomiddleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.PATCH, echo.DELETE, echo.OPTIONS},
		AllowHeaders: []string{
			echo.HeaderAccept,
			echo.HeaderContentType,
			echo.HeaderAuthorization,
			customMiddleware.HeaderWorkspaceKey,
			customMiddleware.HeaderWorkspaceAPIClientID,
			customMiddleware.HeaderWorkspaceAPISecret,
		},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1 := e.Group("/v1")
	v1.Use(customMiddleware.APIKeyAuthMiddleware(rt.apiKeyRepo, rt.workspaceRepo))
	v1.POST("/email", rt.gatewayHandler.HandleSendEmail)
	v1.POST("/sms", rt.gatewayHandler.HandleSendSMS)
	v1.POST("/push", rt.gatewayHandler.HandleSendPush)
	v1.POST("/chat", rt.gatewayHandler.HandleSendChat)

	api := e.Group("/api/v1")

	// Portal API — always enabled.
	pa := rt.portalAuthHandler
	pw := rt.portalWorkspaceHandler
	pi := rt.portalIntegrationHandler
	pt := rt.portalTemplateHandler
	ph := rt.portalHandler

	api.POST("/auth/register", pa.Register)
	api.GET("/auth/verify-email", pa.VerifyEmail)
	api.POST("/auth/login", pa.Login)

	protected := api.Group("")
	protected.Use(customMiddleware.PortalJWT(rt.jwtSecret))
	protected.GET("/auth/me", pa.Me)
	protected.POST("/workspaces/join", pw.JoinWorkspace)
	protected.GET("/workspaces", pw.ListWorkspaces)
	protected.POST("/workspaces", pw.CreateWorkspace)
	protected.GET("/workspaces/:wid", pw.GetWorkspace, customMiddleware.RequirePermission(rt.gate, rt.workspaceRepo, domain.PermissionWorkspacesRead))
	protected.PATCH("/workspaces/:wid", pw.PatchWorkspace, customMiddleware.RequirePermission(rt.gate, rt.workspaceRepo, domain.PermissionWorkspacesWrite))
	protected.GET("/workspaces/:wid/members", pw.ListMembers, customMiddleware.RequirePermission(rt.gate, rt.workspaceRepo, domain.PermissionMembersRead))
	protected.DELETE("/workspaces/:wid/members/:userId", pw.RemoveMember, customMiddleware.RequirePermission(rt.gate, rt.workspaceRepo, domain.PermissionMembersWrite))

	protected.GET("/workspaces/:wid/api-keys", ph.ListAPIKeys, customMiddleware.RequirePermission(rt.gate, rt.workspaceRepo, domain.PermissionAPIKeysRead))
	protected.POST("/workspaces/:wid/api-keys", ph.CreateAPIKey, customMiddleware.RequirePermission(rt.gate, rt.workspaceRepo, domain.PermissionAPIKeysWrite))
	protected.DELETE("/workspaces/:wid/api-keys/:keyId", ph.DeleteAPIKey, customMiddleware.RequirePermission(rt.gate, rt.workspaceRepo, domain.PermissionAPIKeysWrite))
	protected.POST("/workspaces/:wid/api-keys/:keyId/regenerate", ph.RegenerateAPIKey, customMiddleware.RequirePermission(rt.gate, rt.workspaceRepo, domain.PermissionAPIKeysWrite))
	protected.GET("/workspaces/:wid/logs", ph.ListLogs, customMiddleware.RequirePermission(rt.gate, rt.workspaceRepo, domain.PermissionLogsRead))

	protected.GET("/workspaces/:wid/integrations", pi.ListIntegrations, customMiddleware.RequirePermission(rt.gate, rt.workspaceRepo, domain.PermissionIntegrationsRead))
	protected.POST("/workspaces/:wid/integrations", pi.UpsertIntegration, customMiddleware.RequirePermission(rt.gate, rt.workspaceRepo, domain.PermissionIntegrationsWrite))
	protected.DELETE("/workspaces/:wid/integrations/:iid", pi.DeleteIntegration, customMiddleware.RequirePermission(rt.gate, rt.workspaceRepo, domain.PermissionIntegrationsWrite))
	protected.GET("/workspaces/:wid/providers", pi.ListProviders, customMiddleware.RequirePermission(rt.gate, rt.workspaceRepo, domain.PermissionIntegrationsRead))
	protected.GET("/workspaces/:wid/providers/:name/config", pi.GetProviderConfigFields, customMiddleware.RequirePermission(rt.gate, rt.workspaceRepo, domain.PermissionIntegrationsRead))

	protected.GET("/workspaces/:wid/templates", pt.ListTemplates, customMiddleware.RequirePermission(rt.gate, rt.workspaceRepo, domain.PermissionTemplatesRead))
	protected.POST("/workspaces/:wid/templates", pt.CreateTemplate, customMiddleware.RequirePermission(rt.gate, rt.workspaceRepo, domain.PermissionTemplatesWrite))
	protected.PATCH("/workspaces/:wid/templates/:tid", pt.PatchTemplate, customMiddleware.RequirePermission(rt.gate, rt.workspaceRepo, domain.PermissionTemplatesWrite))
	protected.DELETE("/workspaces/:wid/templates/:tid", pt.DeleteTemplate, customMiddleware.RequirePermission(rt.gate, rt.workspaceRepo, domain.PermissionTemplatesWrite))

	protected.GET("/workspaces/:wid/settings", ph.GetSettings, customMiddleware.RequirePermission(rt.gate, rt.workspaceRepo, domain.PermissionSettingsRead))
	protected.PATCH("/workspaces/:wid/settings", ph.PatchSettings, customMiddleware.RequirePermission(rt.gate, rt.workspaceRepo, domain.PermissionSettingsWrite))

	protected.GET("/workspaces/:wid/invitations", pw.ListInvitations, customMiddleware.RequirePermission(rt.gate, rt.workspaceRepo, domain.PermissionInvitationsRead))
	protected.POST("/workspaces/:wid/invitations", pw.CreateInvitation, customMiddleware.RequirePermission(rt.gate, rt.workspaceRepo, domain.PermissionInvitationsWrite))

	// Inbox API — JWT + logs.read (same access model as message logs overview).
	if rt.memberRepo != nil {
		inbox := rt.portalInboxHandler
		inboxGroup := api.Group("/workspaces/:wid/inbox")
		inboxGroup.Use(customMiddleware.PortalJWTBearerOrQuery(rt.jwtSecret))
		inboxGroup.Use(customMiddleware.RequirePermission(rt.gate, rt.workspaceRepo, domain.PermissionLogsRead))
		inboxGroup.GET("/stats", inbox.HandleStats)
		inboxGroup.GET("/emails", inbox.HandleGetEmails)
		inboxGroup.GET("/emails/:id", inbox.HandleGetEmailByID)
		inboxGroup.GET("/sms", inbox.HandleGetSMS)
		inboxGroup.GET("/sms/:id", inbox.HandleGetSMSByID)
		inboxGroup.GET("/push", inbox.HandleGetPush)
		inboxGroup.GET("/push/:id", inbox.HandleGetPushByID)
		inboxGroup.GET("/chat", inbox.HandleGetChat)
		inboxGroup.GET("/chat/:id", inbox.HandleGetChatByID)
		inboxGroup.GET("/events", inbox.HandleSSE)
		inboxGroup.DELETE("/emails/:id", inbox.HandleDeleteEmailByID, customMiddleware.RequirePermission(rt.gate, rt.workspaceRepo, domain.PermissionSendTest))
		inboxGroup.DELETE("/sms/:id", inbox.HandleDeleteSMSByID, customMiddleware.RequirePermission(rt.gate, rt.workspaceRepo, domain.PermissionSendTest))
		inboxGroup.DELETE("/push/:id", inbox.HandleDeletePushByID, customMiddleware.RequirePermission(rt.gate, rt.workspaceRepo, domain.PermissionSendTest))
		inboxGroup.DELETE("/chat/:id", inbox.HandleDeleteChatByID, customMiddleware.RequirePermission(rt.gate, rt.workspaceRepo, domain.PermissionSendTest))
		inboxGroup.DELETE("/messages", inbox.HandleClearAll, customMiddleware.RequirePermission(rt.gate, rt.workspaceRepo, domain.PermissionSendTest))

		internal := api.Group("/workspaces/:wid/internal")
		internal.Use(customMiddleware.PortalJWTBearerOrQuery(rt.jwtSecret))
		internal.Use(customMiddleware.RequirePermission(rt.gate, rt.workspaceRepo, domain.PermissionSendTest))
		internal.POST("/email", inbox.HandleIngestEmail)
		internal.POST("/sms", inbox.HandleIngestSMS)
		internal.POST("/push", inbox.HandleIngestPush)
		internal.POST("/chat", inbox.HandleIngestChat)
	}

	return e
}
