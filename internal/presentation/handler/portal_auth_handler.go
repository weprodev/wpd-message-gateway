package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/weprodev/wpd-message-gateway/internal/core/domain"
	"github.com/weprodev/wpd-message-gateway/internal/core/port"
	"github.com/weprodev/wpd-message-gateway/internal/core/service"
	customMiddleware "github.com/weprodev/wpd-message-gateway/internal/presentation/middleware"
)

type PortalAuthHandler struct {
	svc     *service.PortalService
	authSvc *service.AuthService
}

func NewPortalAuthHandler(svc *service.PortalService, authSvc *service.AuthService) *PortalAuthHandler {
	return &PortalAuthHandler{svc: svc, authSvc: authSvc}
}

type registerBody struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type loginBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type tokenResponse struct {
	Token      string             `json:"token"`
	User       *domain.User       `json:"user"`
	Workspaces []domain.Workspace `json:"workspaces"`
}

type userProfileResponse struct {
	*domain.User
	Workspaces []domain.Workspace `json:"workspaces"`
}

func (h *PortalAuthHandler) Register(c echo.Context) error {
	var body registerBody
	if err := c.Bind(&body); err != nil {
		slog.ErrorContext(c.Request().Context(), "failed to bind register body", "error", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid JSON")
	}
	u, err := h.authSvc.RegisterUser(c.Request().Context(), body.FirstName, body.LastName, body.Email, body.Password)
	if err != nil {
		slog.ErrorContext(c.Request().Context(), "user registration failed", "error", err, "email", body.Email)
		return safeHTTPError(err, http.StatusBadRequest, "registration failed")
	}
	return c.JSON(http.StatusCreated, u)
}

func (h *PortalAuthHandler) VerifyEmail(c echo.Context) error {
	token := c.QueryParam("token")
	if token == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "missing token")
	}

	if err := h.authSvc.VerifyEmail(c.Request().Context(), token); err != nil {
		slog.ErrorContext(c.Request().Context(), "email verification failed", "error", err)
		return safeHTTPError(err, http.StatusBadRequest, "verification failed")
	}

	return c.NoContent(http.StatusOK)
}

func (h *PortalAuthHandler) Login(c echo.Context) error {
	var body loginBody
	if err := c.Bind(&body); err != nil {
		slog.ErrorContext(c.Request().Context(), "failed to bind login body", "error", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid JSON")
	}

	u, token, err := h.authSvc.Login(c.Request().Context(), body.Email, body.Password)
	if err != nil {
		slog.ErrorContext(c.Request().Context(), "login failed", "error", err, "email", body.Email)
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid email or password")
	}

	workspaces, err := h.svc.ListWorkspaces(c.Request().Context(), u.ID.String())
	if err != nil {
		slog.ErrorContext(c.Request().Context(), "failed to list workspaces during login", "error", err, "user_id", u.ID)
		workspaces = []domain.Workspace{}
	}

	return c.JSON(http.StatusOK, tokenResponse{Token: token, User: u, Workspaces: workspaces})
}

func (h *PortalAuthHandler) Me(c echo.Context) error {
	uid := customMiddleware.GetPortalUserID(c.Request().Context())
	u, err := h.svc.UserByID(c.Request().Context(), uid)
	if err != nil {
		slog.ErrorContext(c.Request().Context(), "failed to fetch user profile", "error", err, "user_id", uid)
		if errors.Is(err, port.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "user not found")
		}
		return safeHTTPError(err, http.StatusInternalServerError, "failed to load user")
	}

	workspaces, err := h.svc.ListWorkspaces(c.Request().Context(), uid)
	if err != nil {
		slog.ErrorContext(c.Request().Context(), "failed to list workspaces for profile", "error", err, "user_id", uid)
		workspaces = []domain.Workspace{}
	}

	return c.JSON(http.StatusOK, userProfileResponse{User: u, Workspaces: workspaces})
}
