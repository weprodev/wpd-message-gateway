// Package authgate provides an adapter between the wpd-message-gateway port.AuthorizationGate
// interface and the wpd-gogate RBAC engine.
//
// The adapter is the only place in the gateway that imports wpd-gogate directly.
// All other packages interact with RBAC through the port.AuthorizationGate interface.
package authgate
