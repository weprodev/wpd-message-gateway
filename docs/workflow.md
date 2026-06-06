# Development Workflow

CI/CD pipeline, commit conventions, and release process.

## Workflow Flow

```
┌─────────────────────────────────────────────────────────────────┐
│                         PULL REQUEST                            │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│   ci.yml (runs on PR only)                                      │
│   ┌────────┬────────┬──────────┬──────────┬──────────┬────────┐ │
│   │ 🔍Lint │ 🧪Test │ 🔨Build  │ 🔒Security│ 🌐Bruno  │ 🎨Web │ │
│   └────────┴────────┴──────────┴──────────┴──────────┴────────┘ │
│                            │                                    │
│                       📊 Summary                                │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
                              │
                         [MERGE PR]
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                      PUSH TO MASTER                             │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│   release.yml                                                   │
│   ┌──────────────┐                                              │
│   │ 📦 Release   │ ──(creates tag)──▶ outputs: new_version      │
│   └──────────────┘                                              │
│          │                                                      │
│     needs: release                                              │
│          ▼                                                      │
│   ┌──────────────┐                                              │
│   │ 🐳 Docker    │ ──(builds & pushes image with version tag)   │
│   └──────────────┘                                              │
│          │                                                      │
│          ▼                                                      │
│   ┌──────────────┐                                              │
│   │ 📊 Summary   │                                              │
│   └──────────────┘                                              │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

## Commit Conventions

We use [Conventional Commits](https://www.conventionalcommits.org/) for automatic versioning.

```
<type>(<scope>): <description>
```

| Type                        | Bump  | Example                           |
| --------------------------- | ----- | --------------------------------- |
| `feat`                      | Minor | `feat(sendgrid): add provider`    |
| `fix`                       | Patch | `fix(mailgun): handle rate limit` |
| `feat!`                     | Major | `feat!: new API format`           |
| `docs`, `chore`, `refactor` | Patch | `docs: update readme`             |

## CI Pipeline

Runs on **pull requests** to `master`:

| Job         | What it does                             |
| ----------- | ---------------------------------------- |
| 🔍 Lint     | `gofmt`, `golangci-lint`                 |
| 🧪 Test     | `go test -race` + coverage               |
| 🔨 Build    | Compile binary                           |
| 🔒 Security | `govulncheck`                            |
| 🌐 API Test | Bruno CLI tests (`bru run --env memory`) |
| 🎨 Web UI   | Build Portal frontend                    |

```bash
# Run locally
make audit
```

## Release Process

Automatic on push to `master`. Manual trigger available in Actions → Release.

| Commits contain               | Version bump            |
| ----------------------------- | ----------------------- |
| `BREAKING CHANGE` or `feat!:` | Major (v1 → v2)         |
| `feat:`                       | Minor (v1.0 → v1.1)     |
| `fix:`, `docs:`, etc.         | Patch (v1.0.0 → v1.0.1) |

## Docker Image

On release, the Docker image is published to:

```
ghcr.io/weprodev/wpd-message-gateway:latest
ghcr.io/weprodev/wpd-message-gateway:v1.0.0
```

## Cleanup

Old container images are automatically cleaned up **monthly**:

- 🗑️ Delete untagged images
- 📦 Keep last 10 tagged versions
- 🏷️ Keep last 5 pre-release versions

Manual trigger: Actions → Cleanup → Run workflow

## Dependabot

Dependencies are checked **weekly** for security updates:

- ✅ Minor and patch updates only
- ❌ Major versions require manual review
- 📦 Covers: Go modules, npm, GitHub Actions

## E2E Testing

Use the gateway Docker image to capture and verify messages in your CI tests.

→ See **[E2E Testing Guide](./backend/e2e-testing.md)** for complete examples.

Quick example:

```yaml
services:
  gateway:
    image: ghcr.io/weprodev/wpd-message-gateway:latest
    ports:
      - 10101:10101

steps:
  - run: npm test
    env:
      EMAIL_API: http://localhost:10101

  - run: curl -H "Authorization: Bearer $JWT" -H "X-Api-Client-Id: $ID" -H "X-Api-Client-Secret: $SECRET" "http://localhost:10101/api/v1/workspaces/$WID/inbox/emails" | jq '.emails[0].email.subject'
```

## Branch Strategy

```
master ───●─────●─────●───── (releases)
          │     │
          │     └── feat/sendgrid
          └── fix/rate-limit
```

## PR Checklist

- [ ] **`/smell develop`** — no BLOCKER/HIGH (`.claude/commands/smell.md`, [verification.md](./agents/verification.md))
- [ ] **`make audit`** passes
- [ ] Tests added
- [ ] Commits follow conventions

## Related

- [E2E Testing](./backend/e2e-testing.md) — Test your app's messages
- [Contributing](./backend/contributing.md) — Add new providers
- [Code Conventions](./backend/code-conventions.md) — Coding standards
