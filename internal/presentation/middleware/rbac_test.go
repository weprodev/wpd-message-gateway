package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	gogate "github.com/weprodev/wpd-gogate"

	"github.com/weprodev/wpd-message-gateway/internal/core/domain"
)

func TestRequirePermission(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close() //nolint:errcheck

	gate := gogate.NewGate(db, nil)

	// Build echo instance
	e := echo.New()

	t.Run("Allowed Access", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/workspaces/ws-1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("wid")
		c.SetParamValues("ws-1")

		// Ingest portal user id in request context
		ctx := context.WithValue(req.Context(), PortalUserIDKey, "user-123")
		c.SetRequest(req.WithContext(ctx))

		// Mock the check query inside Gate.Check
		mock.ExpectQuery(`SELECT 'role' AS type, r.name AS value FROM model_has_roles mhr JOIN roles r ON r.id = mhr.role_id WHERE mhr.model_type = \$1 AND mhr.model_id = \$2 AND mhr.team_id IS NOT DISTINCT FROM \$3 UNION ALL SELECT 'permission' AS type, p.name AS value FROM model_has_permissions mhp JOIN permissions p ON p.id = mhp.permission_id WHERE mhp.model_type = \$1 AND mhp.model_id = \$2 AND mhp.team_id IS NOT DISTINCT FROM \$3 AND p.name = \$4`).
			WithArgs("users", "user-123", "ws-1", domain.PermissionWorkspacesRead).
			WillReturnRows(sqlmock.NewRows([]string{"type", "value"}).AddRow("permission", domain.PermissionWorkspacesRead))

		mw := RequirePermission(gate, domain.PermissionWorkspacesRead)
		handler := mw(func(c echo.Context) error {
			return c.String(http.StatusOK, "OK")
		})

		err := handler(c)
		if err != nil {
			t.Fatalf("expected handler to succeed: %v", err)
		}
		if rec.Code != http.StatusOK {
			t.Errorf("expected 200 OK, got %d", rec.Code)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})

	t.Run("Denied Access", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/workspaces/ws-1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("wid")
		c.SetParamValues("ws-1")

		ctx := context.WithValue(req.Context(), PortalUserIDKey, "user-123")
		c.SetRequest(req.WithContext(ctx))

		mock.ExpectQuery(`SELECT 'role' AS type, r.name AS value FROM model_has_roles mhr JOIN roles r ON r.id = mhr.role_id WHERE mhr.model_type = \$1 AND mhr.model_id = \$2 AND mhr.team_id IS NOT DISTINCT FROM \$3 UNION ALL SELECT 'permission' AS type, p.name AS value FROM model_has_permissions mhp JOIN permissions p ON p.id = mhp.permission_id WHERE mhp.model_type = \$1 AND mhp.model_id = \$2 AND mhp.team_id IS NOT DISTINCT FROM \$3 AND p.name = \$4`).
			WithArgs("users", "user-123", "ws-1", domain.PermissionWorkspacesWrite).
			WillReturnRows(sqlmock.NewRows([]string{"type", "value"})) // no rows matched

		mw := RequirePermission(gate, domain.PermissionWorkspacesWrite)
		handler := mw(func(c echo.Context) error {
			return c.String(http.StatusOK, "OK")
		})

		err := handler(c)
		if err == nil {
			t.Fatal("expected handler to fail with forbidden error")
		}

		he, ok := err.(*echo.HTTPError)
		if !ok {
			t.Fatalf("expected HTTPError, got %T", err)
		}
		if he.Code != http.StatusForbidden {
			t.Errorf("expected 403 Forbidden, got %d", he.Code)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})
}
