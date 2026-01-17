# Development Workflow

CI/CD pipeline, commit conventions, and release process.

## Commit Conventions

We use [Conventional Commits](https://www.conventionalcommits.org/) for automatic versioning.

```
<type>(<scope>): <description>
```

| Type | Bump | Example |
|------|------|---------|
| `feat` | Minor | `feat(sendgrid): add provider` |
| `fix` | Patch | `fix(mailgun): handle rate limit` |
| `feat!` | Major | `feat!: new API format` |
| `docs`, `chore`, `refactor` | Patch | `docs: update readme` |

## CI Pipeline

Every push triggers:

| Job | What it does |
|-----|--------------|
| 🔍 Lint | `gofmt`, `golangci-lint` |
| 🧪 Test | `go test -race` + coverage |
| 🔨 Build | Compile binary |
| 🔒 Security | `govulncheck` |
| 🌐 API Test | Bruno tests |
| 🎨 Web UI | Build DevBox frontend |

```bash
# Run locally
make audit
```

## Release Process

Automatic on `main` when CI passes. Manual trigger available in Actions → Release.

| Commits contain | Version bump |
|-----------------|--------------|
| `BREAKING CHANGE` or `feat!:` | Major (v1 → v2) |
| `feat:` | Minor (v1.0 → v1.1) |
| `fix:`, `docs:`, etc. | Patch (v1.0.0 → v1.0.1) |

## Docker Image

On release, the Docker image is published to:

```
ghcr.io/weprodev/wpd-message-gateway:latest
ghcr.io/weprodev/wpd-message-gateway:v1.0.0
```

## E2E Testing

Use the gateway Docker image to capture and verify messages in your CI tests.

→ See **[E2E Testing Guide](./e2e-testing.md)** for complete examples.

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
  
  - run: curl http://localhost:10101/api/v1/emails | jq '.emails[0].email.subject'
```

## Branch Strategy

```
main ─────●─────●─────●───── (releases)
          │     │
          │     └── feat/sendgrid
          └── fix/rate-limit
```

## PR Checklist

- [ ] `make audit` passes
- [ ] Tests added
- [ ] Commits follow conventions

## Related

- [E2E Testing](./e2e-testing.md) — Test your app's messages
- [Contributing](./contributing.md) — Add new providers
- [Code Conventions](./code-conventions.md) — Coding standards
