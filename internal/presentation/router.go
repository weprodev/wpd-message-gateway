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
	gatewayHandler     *handler.GatewayHandler
	portalInboxHandler *handler.PortalInboxHandler
	portalHandler      *handler.PortalHandler
	jwtSecret          string
	apiKeyRepo         port.APIKeyRepository
	workspaceRepo      port.WorkspaceRepository
	memberRepo         port.WorkspaceMemberRepository
	logger             *pkglogger.Logger
	gate               *gogate.Gate
}

// NewRouter creates a new router bundle.
func NewRouter(
	gateway *handler.GatewayHandler,
	portalInbox *handler.PortalInboxHandler,
	apiKeyRepo port.APIKeyRepository,
	workspaceRepo port.WorkspaceRepository,
	memberRepo port.WorkspaceMemberRepository,
	portal *handler.PortalHandler,
	jwtSecret string,
	logger *pkglogger.Logger,
	gate *gogate.Gate,
) *Router {
	return &Router{
		gatewayHandler:     gateway,
		portalInboxHandler: portalInbox,
		portalHandler:      portal,
		jwtSecret:          jwtSecret,
		apiKeyRepo:         apiKeyRepo,
		workspaceRepo:      workspaceRepo,
		memberRepo:         memberRepo,
		logger:             logger,
		gate:               gate,
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
			customMiddleware.HeaderInternalSecret,
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
	ph := rt.portalHandler
	api.POST("/auth/register", ph.Register)
	api.GET("/auth/verify-email", ph.VerifyEmail)
	api.POST("/auth/login", ph.Login)

	protected := api.Group("")
	protected.Use(customMiddleware.PortalJWT(rt.jwtSecret))
	protected.GET("/auth/me", ph.Me)
	protected.POST("/workspaces/join", ph.JoinWorkspace)
	protected.GET("/workspaces", ph.ListWorkspaces)
	protected.POST("/workspaces", ph.CreateWorkspace)
	protected.GET("/workspaces/:wid", ph.GetWorkspace, customMiddleware.RequirePermission(rt.gate, domain.PermissionWorkspacesRead))
	protected.PATCH("/workspaces/:wid", ph.PatchWorkspace, customMiddleware.RequirePermission(rt.gate, domain.PermissionWorkspacesWrite))
	protected.GET("/workspaces/:wid/members", ph.ListMembers, customMiddleware.RequirePermission(rt.gate, domain.PermissionMembersRead))
	protected.DELETE("/workspaces/:wid/members/:userId", ph.RemoveMember, customMiddleware.RequirePermission(rt.gate, domain.PermissionMembersWrite))
	protected.GET("/workspaces/:wid/api-keys", ph.ListAPIKeys, customMiddleware.RequirePermission(rt.gate, domain.PermissionAPIKeysRead))
	protected.POST("/workspaces/:wid/api-keys", ph.CreateAPIKey, customMiddleware.RequirePermission(rt.gate, domain.PermissionAPIKeysWrite))
	protected.DELETE("/workspaces/:wid/api-keys/:keyId", ph.DeleteAPIKey, customMiddleware.RequirePermission(rt.gate, domain.PermissionAPIKeysWrite))
	protected.POST("/workspaces/:wid/api-keys/:keyId/regenerate", ph.RegenerateAPIKey, customMiddleware.RequirePermission(rt.gate, domain.PermissionAPIKeysWrite))
	protected.GET("/workspaces/:wid/logs", ph.ListLogs, customMiddleware.RequirePermission(rt.gate, domain.PermissionLogsRead))
	protected.GET("/workspaces/:wid/integrations", ph.ListIntegrations, customMiddleware.RequirePermission(rt.gate, domain.PermissionIntegrationsRead))
	protected.POST("/workspaces/:wid/integrations", ph.UpsertIntegration, customMiddleware.RequirePermission(rt.gate, domain.PermissionIntegrationsWrite))
	protected.DELETE("/workspaces/:wid/integrations/:iid", ph.DeleteIntegration, customMiddleware.RequirePermission(rt.gate, domain.PermissionIntegrationsWrite))
	protected.GET("/workspaces/:wid/templates", ph.ListTemplates, customMiddleware.RequirePermission(rt.gate, domain.PermissionTemplatesRead))
	protected.POST("/workspaces/:wid/templates", ph.CreateTemplate, customMiddleware.RequirePermission(rt.gate, domain.PermissionTemplatesWrite))
	protected.PATCH("/workspaces/:wid/templates/:tid", ph.PatchTemplate, customMiddleware.RequirePermission(rt.gate, domain.PermissionTemplatesWrite))
	protected.DELETE("/workspaces/:wid/templates/:tid", ph.DeleteTemplate, customMiddleware.RequirePermission(rt.gate, domain.PermissionTemplatesWrite))
	protected.GET("/workspaces/:wid/settings", ph.GetSettings, customMiddleware.RequirePermission(rt.gate, domain.PermissionSettingsRead))
	protected.PATCH("/workspaces/:wid/settings", ph.PatchSettings, customMiddleware.RequirePermission(rt.gate, domain.PermissionSettingsWrite))
	protected.GET("/workspaces/:wid/invitations", ph.ListInvitations, customMiddleware.RequirePermission(rt.gate, domain.PermissionInvitationsRead))
	protected.POST("/workspaces/:wid/invitations", ph.CreateInvitation, customMiddleware.RequirePermission(rt.gate, domain.PermissionInvitationsWrite))
	protected.POST("/workspaces/:wid/send-test/:channel", ph.SendTest, customMiddleware.RequirePermission(rt.gate, domain.PermissionSendTest))

	// Inbox API — always enabled (portal is always on).
	if rt.memberRepo != nil {
		inbox := rt.portalInboxHandler
		inboxGroup := api.Group("/workspaces/:wid/inbox")
		inboxGroup.Use(customMiddleware.PortalJWTBearerOrQuery(rt.jwtSecret))
		inboxGroup.Use(customMiddleware.RequireWorkspaceMember(rt.memberRepo))
		inboxGroup.Use(customMiddleware.RequireWorkspaceAPIKey(rt.apiKeyRepo))
		inboxGroup.GET("/stats", inbox.HandleStats)
		inboxGroup.GET("/emails", inbox.HandleGetEmails)
		inboxGroup.GET("/emails/:id", inbox.HandleGetEmailByID)
		inboxGroup.DELETE("/emails/:id", inbox.HandleDeleteEmailByID)
		inboxGroup.GET("/sms", inbox.HandleGetSMS)
		inboxGroup.GET("/sms/:id", inbox.HandleGetSMSByID)
		inboxGroup.DELETE("/sms/:id", inbox.HandleDeleteSMSByID)
		inboxGroup.GET("/push", inbox.HandleGetPush)
		inboxGroup.GET("/push/:id", inbox.HandleGetPushByID)
		inboxGroup.DELETE("/push/:id", inbox.HandleDeletePushByID)
		inboxGroup.GET("/chat", inbox.HandleGetChat)
		inboxGroup.GET("/chat/:id", inbox.HandleGetChatByID)
		inboxGroup.DELETE("/chat/:id", inbox.HandleDeleteChatByID)
		inboxGroup.DELETE("/messages", inbox.HandleClearAll)
		inboxGroup.GET("/events", inbox.HandleSSE)

		internal := api.Group("/workspaces/:wid/internal")
		internal.Use(customMiddleware.InternalIngestSecret())
		internal.POST("/email", inbox.HandleIngestEmail)
		internal.POST("/sms", inbox.HandleIngestSMS)
		internal.POST("/push", inbox.HandleIngestPush)
		internal.POST("/chat", inbox.HandleIngestChat)
	}

	return e
}
