.PHONY: it install upgrade setup specify spec clarify clr plan pln tasks tsk implement impl analyze alyz checklist chk pr sync agents agents-kill start stop test audit build clean docker-check dev dev-down db-init seed-demo help ui-install ui ui-build ui-format ui-test ui-lint storybook

# ============================================================================
# ANSI Color Codes
# ============================================================================
CYAN    := \033[36m
GREEN   := \033[32m
YELLOW  := \033[33m
BLUE    := \033[34m
MAGENTA := \033[35m
BOLD    := \033[1m
RESET   := \033[0m

# Go packages only (excludes frontend/node_modules vendored Go stubs)
GO_PKGS := ./cmd/... ./internal/... ./pkg/...

.DEFAULT_GOAL := help

# ============================================================================
# Core Commands
# ============================================================================

## Upgrade all dependencies (Go + frontend)
upgrade:
	@printf "\n"
	@printf "$(BOLD)$(CYAN)⬆️  Upgrading all dependencies...$(RESET)\n"
	@printf "$(CYAN)━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━$(RESET)\n"
	@printf "\n"
	@printf "$(BOLD)$(YELLOW)📦 Upgrading Go dependencies...$(RESET)\n"
	@go get -u ./...
	@go mod tidy
	@printf "$(GREEN)✅ Go dependencies upgraded!$(RESET)\n"
	@printf "\n"
	@printf "$(BOLD)$(YELLOW)🌐 Upgrading frontend dependencies...$(RESET)\n"
	@cd frontend && npm update
	@cd frontend && npm audit fix --force 2>/dev/null || true
	@printf "$(GREEN)✅ Frontend dependencies upgraded!$(RESET)\n"
	@printf "\n"
	@printf "$(CYAN)━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━$(RESET)\n"
	@printf "$(BOLD)$(GREEN)✅ All dependencies upgraded!$(RESET)\n"
	@printf "$(CYAN)━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━$(RESET)\n"
	@printf "\n"
	@printf "$(BOLD)$(MAGENTA)💡 Tip:$(RESET) Run $(YELLOW)make audit$(RESET) to verify everything works\n"
	@printf "\n"

## Install all dependencies (Go + tools + frontend)
install:
	@printf "\n"
	@printf "$(BOLD)$(CYAN)📦 Installing Go dependencies...$(RESET)\n"
	@go mod download
	@go mod tidy
	@printf "$(GREEN)✅ Go dependencies installed!$(RESET)\n"
	@printf "\n"
	@printf "$(BOLD)$(CYAN)🔧 Installing development tools...$(RESET)\n"
	@go install golang.org/x/tools/cmd/goimports@latest
	@go install golang.org/x/vuln/cmd/govulncheck@latest
	@go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.8.0
	@printf "$(GREEN)✅ Development tools installed!$(RESET)\n"
	@printf "\n"
	@printf "$(BOLD)$(CYAN)🌐 Installing frontend dependencies...$(RESET)\n"
	@cd frontend && npm install
	@printf "$(GREEN)✅ Frontend dependencies installed!$(RESET)\n"
	@printf "\n"
	@printf "$(BOLD)$(CYAN)━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━$(RESET)\n"
	@printf "$(BOLD)$(GREEN)🎉 Installation complete!$(RESET)\n"
	@printf "$(BOLD)$(CYAN)━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━$(RESET)\n"
	@printf "\n"
	@printf "$(BOLD)$(MAGENTA)💡 Next step:$(RESET) Run $(YELLOW)make start$(RESET) to begin development\n"
	@printf "\n"

## Initialize configuration for FlowAI
setup:
	@flowai init

# ============================================================================
# Spec Kit / AI Workflow Commands
# ============================================================================

## Spec Kit: Interactive workflow orchestrator
it:
	@flowai run master

## Spec Kit: create a new feature spec (FEATURE="...")
specify spec:
	@flowai run spec

## Spec Kit: produce plan.md (+ design artifacts)
plan pln:
	@flowai run plan

## Spec Kit: produce tasks.md
tasks tsk:
	@flowai run tasks

## Spec Kit: implement tasks.md phase-by-phase
implement impl:
	@flowai run implement

## Spec Kit: review implementation natively
review rv:
	@flowai run review

## Launch tmux multi-agent session (requires make setup first)
agents:
	@flowai start

## Kill the active tmux agent session
agents-kill:
	@flowai kill

# ============================================================================
# Local Server Operations
# ============================================================================


## Start development environment (Gateway + Portal UI)
start: stop
	@printf "\n"
	@printf "$(BOLD)$(CYAN)🚀 Starting development environment...$(RESET)\n"
	@printf "$(CYAN)━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━$(RESET)\n"
	@if [ ! -d "frontend/node_modules" ]; then \
		printf "$(YELLOW)📦 Installing frontend dependencies...$(RESET)\n"; \
		cd frontend && npm install; \
	fi
	@# Build fresh binary - force rebuild to avoid cache issues
	@printf "$(YELLOW)🔨 Building server...$(RESET)\n"
	@rm -f ./bin/server
	@go build -a -o ./bin/server ./cmd/server
	@# Validate config before starting both processes
	@printf "$(YELLOW)🔍 Validating configuration...$(RESET)\n"
	@env -u MESSAGE_DEFAULT_EMAIL_PROVIDER -u MESSAGE_DEFAULT_SMS_PROVIDER \
		-u MESSAGE_DEFAULT_PUSH_PROVIDER -u MESSAGE_DEFAULT_CHAT_PROVIDER \
		CONFIG_PATH=configs/local.yml ./bin/server & PID=$$!; sleep 1; \
		if kill -0 $$PID 2>/dev/null; then \
			kill $$PID 2>/dev/null; \
			printf "$(GREEN)✅ Configuration valid!$(RESET)\n"; \
		else \
			printf "\n$(BOLD)$(MAGENTA)💡 Tip: Copy configs/local.example.yml to configs/local.yml and configure your providers$(RESET)\n\n"; \
			exit 1; \
		fi
	@printf "\n"
	@printf "   $(BOLD)Gateway API:$(RESET)  http://localhost:10101\n"
	@printf "   $(BOLD)Portal UI:$(RESET)    http://localhost:10104\n"
	@printf "\n"
	@printf "$(YELLOW)Press Ctrl+C to stop both servers$(RESET)\n"
	@printf "$(CYAN)━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━$(RESET)\n"
	@trap 'kill %1 2>/dev/null' EXIT; \
		env -u MESSAGE_DEFAULT_EMAIL_PROVIDER -u MESSAGE_DEFAULT_SMS_PROVIDER \
		-u MESSAGE_DEFAULT_PUSH_PROVIDER -u MESSAGE_DEFAULT_CHAT_PROVIDER \
		CONFIG_PATH=configs/local.yml ./bin/server & \
		cd frontend && npm run dev

## Stop any running gateway processes
stop:
	@if lsof -i :10101 >/dev/null 2>&1; then \
		printf "$(YELLOW)🛑 Stopping existing server on port 10101...$(RESET)\n"; \
		lsof -ti :10101 | xargs kill -9 2>/dev/null || true; \
		sleep 1; \
	fi
	@if lsof -i :10104 >/dev/null 2>&1; then \
		printf "$(YELLOW)🛑 Stopping existing frontend on port 10104...$(RESET)\n"; \
		lsof -ti :10104 | xargs kill -9 2>/dev/null || true; \
	fi

## Run tests
test:
	@printf "\n"
	@printf "$(BOLD)$(CYAN)🧪 Running tests...$(RESET)\n"
	@go test $(GO_PKGS)
	@printf "$(GREEN)✅ All tests passed!$(RESET)\n"

## Full quality check: format + lint + test (Go + frontend), govulncheck, then compile checks (Go + Vite + Storybook)
audit:
	@printf "\n"
	@printf "$(BOLD)$(CYAN)🔍 Running full audit (Go + frontend)...$(RESET)\n"
	@printf "$(CYAN)━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━$(RESET)\n"
	@printf "\n"
	@printf "$(BOLD)$(YELLOW)🎨 Formatting Go...$(RESET)\n"
	@find . -name '*.go' -not -path './.history/*' -not -path './vendor/*' | xargs goimports -local github.com/weprodev/wpd-message-gateway -w 2>/dev/null || \
		find . -name '*.go' -not -path './.history/*' -not -path './vendor/*' | xargs gofmt -w
	@go mod tidy
	@printf "$(GREEN)✅ Go formatted!$(RESET)\n"
	@printf "\n"
	@printf "$(BOLD)$(YELLOW)🎨 Formatting frontend (ESLint --fix)...$(RESET)\n"
	@cd frontend && npm run format
	@printf "$(GREEN)✅ Frontend formatted!$(RESET)\n"
	@printf "\n"
	@printf "$(BOLD)$(YELLOW)🔍 Linting Go...$(RESET)\n"
	@golangci-lint run $(GO_PKGS)
	@printf "$(GREEN)✅ Go lint passed!$(RESET)\n"
	@printf "\n"
	@printf "$(BOLD)$(YELLOW)🔍 Linting frontend...$(RESET)\n"
	@cd frontend && npm run lint
	@printf "$(GREEN)✅ Frontend lint passed!$(RESET)\n"
	@printf "\n"
	@printf "$(BOLD)$(YELLOW)🧪 Testing Go...$(RESET)\n"
	@go test $(GO_PKGS)
	@printf "$(GREEN)✅ Go tests passed!$(RESET)\n"
	@printf "\n"
	@printf "$(BOLD)$(YELLOW)🧪 Testing frontend...$(RESET)\n"
	@cd frontend && npm run test
	@printf "$(GREEN)✅ Frontend tests passed!$(RESET)\n"
	@printf "\n"
	@printf "$(BOLD)$(YELLOW)🔒 govulncheck (Go)...$(RESET)\n"
	@govulncheck $(GO_PKGS)
	@printf "$(GREEN)✅ No Go vulnerabilities reported!$(RESET)\n"
	@printf "\n"
	@printf "$(BOLD)$(YELLOW)🔨 Verifying builds (Go + Vite app + Storybook)...$(RESET)\n"
	@go build $(GO_PKGS)
	@cd frontend && npm run build:all
	@printf "$(GREEN)✅ All builds succeeded!$(RESET)\n"
	@printf "\n"
	@printf "$(CYAN)━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━$(RESET)\n"
	@printf "$(BOLD)$(GREEN)✅ Audit complete.$(RESET)\n"
	@printf "$(CYAN)━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━$(RESET)\n"

## Build: Go packages + frontend production bundle + static Storybook
build:
	@printf "\n"
	@printf "$(BOLD)$(CYAN)🔨 Building Go + frontend (app + Storybook)...$(RESET)\n"
	@go build $(GO_PKGS)
	@cd frontend && npm run build:all
	@printf "$(GREEN)✅ Build successful!$(RESET)\n"

## Clean build artifacts
clean:
	@printf "\n"
	@printf "$(BOLD)$(CYAN)🧹 Cleaning build artifacts...$(RESET)\n"
	@rm -f coverage.out coverage.html
	@rm -rf ./bin
	@go clean -cache -testcache
	@printf "$(GREEN)✅ Cleaned!$(RESET)\n"

# ============================================================================
# UI (Frontend) Commands
# ============================================================================

## Install front-end dependencies
ui-install:
	@printf "\n"
	@printf "$(BOLD)$(CYAN)📦 Installing UI dependencies...$(RESET)\n"
	@cd frontend && npm install
	@printf "$(GREEN)✅ Installed UI dependencies!$(RESET)\n"

## Start UI in development mode
ui:
	@printf "\n"
	@printf "$(BOLD)$(CYAN)🚀 Starting local UI server...$(RESET)\n"
	@cd frontend && npm run dev

## Build UI: Vite app + static Storybook (same as npm run build:all)
ui-build:
	@printf "\n"
	@printf "$(BOLD)$(CYAN)🔨 Building UI (app + Storybook)...$(RESET)\n"
	@cd frontend && npm run build:all
	@printf "$(GREEN)✅ UI build successful!$(RESET)\n"

## Format frontend (ESLint --fix)
ui-format:
	@printf "\n"
	@printf "$(BOLD)$(CYAN)🎨 Formatting frontend...$(RESET)\n"
	@cd frontend && npm run format
	@printf "$(GREEN)✅ Frontend formatted!$(RESET)\n"

## Run UI tests
ui-test:
	@printf "\n"
	@printf "$(BOLD)$(CYAN)🧪 Running UI tests...$(RESET)\n"
	@cd frontend && npm run test
	@printf "$(GREEN)✅ UI Tests passed!$(RESET)\n"

## Run UI linter
ui-lint:
	@printf "\n"
	@printf "$(BOLD)$(CYAN)🔍 Running UI linter...$(RESET)\n"
	@cd frontend && npm run lint
	@printf "$(GREEN)✅ UI Linter passed!$(RESET)\n"

## Start Storybook server
storybook:
	@printf "\n"
	@printf "$(BOLD)$(CYAN)📚 Starting Storybook server...$(RESET)\n"
	@printf "$(YELLOW)Storybook UI:$(RESET) http://localhost:6006\n"
	@cd frontend && npm run storybook

# ============================================================================
# Docker
# ============================================================================

## Verify Docker daemon is reachable (used by dev, dev-down targets)
docker-check:
	@docker info >/dev/null 2>&1 || ( \
		printf "\n$(BOLD)$(YELLOW)Docker is not running.$(RESET)\n"; \
		printf "Start Docker Desktop (macOS: open Docker from Applications), wait until it finishes starting, then retry.\n"; \
		printf "To run without Docker: $(YELLOW)make start$(RESET) (Gateway + Portal UI on localhost).\n\n"; \
		exit 1; \
	)

## Start Gateway, DB, and Portal UI via Docker Compose (with hot-reloading)
dev: docker-check
	@printf "\n"
	@printf "$(BOLD)$(CYAN)🐳 Starting Gateway, database, and Portal UI via Docker...$(RESET)\n"
	@docker compose up -d --build --renew-anon-volumes
	@printf "$(GREEN)✅ All services started!$(RESET)\n"
	@printf "\n"
	@printf "   ╭───────────────────────────────────────────────────╮\n"
	@printf "   │                                                   │\n"
	@printf "   │  $(BOLD)Gateway API:$(RESET)  http://localhost:10101             │\n"
	@printf "   │  $(BOLD)Portal UI:$(RESET)    http://localhost:10104             │\n"
	@printf "   │                                                   │\n"
	@printf "   │  $(BOLD)Demo Account$(RESET)                                     │\n"
	@printf "   │  Email:    demo@weprodev.com                      │\n"
	@printf "   │  Password: secret                                 │\n"
	@printf "   │                                                   │\n"
	@printf "   ╰───────────────────────────────────────────────────╯\n"
	@printf "\n"

## Stop Docker Compose
dev-down: docker-check
	@printf "\n"
	@printf "$(BOLD)$(CYAN)🛑 Stopping Docker...$(RESET)\n"
	@docker compose down
	@printf "$(GREEN)✅ Stopped!$(RESET)\n"

## Apply schema when Postgres volume exists but tables are missing (then restart gateway)
db-init: docker-check
	@printf "\n"
	@printf "$(BOLD)$(CYAN)🗄️  Initializing database schema...$(RESET)\n"
	@docker compose exec -T -e POSTGRES_USER=gateway -e POSTGRES_DB=gateway db \
		bash /docker-entrypoint-initdb.d/init-db.sh
	@docker compose restart gateway
	@printf "$(GREEN)✅ Database ready — gateway will apply seeds on startup (demo@weprodev.com / secret)$(RESET)\n"
	@printf "\n"

## Re-apply SQL seeds by restarting the gateway (seeds run on every startup)
seed-demo: docker-check
	@printf "\n"
	@printf "$(BOLD)$(CYAN)🌱 Re-applying database seeds (gateway restart)...$(RESET)\n"
	@docker compose restart gateway
	@printf "$(GREEN)✅ Seeds applied on gateway startup (demo@weprodev.com / secret)$(RESET)\n"
	@printf "\n"


# ============================================================================
# Help
# ============================================================================

help:
	@printf "\n"
	@printf "$(BOLD)$(CYAN)━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━$(RESET)\n"
	@printf "$(BOLD)$(CYAN)           WPD Message Gateway                              $(RESET)\n"
	@printf "$(BOLD)$(CYAN)━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━$(RESET)\n"
	@printf "\n"
	@printf "$(BOLD)$(GREEN)🚀 Core Development$(RESET)\n"
	@printf "   $(YELLOW)make install$(RESET)      Install all dependencies\n"
	@printf "   $(YELLOW)make upgrade$(RESET)      Upgrade all dependencies\n"
	@printf "   $(YELLOW)make setup$(RESET)        Interactive setup (macOS: tools + models)\n"
	@printf "   $(YELLOW)make start$(RESET)        Start Gateway + Portal UI\n"
	@printf "   $(YELLOW)make stop$(RESET)         Stop running servers\n"
	@printf "   $(YELLOW)make test$(RESET)         Run Go tests only\n"
	@printf "   $(YELLOW)make audit$(RESET)        Full check: fmt+lint+test (Go+UI), govulncheck, builds\n"
	@printf "\n"
	@printf "$(BOLD)$(GREEN)🤖 Spec Kit / AI Workflow$(RESET)\n"
	@printf "   $(YELLOW)make it$(RESET)           Start interactive end-to-end AI orchestrator\n"
	@printf "   $(YELLOW)make specify|spec$(RESET)   Create spec (FEATURE=\"...\")\n"
	@printf "   $(YELLOW)make clarify|clr$(RESET)    Ask questions to clarify spec\n"
	@printf "   $(YELLOW)make plan|pln$(RESET)       Generate plan.md (+ design artifacts)\n"
	@printf "   $(YELLOW)make tasks|tsk$(RESET)      Generate implementation tasks.md\n"
	@printf "   $(YELLOW)make analyze|alyz$(RESET)   Read-only consistency analysis\n"
	@printf "   $(YELLOW)make implement|impl$(RESET) Implement tasks.md phase-by-phase\n"
	@printf "   $(YELLOW)make checklist|chk$(RESET)  Requirements-quality check (DOMAIN=\"...\")\n"
	@printf "\n"
	@printf "$(BOLD)$(GREEN)🎨 UI / Frontend Development$(RESET)\n"
	@printf "   $(YELLOW)make ui$(RESET)           Start UI development server\n"
	@printf "   $(YELLOW)make storybook$(RESET)    Start Storybook server (port 6006)\n"
	@printf "   $(YELLOW)make ui-format$(RESET)    ESLint --fix in frontend\n"
	@printf "   $(YELLOW)make ui-test$(RESET)      Run frontend tests\n"
	@printf "   $(YELLOW)make ui-lint$(RESET)      Run frontend linter\n"
	@printf "   $(YELLOW)make ui-build$(RESET)     Build Vite app + static Storybook\n"
	@printf "\n"
	@printf "$(BOLD)$(GREEN)🔨 Build$(RESET)\n"
	@printf "   $(YELLOW)make build$(RESET)        Go compile + frontend build:all (no tests)\n"
	@printf "   $(YELLOW)make clean$(RESET)        Clean build artifacts\n"
	@printf "\n"
	@printf "$(BOLD)$(GREEN)🐳 Docker$(RESET)\n"
	@printf "   $(YELLOW)make dev$(RESET)          Start Gateway, DB, and UI via Docker (with hot-reloading)\n"
	@printf "   $(YELLOW)make dev-down$(RESET)     Stop Docker containers\n"
	@printf "   $(YELLOW)make db-init$(RESET)      Apply schema if DB volume is empty; gateway seeds on start\n"
	@printf "   $(YELLOW)make seed-demo$(RESET)    Re-apply seeds via gateway restart\n"
	@printf "\n"

	@printf "$(BOLD)$(CYAN)━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━$(RESET)\n"
	@printf "$(BOLD)$(MAGENTA)💡 Quick Start:$(RESET) make install && make start\n"
	@printf "$(BOLD)$(CYAN)━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━$(RESET)\n"
	@printf "\n"
