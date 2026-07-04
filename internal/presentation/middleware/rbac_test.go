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
	"github.com/weprodev/wpd-message-gateway/internal/core/port"
)

type mockWorkspaceRepository struct {
	port.WorkspaceRepository
	GetByIDFunc func(ctx context.Context, id string) (*domain.Workspace, error)
}

func (m *mockWorkspaceRepository) GetByID(ctx context.Context, id string) (*domain.Workspace, error) {
	return m.GetByIDFunc(ctx, id)
}

func TestRequirePermission(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close() //nolint:errcheck

	gateCfg := gogate.DefaultConfig()
	gateCfg.DefaultGuardName = domain.RBACGuardName
	gate := gogate.NewGate(db, &gateCfg)

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
		mock.ExpectQuery(`SELECT 'role' AS type, r.name AS value FROM model_has_roles mhr JOIN roles r ON r.id = mhr.role_id WHERE mhr.model_type = \$1 AND mhr.model_id = \$2 AND mhr.team_id = \$5 AND r.guard_name = \$3 UNION ALL SELECT 'permission' AS type, p.name AS value FROM model_has_permissions mhp JOIN permissions p ON p.id = mhp.permission_id WHERE mhp.model_type = \$1 AND mhp.model_id = \$2 AND mhp.team_id = \$5 AND p.name = \$4 AND p.guard_name = \$3`).
			WithArgs("users", "user-123", "msg_web", domain.PermissionWorkspacesRead, "ws-1").
			WillReturnRows(sqlmock.NewRows([]string{"type", "value"}).AddRow("permission", domain.PermissionWorkspacesRead))

		wsRepo := &mockWorkspaceRepository{
			GetByIDFunc: func(ctx context.Context, id string) (*domain.Workspace, error) {
				return &domain.Workspace{
					ID:         id,
					Visibility: "private",
				}, nil
			},
		}

		mw := RequirePermission(gate, wsRepo, domain.PermissionWorkspacesRead)
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

		mock.ExpectQuery(`SELECT 'role' AS type, r.name AS value FROM model_has_roles mhr JOIN roles r ON r.id = mhr.role_id WHERE mhr.model_type = \$1 AND mhr.model_id = \$2 AND mhr.team_id = \$5 AND r.guard_name = \$3 UNION ALL SELECT 'permission' AS type, p.name AS value FROM model_has_permissions mhp JOIN permissions p ON p.id = mhp.permission_id WHERE mhp.model_type = \$1 AND mhp.model_id = \$2 AND mhp.team_id = \$5 AND p.name = \$4 AND p.guard_name = \$3`).
			WithArgs("users", "user-123", "msg_web", domain.PermissionWorkspacesWrite, "ws-1").
			WillReturnRows(sqlmock.NewRows([]string{"type", "value"})) // no rows matched

		wsRepo := &mockWorkspaceRepository{
			GetByIDFunc: func(ctx context.Context, id string) (*domain.Workspace, error) {
				return &domain.Workspace{
					ID:         id,
					Visibility: "private",
				}, nil
			},
		}

		mw := RequirePermission(gate, wsRepo, domain.PermissionWorkspacesWrite)
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

	t.Run("Public Workspace Read Access Bypass", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/workspaces/ws-1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("wid")
		c.SetParamValues("ws-1")

		ctx := context.WithValue(req.Context(), PortalUserIDKey, "user-123")
		c.SetRequest(req.WithContext(ctx))

		wsRepo := &mockWorkspaceRepository{
			GetByIDFunc: func(ctx context.Context, id string) (*domain.Workspace, error) {
				return &domain.Workspace{
					ID:         id,
					Visibility: "public",
				}, nil
			},
		}

		// RequirePermission is called with read permission, but no SQL query is expected
		// because the public read bypass logic should authorize the request instantly.
		mw := RequirePermission(gate, wsRepo, domain.PermissionWorkspacesRead)
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
	})

	t.Run("Public Workspace Denies Sensitive Read", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/workspaces/ws-1/logs", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("wid")
		c.SetParamValues("ws-1")

		ctx := context.WithValue(req.Context(), PortalUserIDKey, "user-123")
		c.SetRequest(req.WithContext(ctx))

		wsRepo := &mockWorkspaceRepository{
			GetByIDFunc: func(ctx context.Context, id string) (*domain.Workspace, error) {
				return &domain.Workspace{
					ID:         id,
					Visibility: "public",
				}, nil
			},
		}

		mock.ExpectQuery(`SELECT 'role' AS type, r.name AS value FROM model_has_roles mhr JOIN roles r ON r.id = mhr.role_id WHERE mhr.model_type = \$1 AND mhr.model_id = \$2 AND mhr.team_id = \$5 AND r.guard_name = \$3 UNION ALL SELECT 'permission' AS type, p.name AS value FROM model_has_permissions mhp JOIN permissions p ON p.id = mhp.permission_id WHERE mhp.model_type = \$1 AND mhp.model_id = \$2 AND mhp.team_id = \$5 AND p.name = \$4 AND p.guard_name = \$3`).
			WithArgs("users", "user-123", "msg_web", domain.PermissionLogsRead, "ws-1").
			WillReturnRows(sqlmock.NewRows([]string{"type", "value"}))

		mw := RequirePermission(gate, wsRepo, domain.PermissionLogsRead)
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
