package domain

import "testing"

func TestIsWorkspaceRole(t *testing.T) {
	t.Parallel()

	for _, role := range WorkspaceRoles {
		if !IsWorkspaceRole(role) {
			t.Fatalf("IsWorkspaceRole(%q) = false, want true", role)
		}
	}
	if IsWorkspaceRole("superadmin") {
		t.Fatal("IsWorkspaceRole(superadmin) = true, want false")
	}
}

func TestIsPublicGuestPermission(t *testing.T) {
	t.Parallel()

	if !IsPublicGuestPermission(PermissionWorkspacesRead) {
		t.Fatal("expected workspaces.read for public guests")
	}
	if !IsPublicGuestPermission(PermissionTemplatesRead) {
		t.Fatal("expected templates.read for public guests")
	}
	if IsPublicGuestPermission(PermissionLogsRead) {
		t.Fatal("logs.read must not be granted to public guests")
	}
}
