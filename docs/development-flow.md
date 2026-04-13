# Development Flow (Spec Kit + GitHub)

This repo uses GitHub Spec Kit to drive a predictable flow from idea → spec → plan → tasks → implementation → review → PR.

Key rule: **before creating a PR**, we update the **GitHub issue** so it reflects the latest spec/plan/tasks (links + current status).


```mermaid
flowchart TD
  A[make setup<br/>one-time: tools + defaults] -->|writes| A1[.specify/memory/setup.json<br/>default_branch + model/tool prefs]

  B[make specify] --> B1[Prompt: feature description]
  B1 --> B2[Fetch + sync default branch<br/>(master/main from setup.json)]
  B2 --> B3[Create feature branch + specs/<branch>/spec.md<br/>(create-new-feature.sh)]
  B3 --> B4{gh authenticated?}
  B4 -->|yes| B5[Create GitHub Issue<br/>(title + body + links)]
  B4 -->|no| B6[Skip issue creation<br/>show: gh auth login]
  B5 --> B7[Write specs/<branch>/meta.json<br/>(issue_number, issue_url)]
  B6 --> B7
  B7 --> B8[Copy slash command to clipboard<br/>/speckit.specify <description>]

  B8 --> C[/speckit.specify<br/>writes spec content into spec.md]
  C --> D[/speckit.clarify<br/>ask ≤5 questions; updates spec.md]
  D --> E[/speckit.plan<br/>plan.md + research/data-model/contracts/quickstart<br/>+ update-agent-context]
  E --> F[/speckit.tasks<br/>tasks.md: atomic, dependency-ordered]
  F --> G[/speckit.analyze<br/>read-only consistency check]
  G --> H[/speckit.implement<br/>execute tasks; mark tasks complete]

  H --> I[make checklist DOMAIN=security|api|ux]
  I --> J[make audit]

  J --> K{Ready to open PR?}
  K -->|no| H

  K -->|yes| U[make sync<br/>Auto-pushes spec.md, plan.md, tasks.md<br/>to GitHub Issue body]
  U --> L[make pr]
  L --> L1{meta.json has issue_number?}
  L1 -->|yes| L2[Create PR with "Closes #<issue>"<br/>auto-links to issue]
  L1 -->|no| L3[Create PR without issue linkage]
```

## AI-Assisted Development

This project uses [GitHub Spec Kit](https://github.com/github/spec-kit) for specification-driven development. Each feature goes through: **specify → plan → tasks → implement → review**.

<details>
<summary><strong>📖 How it works (click to expand)</strong></summary>

### Feature Workflow

1. **`/speckit.specify <description>`** — Create spec + feature branch → `specs/<branch>/spec.md`
2. **`/speckit.clarify`** *(optional)* — Ask high-impact questions, refine the spec
3. **`/speckit.plan`** — Generate implementation plan + design artifacts
4. **`/speckit.tasks`** — Generate file-path-specific, dependency-ordered tasks
5. **`/speckit.analyze`** *(optional)* — Consistency check (gap detection)
6. **`/speckit.implement`** — Execute tasks phase-by-phase
7. **`/speckit.checklist <domain>`** + `make audit` — Pre-PR quality gates

### Agent Mapping

| Spec Kit Command | Agent |
|------------------|-------|
| `/speckit.specify`, `/speckit.plan` | [Master Agent](docs/agents/master-agent.md) |
| `/speckit.tasks`, `/speckit.implement` | [Delivery Agent](docs/agents/delivery-agent.md) |
| `/speckit.checklist`, `/speckit.analyze` | [Review Agent](docs/agents/review-agent.md) |

</details>
