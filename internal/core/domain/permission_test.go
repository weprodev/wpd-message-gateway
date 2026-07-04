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
