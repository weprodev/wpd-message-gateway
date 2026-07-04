package domain

import (
	"slices"
	"testing"
)

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

func TestReadPermissions_excludeWritePermissions(t *testing.T) {
	t.Parallel()

	for _, readPermission := range ReadPermissions {
		if IsWritePermission(readPermission) {
			t.Fatalf("read permission %q must not be a write permission", readPermission)
		}
	}
}

func TestWritePermissions_areNotReadPermissions(t *testing.T) {
	t.Parallel()

	for _, writePermission := range WritePermissions {
		if IsReadPermission(writePermission) {
			t.Fatalf("write permission %q must not be a read permission", writePermission)
		}
	}
}

func TestAllPermissions_coversPortalCatalog(t *testing.T) {
	t.Parallel()

	if len(AllPermissions) != len(ReadPermissions)+len(WritePermissions) {
		t.Fatalf("AllPermissions length = %d, want %d", len(AllPermissions), len(ReadPermissions)+len(WritePermissions))
	}

	seen := make(map[string]struct{}, len(AllPermissions))
	for _, permission := range AllPermissions {
		if _, ok := seen[permission]; ok {
			t.Fatalf("duplicate permission %q in AllPermissions", permission)
		}
		seen[permission] = struct{}{}
	}
}

func TestViewerReadPermissions_matchSeedCatalog(t *testing.T) {
	t.Parallel()

	expected := []string{
		PermissionWorkspacesRead,
		PermissionMembersRead,
		PermissionAPIKeysRead,
		PermissionLogsRead,
		PermissionIntegrationsRead,
		PermissionTemplatesRead,
		PermissionSettingsRead,
		PermissionInvitationsRead,
	}

	if len(ReadPermissions) != len(expected) {
		t.Fatalf("ReadPermissions length = %d, want %d", len(ReadPermissions), len(expected))
	}

	for _, permission := range expected {
		if !slices.Contains(ReadPermissions, permission) {
			t.Fatalf("ReadPermissions missing %q (viewer seed catalog)", permission)
		}
	}
}
