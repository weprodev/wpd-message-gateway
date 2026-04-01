# Implementation Plan: [FEATURE]

**Branch**: `[###-feature-name]` | **Date**: [DATE] | **Spec**: [link]
**Input**: Feature specification from `/specs/[###-feature-name]/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/plan-template.md` for the execution workflow.

## Summary

[Extract from feature spec: primary requirement + technical approach from research]

## Technical Context

This repository is a **Go HTTP API + Go SDK + React Portal UI** system.

Fill these fields with feature-specific details (leave `NEEDS CLARIFICATION` only when the missing detail would change architecture or acceptance tests).

**Backend**:
- **Language/Version**: Go (per `go.mod`)
- **Architecture**: Hexagonal / DDD (Domain + Ports + Services; Presentation/Infrastructure outside)
- **Key packages**: `internal/core/*`, `internal/presentation/*`, `internal/infrastructure/*`, `pkg/*`

**Frontend**:
- **Framework**: React + Vite + TypeScript + Tailwind (Portal UI in `frontend/`)
- **UI System**: `shadcn/ui` composition; Storybook required for UI component changes

**Storage**:
- **Database**: PostgreSQL (migrations in `database/migrations/`)
- **Credentials**: provider secrets stored encrypted in DB (Portal-managed), not in config files

**Testing/Quality Gate**:
- **Command**: `make audit` (required before PR)

## Scope & Non-Goals

- **In scope**: [explicitly list what this feature changes]
- **Out of scope**: [explicit exclusions to prevent scope creep]

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

[Gates determined based on constitution file]

## Project Structure

### Documentation (this feature)

```text
specs/[###-feature]/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
cmd/
internal/
├── app/
├── core/
│   ├── domain/
│   ├── port/
│   └── service/
├── infrastructure/
└── presentation/
pkg/
database/
frontend/
```

**Files touched by this feature (expected)**:
- Backend: [list expected packages/files]
- Frontend: [list expected pages/components/stories]
- Database: [list migrations/seed changes, if any]

## Verification Plan (Mandatory)

At minimum:
- `make audit`

If feature adds/changing endpoints:
- Add/adjust Bruno requests under `tests/bruno/` (or document why not)

If feature changes UI components:
- Update Storybook stories under `frontend/src/components/**/*.stories.tsx`

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| [e.g., 4th project] | [current need] | [why 3 projects insufficient] |
| [e.g., Repository pattern] | [specific problem] | [why direct DB access insufficient] |
