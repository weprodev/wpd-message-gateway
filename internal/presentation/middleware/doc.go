// Package middleware contains Echo middleware functions used by the message gateway HTTP server.
//
// Auth middlewares: APIKeyAuthMiddleware (API key validation for /v1/ routes),
// PortalJWT and PortalJWTBearerOrQuery (portal JWT validation), RequireWorkspaceMember,
// RequireWorkspaceAPIKey, RequirePermission (RBAC).
//
// All middleware must read from the request context only — never from the database directly.
// Business logic belongs in the service layer, not here.
package middleware
