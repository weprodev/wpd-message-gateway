package authgate

import (
	"context"

	gogate "github.com/weprodev/wpd-gogate"

	"github.com/weprodev/wpd-message-gateway/internal/core/port"
)

type GoGateAdapter struct {
	gate *gogate.Gate
}

func NewGoGateAdapter(gate *gogate.Gate) port.AuthorizationGate {
	return &GoGateAdapter{gate: gate}
}

func (a *GoGateAdapter) AssignRole(ctx context.Context, modelType, modelID, teamID, roleName string) error {
	return a.gate.Model(modelType, modelID, teamID).AssignRole(ctx, roleName)
}

func (a *GoGateAdapter) RemoveRole(ctx context.Context, modelType, modelID, teamID, roleName string) error {
	return a.gate.Model(modelType, modelID, teamID).RemoveRole(ctx, roleName)
}

func (a *GoGateAdapter) GetRoleNames(ctx context.Context, modelType, modelID, teamID string) ([]string, error) {
	return a.gate.Model(modelType, modelID, teamID).GetRoleNames(ctx)
}

func (a *GoGateAdapter) GetAllPermissions(ctx context.Context, modelType, modelID, teamID string) ([]string, error) {
	return a.gate.Model(modelType, modelID, teamID).GetAllPermissions(ctx)
}
