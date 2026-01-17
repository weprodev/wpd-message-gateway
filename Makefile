.PHONY: install start stop test audit build clean dev dev-down mailpit mailpit-down help

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

.DEFAULT_GOAL := help

# ============================================================================
# Core Commands
# ============================================================================

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
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@printf "$(GREEN)✅ Development tools installed!$(RESET)\n"
	@printf "\n"
	@printf "$(BOLD)$(CYAN)🌐 Installing frontend dependencies...$(RESET)\n"
	@cd web && npm install
	@printf "$(GREEN)✅ Frontend dependencies installed!$(RESET)\n"
	@printf "\n"
	@printf "$(BOLD)$(CYAN)━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━$(RESET)\n"
	@printf "$(BOLD)$(GREEN)🎉 Installation complete!$(RESET)\n"
	@printf "$(BOLD)$(CYAN)━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━$(RESET)\n"
	@printf "\n"
	@printf "$(BOLD)$(MAGENTA)💡 Next step:$(RESET) Run $(YELLOW)make start$(RESET) to begin development\n"
	@printf "\n"

## Start development environment (Gateway + DevBox UI)
start: stop
	@printf "\n"
	@printf "$(BOLD)$(CYAN)🚀 Starting development environment...$(RESET)\n"
	@printf "$(CYAN)━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━$(RESET)\n"
	@if [ ! -d "web/node_modules" ]; then \
		printf "$(YELLOW)📦 Installing frontend dependencies...$(RESET)\n"; \
		cd web && npm install; \
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
	@printf "   $(BOLD)DevBox UI:$(RESET)    http://localhost:10104\n"
	@printf "\n"
	@printf "$(YELLOW)Press Ctrl+C to stop both servers$(RESET)\n"
	@printf "$(CYAN)━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━$(RESET)\n"
	@trap 'kill %1 2>/dev/null' EXIT; \
		env -u MESSAGE_DEFAULT_EMAIL_PROVIDER -u MESSAGE_DEFAULT_SMS_PROVIDER \
		-u MESSAGE_DEFAULT_PUSH_PROVIDER -u MESSAGE_DEFAULT_CHAT_PROVIDER \
		CONFIG_PATH=configs/local.yml ./bin/server & \
		cd web && npm run dev

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
	@go test ./...
	@printf "$(GREEN)✅ All tests passed!$(RESET)\n"

## Full quality check: format, lint, test, security scan
audit:
	@printf "\n"
	@printf "$(BOLD)$(CYAN)🔍 Running full audit...$(RESET)\n"
	@printf "$(CYAN)━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━$(RESET)\n"
	@printf "\n"
	@printf "$(BOLD)$(YELLOW)🎨 Formatting code...$(RESET)\n"
	@goimports -local github.com/weprodev/wpd-message-gateway -w . 2>/dev/null || gofmt -w .
	@go mod tidy
	@printf "$(GREEN)✅ Code formatted!$(RESET)\n"
	@printf "\n"
	@printf "$(BOLD)$(YELLOW)🔍 Running linter...$(RESET)\n"
	@golangci-lint run ./...
	@printf "$(GREEN)✅ Linting passed!$(RESET)\n"
	@printf "\n"
	@printf "$(BOLD)$(YELLOW)🧪 Running tests...$(RESET)\n"
	@go test ./...
	@printf "$(GREEN)✅ All tests passed!$(RESET)\n"
	@printf "\n"
	@printf "$(BOLD)$(YELLOW)🔒 Running security scan...$(RESET)\n"
	@govulncheck ./...
	@printf "$(GREEN)✅ No vulnerabilities found!$(RESET)\n"
	@printf "\n"
	@printf "$(CYAN)━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━$(RESET)\n"
	@printf "$(BOLD)$(GREEN)✅ All audit checks passed!$(RESET)\n"
	@printf "$(CYAN)━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━$(RESET)\n"

## Build all packages
build:
	@printf "\n"
	@printf "$(BOLD)$(CYAN)🔨 Building all packages...$(RESET)\n"
	@go build ./...
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
# Docker
# ============================================================================

