# Development Workflow

This guide covers the CI/CD pipeline, commit conventions, and release process.

## Commit Conventions

We use [Conventional Commits](https://www.conventionalcommits.org/) for automatic versioning.

### Format

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

### Types

| Type | Description | Version Bump |
|------|-------------|--------------|
| `feat` | New feature | Minor |
| `fix` | Bug fix | Patch |
| `docs` | Documentation only | Patch |
| `chore` | Maintenance tasks | Patch |
| `refactor` | Code refactoring | Patch |
| `test` | Adding tests | Patch |
| `feat!` | Breaking change | Major |

### Examples

```bash
# Patch release (v1.0.0 → v1.0.1)
git commit -m "fix: resolve email validation error"
git commit -m "docs: update README"
git commit -m "chore: update dependencies"

# Minor release (v1.0.0 → v1.1.0)
git commit -m "feat: add SendGrid email provider"
git commit -m "feat(sms): add Twilio support"

# Major release (v1.0.0 → v2.0.0)
git commit -m "feat!: change API response format"
git commit -m "refactor: rename Email.Body to Email.HTML

BREAKING CHANGE: Email.Body field renamed to Email.HTML"
```

## CI Pipeline

Every push and pull request triggers the CI workflow:

```
┌──────────────────────────────────────────────────────────┐
│                    CI Pipeline                           │
├──────────────────────────────────────────────────────────┤
│                                                          │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐      │
│  │   Lint &    │  │    Unit     │  │   Build     │      │
│  │   Format    │  │   Tests     │  │   Check     │      │
│  └─────────────┘  └─────────────┘  └─────────────┘      │
│                                                          │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐      │
│  │  Security   │  │    API      │  │  Frontend   │      │
│  │    Scan     │  │   Tests     │  │   Build     │      │
│  └─────────────┘  └─────────────┘  └─────────────┘      │
│                                                          │
└──────────────────────────────────────────────────────────┘
```

### Jobs

| Job | Description | Tools |
|-----|-------------|-------|
| **Lint & Format** | Code style check | `gofmt`, `golangci-lint` |
| **Unit Tests** | Run Go tests | `go test -race` |
| **Build** | Compile binaries | `go build` |
| **Security Scan** | Vulnerability check | `govulncheck` |
| **API Tests** | Run Bruno tests | `bru run` |
| **Frontend Build** | Build React UI | `npm run build` |

### Running Locally

```bash
# Run all checks (same as CI)
make audit

# Individual checks
make lint
make test
make build
```

## Release Process

Releases are created automatically when CI passes on `main`.

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│  CI Passes  │ ──▶ │  Detect     │ ──▶ │  Create     │
│  on main    │     │  Version    │     │  Release    │
└─────────────┘     └─────────────┘     └─────────────┘
```

### Automatic Versioning

The release workflow:

1. Gets the latest Git tag (e.g., `v1.2.3`)
2. Analyzes commits since last tag
3. Determines version bump type from commit messages
4. Creates new tag and GitHub release

### Version Detection

| Commits contain | Bump | Example |
|-----------------|------|---------|
| `BREAKING CHANGE` or `feat!:` | Major | `v1.0.0` → `v2.0.0` |
| `feat:` | Minor | `v1.0.0` → `v1.1.0` |
| `fix:`, `docs:`, etc. | Patch | `v1.0.0` → `v1.0.1` |

### Manual Release

You can trigger a release manually:

1. Go to **Actions** → **Release**
2. Click **Run workflow**
3. Select version type: `patch`, `minor`, or `major`
4. Click **Run workflow**

## Branch Strategy

```
main ─────●─────●─────●─────●───── (releases)
          │     │     │
          │     │     └── fix/email-validation
          │     └── feat/sendgrid-provider
          └── docs/update-readme
```

- **main**: Production-ready code, releases are tagged here
- **feature branches**: `feat/`, `fix/`, `docs/`, `chore/`

## Pull Request Workflow

1. Create branch from `main`
2. Make changes with conventional commits
3. Push and create PR
4. CI runs automatically
5. Get review and approval
6. Merge to `main`
7. Release is created automatically

### PR Checklist

```markdown
- [ ] `make audit` passes locally
- [ ] Tests added for new features
- [ ] Documentation updated if needed
- [ ] Commit messages follow conventions
```

## Troubleshooting

### CI Failed

```bash
# Check formatting
gofmt -d .

# Check linting
golangci-lint run ./...

# Run tests
go test -v ./...
```

### Release Not Created

1. Check if CI passed on `main`
2. Verify commit messages follow conventions
3. Check Actions tab for errors

### Manual Tag Creation

If automatic release fails:

```bash
# Create tag manually
git tag -a v1.2.3 -m "Release v1.2.3"
git push origin v1.2.3
```

## Related Documentation

- [Contributing Guide](./contributing.md) - How to contribute code
- [Code Conventions](./code-conventions.md) - Coding standards
