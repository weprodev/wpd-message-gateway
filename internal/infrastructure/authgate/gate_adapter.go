package authgate

import (
	"context"
	"sort"

	"github.com/lib/pq"

	"github.com/weprodev/go-pkg/pgsql"
	gogate "github.com/weprodev/wpd-gogate"

	"github.com/weprodev/wpd-message-gateway/internal/core/port"
)

type GoGateAdapter struct {
	gate *gogate.Gate
	db   pgsql.DBTX
}

func NewGoGateAdapter(gate *gogate.Gate, db pgsql.DBTX) port.AuthorizationGate {
	return &GoGateAdapter{gate: gate, db: db}
}

func (a *GoGateAdapter) AssignRole(ctx context.Context, modelType, modelID, teamID, roleName string) error {
	return a.gate.Model(modelType, modelID, teamID).AssignRole(ctx, roleName, "")
}

func (a *GoGateAdapter) RemoveRole(ctx context.Context, modelType, modelID, teamID, roleName string) error {
	return a.gate.Model(modelType, modelID, teamID).RemoveRole(ctx, roleName, "")
}

func (a *GoGateAdapter) GetRoleNames(ctx context.Context, modelType, modelID, teamID string) ([]string, error) {
	return a.gate.Model(modelType, modelID, teamID).GetRoleNames(ctx)
}

func (a *GoGateAdapter) GetAllPermissions(ctx context.Context, modelType, modelID, teamID string) ([]string, error) {
	return a.gate.Model(modelType, modelID, teamID).GetAllPermissions(ctx)
}

func (a *GoGateAdapter) GetPermissionsForTeams(
	ctx context.Context,
	modelType, modelID string,
	teamIDs []string,
) (map[string][]string, error) {
	if len(teamIDs) == 0 {
		return map[string][]string{}, nil
	}

	query := `
		SELECT team_id::text, permission FROM (
			SELECT mhr.team_id, p.name AS permission
			FROM model_has_roles mhr
			JOIN roles r ON r.id = mhr.role_id
			JOIN role_has_permissions rhp ON rhp.role_id = r.id
			JOIN permissions p ON p.id = rhp.permission_id
			WHERE mhr.model_type = $1 AND mhr.model_id::text = $2 AND mhr.team_id = ANY($3::uuid[])
			UNION
			SELECT mhp.team_id, p.name AS permission
			FROM model_has_permissions mhp
			JOIN permissions p ON p.id = mhp.permission_id
			WHERE mhp.model_type = $1 AND mhp.model_id::text = $2 AND mhp.team_id = ANY($3::uuid[])
		) AS combined
		ORDER BY team_id, permission
	`

	rows, err := a.db.QueryContext(ctx, query, modelType, modelID, pq.Array(teamIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close() //nolint:errcheck

	out := make(map[string][]string, len(teamIDs))
	for rows.Next() {
		var teamID, permission string
		if err := rows.Scan(&teamID, &permission); err != nil {
			return nil, err
		}
		out[teamID] = append(out[teamID], permission)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	for teamID, perms := range out {
		sort.Strings(perms)
		out[teamID] = dedupeSortedStrings(perms)
	}
	return out, nil
}

func dedupeSortedStrings(values []string) []string {
	if len(values) == 0 {
		return values
	}
	out := values[:1]
	for _, v := range values[1:] {
		if v != out[len(out)-1] {
			out = append(out, v)
		}
	}
	return out
}
