# FlowAI Knowledge Graph Report

> Built: 2026-04-12T15:26:46Z
> Nodes: 30 · Edges: 103 · Communities: 23
> Specs: 2 · Spec→Code edges: 1

---

## God Nodes

The highest-degree nodes in the graph — the architectural hubs everything
depends on. Start here when you don't know where to look.

- **README.md** (docs/README.md) — degree 15
- **README.md** (docs/frontend/README.md) — degree 10
- **SKILL.md** (docs/frontend/shadcn/SKILL.md) — degree 10
- **delivery-agent.md** (docs/agents/delivery-agent.md) — degree 9
- **master-agent.md** (docs/agents/master-agent.md) — degree 7
- **frontend-engineer.md** (docs/frontend/frontend-engineer.md) — degree 6
- **review-agent.md** (docs/agents/review-agent.md) — degree 5
- **architecture.md** (docs/backend/architecture.md) — degree 5
- **backend-engineer.md** (docs/backend/backend-engineer.md) — degree 5
- **contributing.md** (docs/backend/contributing.md) — degree 5

---

## Spec Coverage (Spec-Driven Development)

FlowAI uses Spec-Driven Development. Spec nodes carry higher trust than
regular source files — they are the **authoritative source of intent**.

**2 spec documents · 1 spec→code edges · 0 code→spec (IMPLEMENTS) edges**

### Spec Status Dashboard

| Status | Count | Meaning |
|---|---|---|
| 🟢 implemented | 0 | **Frontmatter** declares implemented (workflow / intent) |
| 🔵 in-progress  | 0 | Actively being implemented |
| ⬜ planned      | 0 | Accepted, not yet started |
| 🔴 deprecated   | 0 | No longer relevant |

> **SDD vs git:** Status here comes from spec YAML. **Git evidence** (commits that reference a spec ID) is compiled separately — see **Project evolution** below after `flowai graph chronicle`.

### Spec Inventory

- **Implementation Plan: Improve Readme File Content** (`specs/018-improve-readme-file/plan.md`) [unknown] — IDs: none · 0 criteria
- **Specification: Refactor Authentication System** (`specs/018-improve-readme-file/spec.md`) [unknown] — IDs: none · 0 criteria

> [!TIP]
> Run `flowai graph chronicle` to mine git history for implementation evidence.
> Run `flowai graph lint` to detect gaps: unimplemented specs, unspecified code, zombie specs.

### Project evolution (compiled history)

Aligned with **persistent wiki** ideas: the graph stores a **denormalized timeline** on each spec node (`evolution[]`) and **IMPLEMENTS** edges from commits — so agents read `graph.json` / this report instead of re-walking `git log`.

| Metric | Value |
|---|---|
| Evolution events (total) | 0 |
| Specs with ≥1 git-linked event | 0 |
| Code→spec IMPLEMENTS edges | 0 |

When these counts are non-zero, chronicle has linked **repository activity** to **spec IDs** (works for `src/`, `lib/`, `internal/`, `apps/*`, packages, etc. — any touched path except spec-only dirs). Reference IDs in commits, e.g. `Implements UC-AUTH-001`.

---

## Architectural Insights

Key design decisions and patterns extracted from source and documentation.
Treat **INFERRED** insights as hypotheses until verified.

_No insights extracted yet. Run `flowai graph ingest <spec-file>` to populate._

---

## Ambiguous Relationships

These relationships were flagged during extraction as uncertain.
They may reflect undocumented dependencies or extraction errors.

_No ambiguous edges — clean graph._

---

## Suggested Queries

These questions can be answered efficiently using this knowledge graph:

1. What are the core phases in the FlowAI pipeline and how do they signal each other?
2. Which modules does the AI tool dispatcher (`ai.sh`) depend on?
3. What does the skill resolution chain look like end-to-end?
4. Which specs have no corresponding implementation coverage?
5. Are there implementation divergences from existing spec documents?

---

## Community Structure

### Centrality Classes

| Class | Description |
|---|---|
| god   | Central hubs with ≥10 edges — architectural load-bearers |
| hub   | Well-connected modules with 5-9 edges |
| leaf  | Peripheral files with <5 edges |

### Detected Communities (Label Propagation)

Nodes are grouped by `community_id` — clusters of related modules identified via label propagation. Agents can use `community_id` in `graph.json` to find related code.

- **.architecture** (3 members)
- **.....frontend.** (2 members)
- **..backend.architecture** (2 members)
- **..workflow** (2 members)
- **.SKILL.md#updating-components** (2 members)
- **specs.018-improve-readme-file.meta** (2 members)
- **#critical-rules** (1 members)
- **.....frontend.README** (1 members)
- **.....internal.app.imports.go** (1 members)
- **..backend.backend-engineer** (1 members)

---

## Navigation Protocol

```
GRAPH_REPORT.md  →  index.md  →  wiki/<topic>.md  →  graph.json  →  source files
```

**For this project's SDD workflow:**
```
Specs (.specify/, specs/)  ─ SPECIFIES edges ─►  Implementation (src/)
                            ``flowai graph lint`` detects divergence
```

1. Read this file first for architectural orientation
2. Consult spec nodes for authoritative intent before touching source
3. Use `index.md` to find specific wiki pages by concept
4. Use `graph.json` for multi-hop dependency and spec-traceability queries
5. Only open raw source files for implementation details

---

_Generated by FlowAI graph engine. Run `flowai graph update` for an incremental refresh._
_Run `flowai graph chronicle` to mine git history for implementation evidence._
