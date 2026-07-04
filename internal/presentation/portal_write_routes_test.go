package presentation

import (
	"net/http"
	"testing"

	"github.com/weprodev/wpd-message-gateway/internal/core/domain"
)

// portalWriteRoutes documents mutating portal endpoints and their required write permissions.
// Keep in sync with Router.Setup — every state-changing workspace route must appear here.
var portalWriteRoutes = []struct {
	method     string
	path       string
	permission string
}{
	{http.MethodPatch, "/api/v1/workspaces/:wid", domain.PermissionWorkspacesWrite},
	{http.MethodDelete, "/api/v1/workspaces/:wid/members/:userId", domain.PermissionMembersWrite},
	{http.MethodPost, "/api/v1/workspaces/:wid/api-keys", domain.PermissionAPIKeysWrite},
	{http.MethodDelete, "/api/v1/workspaces/:wid/api-keys/:keyId", domain.PermissionAPIKeysWrite},
	{http.MethodPost, "/api/v1/workspaces/:wid/api-keys/:keyId/regenerate", domain.PermissionAPIKeysWrite},
	{http.MethodPost, "/api/v1/workspaces/:wid/integrations", domain.PermissionIntegrationsWrite},
	{http.MethodDelete, "/api/v1/workspaces/:wid/integrations/:iid", domain.PermissionIntegrationsWrite},
	{http.MethodPost, "/api/v1/workspaces/:wid/templates", domain.PermissionTemplatesWrite},
	{http.MethodPatch, "/api/v1/workspaces/:wid/templates/:tid", domain.PermissionTemplatesWrite},
	{http.MethodDelete, "/api/v1/workspaces/:wid/templates/:tid", domain.PermissionTemplatesWrite},
	{http.MethodPatch, "/api/v1/workspaces/:wid/settings", domain.PermissionSettingsWrite},
	{http.MethodPost, "/api/v1/workspaces/:wid/invitations", domain.PermissionInvitationsWrite},
	{http.MethodDelete, "/api/v1/workspaces/:wid/inbox/emails/:id", domain.PermissionInboxWrite},
	{http.MethodDelete, "/api/v1/workspaces/:wid/inbox/sms/:id", domain.PermissionInboxWrite},
	{http.MethodDelete, "/api/v1/workspaces/:wid/inbox/push/:id", domain.PermissionInboxWrite},
	{http.MethodDelete, "/api/v1/workspaces/:wid/inbox/chat/:id", domain.PermissionInboxWrite},
	{http.MethodDelete, "/api/v1/workspaces/:wid/inbox/messages", domain.PermissionInboxWrite},
	{http.MethodPost, "/api/v1/workspaces/:wid/internal/email", domain.PermissionInboxWrite},
	{http.MethodPost, "/api/v1/workspaces/:wid/internal/sms", domain.PermissionInboxWrite},
	{http.MethodPost, "/api/v1/workspaces/:wid/internal/push", domain.PermissionInboxWrite},
	{http.MethodPost, "/api/v1/workspaces/:wid/internal/chat", domain.PermissionInboxWrite},
}

func TestPortalWriteRoutes_requireWritePermissions(t *testing.T) {
	t.Parallel()

	for _, route := range portalWriteRoutes {
		if !domain.IsWritePermission(route.permission) {
			t.Fatalf("%s %s must require a write permission, got %q", route.method, route.path, route.permission)
		}
	}
}