## Start Gateway via Docker Compose
dev:
	@printf "\n"
	@printf "$(BOLD)$(CYAN)🐳 Starting Gateway via Docker...$(RESET)\n"
	@docker compose up -d
	@printf "$(GREEN)✅ Gateway started!$(RESET)\n"
	@printf "\n"
	@printf "   $(BOLD)Gateway API:$(RESET)  http://localhost:10101\n"
	@printf "   $(BOLD)DevBox UI:$(RESET)    Run $(YELLOW)make start$(RESET) → http://localhost:10104\n"
	@printf "\n"

## Stop Docker Compose
dev-down:
	@printf "\n"
	@printf "$(BOLD)$(CYAN)🛑 Stopping Docker...$(RESET)\n"
	@docker compose down
	@printf "$(GREEN)✅ Stopped!$(RESET)\n"

# ============================================================================
# Mailpit (Optional - for SMTP provider testing)
# ============================================================================

## Start Mailpit (SMTP testing server)
mailpit:
	@printf "\n"
	@printf "$(BOLD)$(CYAN)📬 Starting Mailpit...$(RESET)\n"
	@docker compose -f docker-compose.mailpit.yml up -d
	@printf "$(GREEN)✅ Mailpit started!$(RESET)\n"
	@printf "\n"
	@printf "   $(BOLD)SMTP Server:$(RESET)  localhost:10102\n"
	@printf "   $(BOLD)Web UI:$(RESET)       http://localhost:10103\n"
	@printf "\n"
	@printf "$(YELLOW)To forward emails to Mailpit, set in configs/local.yml:$(RESET)\n"
	@printf "   mailpit:\n"
	@printf "     enabled: true\n"
	@printf "\n"

## Stop Mailpit
mailpit-down:
	@printf "\n"
	@printf "$(BOLD)$(CYAN)🛑 Stopping Mailpit...$(RESET)\n"
	@docker compose -f docker-compose.mailpit.yml down
	@printf "$(GREEN)✅ Mailpit stopped!$(RESET)\n"

# ============================================================================
# Help
# ============================================================================

help:
	@printf "\n"
	@printf "$(BOLD)$(CYAN)━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━$(RESET)\n"
	@printf "$(BOLD)$(CYAN)           WPD Message Gateway                              $(RESET)\n"
	@printf "$(BOLD)$(CYAN)━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━$(RESET)\n"
	@printf "\n"
	@printf "$(BOLD)$(GREEN)🚀 Development$(RESET)\n"
	@printf "   $(YELLOW)make install$(RESET)      Install all dependencies\n"
	@printf "   $(YELLOW)make start$(RESET)        Start Gateway + DevBox UI\n"
	@printf "   $(YELLOW)make stop$(RESET)         Stop running servers\n"
	@printf "   $(YELLOW)make test$(RESET)         Run tests\n"
	@printf "   $(YELLOW)make audit$(RESET)        Full check (fmt + lint + test + security)\n"
	@printf "\n"
	@printf "$(BOLD)$(GREEN)🔨 Build$(RESET)\n"
	@printf "   $(YELLOW)make build$(RESET)        Build all packages\n"
	@printf "   $(YELLOW)make clean$(RESET)        Clean build artifacts\n"
	@printf "\n"
	@printf "$(BOLD)$(GREEN)🐳 Docker$(RESET)\n"
	@printf "   $(YELLOW)make dev$(RESET)          Start Gateway via Docker\n"
	@printf "   $(YELLOW)make dev-down$(RESET)     Stop Docker\n"
	@printf "\n"
	@printf "$(BOLD)$(GREEN)📬 Optional$(RESET)\n"
	@printf "   $(YELLOW)make mailpit$(RESET)      Start Mailpit (SMTP testing)\n"
	@printf "   $(YELLOW)make mailpit-down$(RESET) Stop Mailpit\n"
	@printf "\n"
	@printf "$(BOLD)$(CYAN)━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━$(RESET)\n"
	@printf "$(BOLD)$(MAGENTA)💡 Quick Start:$(RESET) make install && make start\n"
	@printf "$(BOLD)$(CYAN)━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━$(RESET)\n"
	@printf "\n"
