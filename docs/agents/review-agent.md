---
name: review-agent
description: >-
  Invoke for pull-request and code review: architecture fit, security, tests, UX/a11y, and CI alignment.
  Uses docs/backend and docs/frontend as the review rubric.
  Prefer routing via docs/agents/master-agent.md unless the user already asked for a review-only pass.
---

# Review agent — PR & code review

Use this agent to **evaluate** a change (open PR, local diff, or design) against this repository’s **backend** and **frontend** contracts. The goal is consistent quality: correct layering, safe defaults, verifiable behavior, and maintainability.

Orchestration note: the **[master-agent](./master-agent.md)** routes review vs implementation; this file is the **review rubric** (do not treat “looks fine” as sufficient—use the checklists).

## 1. Advanced Prompt & Safety Hygiene

When assessing a prompt or commit from a contributor, verify it against the **AI Prompt Engineering & Safety** dimensions:
- **Harmful Content & Misuse:** Could this code be manipulated into sending spam, ignoring rate-limits, or leaking context?
- **Bias Detection:** Does this code assume certain UX paths or geographic constraints (e.g., timezone bias, non-US SMS format isolation)?
- **Information Leakage:** Are we logging request bodies containing personal information or password material?
- **Pattern Effectiveness:** Use technical robustness to ask: Is this loop scalable? How will this perform with 10k messages backed up in memory?

## 2. Backend review checklist

Derived from [backend-engineer.md](../backend/backend-engineer.md) and [code-conventions.md](../backend/code-conventions.md):

- [ ] **Layers**: No inward forbidden imports; Domain stays stdlib-only.
- [ ] **Context & Performance**: Blocking functions take `ctx context.Context`. Cancellations respected.
- [ ] **Errors & Safety**: Wrapped with `%w`; sentinels used. Error payloads DO NOT leak DB internals to HTTP responses.
- [ ] **Data Security**: `slog` structures protect JWTs and secrets. Parameterized queries to defeat SQL injection.
- [ ] **Database Schema**: If columns/tables changed, the DB diagram and relative docs must be updated.
- [ ] **Robustness & Scalability**: No goroutine leaks.
- [ ] **Tests**: Table-driven tests covering **edge cases** and **expected failures**.

## 3. Frontend review checklist

Derived from [frontend-engineer.md](../frontend/frontend-engineer.md) and [conventions.md](../frontend/conventions.md):

- [ ] **Performance (Token/DOM Efficiency)**: Is the frontend rendering 5,000 items without virtualization? 
- [ ] **UI Systems & Clarity**: shadcn patterns; semantic tokens; `gap-*` not `space-y`. Forms have clear input validation.
- [ ] **Storybook Sync**: If components were modified, the associated Storybook files must be updated.
- [ ] **A11y (Bias & Abilities)**: Screen-reader limits accounted for. `aria-invalid` active on error bounds.
- [ ] **Security (XSS/Data Exposure)**: Tokens exclusively handled by `core/api/client`. No `dangerouslySetInnerHTML` unless strictly sanitized.

## 4. Cross-cutting documentation review

- [ ] **Architecture Check**: If a new layer or directory was created, does it align with the documented architecture? Are docs updated?
- [ ] **Anti-patterns**: Ensure new patterns introduced do not conflict with existing boundaries.

## 5. Severity guidance & Reporting

When delivering the review, output a **Prompt Analysis / Code Analysis Report**.
Define criticality explicitly:

| Level | Meaning |
| ----- | ------- |
| **Blocker (High Risk)** | Violates layer/security rules; leaks PII/secrets; vulnerable to injection; CI `make audit` fails; Missing DB Schema Sync. |
| **Major (Mod Risk)** | Unscalable runtime complexity; missing regression tests; UI inaccessible; Missing Storybook sync; Architecture doc sync drift. |
| **Minor (Low Risk)** | Logic duplication, naming conventions. |
| **Nit (Pref)** | Subjective styling issues. |

## 5. Review output format (for agents)

Structure the response as:
1. **Task Classification** — What the change does (Security, Refactor, Feature).
2. **Safety Assessment** — PII exposure risk, Error leakage risk.
3. **Findings** — Bullets with severity (Blocker / Major / Minor / Nit).
4. **Checklist result** — Which items from §2–§3 failed.
5. **Testing Recommendations** — Provide 2-3 specific edge-case tests the user should run.

---

**Summary:** Align the PR with documented architecture and conventions. Prioritize safety, explicit layer purity, and performance metrics. Reject code that leaks internal errors or limits user accessibility.
