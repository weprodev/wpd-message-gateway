.PHONY: all clean help setup setup-tools setup-go test lint vulncheck audit build fmt system-info

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

# Detect OS
UNAME_S := $(shell uname -s)
GOBIN := $(shell go env GOPATH)/bin

# OS-specific settings
ifeq ($(UNAME_S),Linux)
	OS := linux
else ifeq ($(UNAME_S),Darwin)
	OS := macos
else
	OS := unknown
endif

# Default target: show help
.DEFAULT_GOAL := help

# ============================================================================
# Setup commands
# ============================================================================
setup: setup-tools setup-go
	@printf "\n"
	@printf "$(BOLD)$(GREEN)🎉 Setup completed successfully!$(RESET)\n"
	@printf "$(BOLD)$(CYAN)💡 Run 'make test' to run tests$(RESET)\n"

setup-tools:
	@printf "\n"
	@printf "$(BOLD)$(CYAN)🔧 Installing development tools for $(OS)...$(RESET)\n"
	@go install golang.org/x/tools/cmd/goimports@latest
	@go install golang.org/x/vuln/cmd/govulncheck@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@printf "$(GREEN)✅ Development tools installed successfully!$(RESET)\n"

setup-go:
	@printf "\n"
	@printf "$(BOLD)$(CYAN)🔧 Setting up Go environment...$(RESET)\n"
	@go mod tidy
	@printf "$(GREEN)✅ Go environment setup completed!$(RESET)\n"

# ============================================================================
# Build commands
# ============================================================================
build:
	@printf "\n"
	@printf "$(BOLD)$(CYAN)🔨 Building all packages...$(RESET)\n"
	@go build ./...
	@printf "$(GREEN)✅ All packages built successfully!$(RESET)\n"

fmt:
	@printf "\n"
	@printf "$(BOLD)$(CYAN)🎨 Formatting code...$(RESET)\n"
	@goimports -local github.com/weprodev/go-message-gateway -w . 2>/dev/null || gofmt -w .
	@go mod tidy
	@printf "$(GREEN)✅ Code formatted!$(RESET)\n"

# ============================================================================
# Testing and security
# ============================================================================
test:
	@printf "\n"
	@printf "$(BOLD)$(CYAN)🧪 Running tests...$(RESET)\n"
	@go test -v ./...
	@printf "$(GREEN)✅ All tests passed!$(RESET)\n"

test-short:
	@printf "\n"
	@printf "$(BOLD)$(CYAN)🧪 Running tests (short mode)...$(RESET)\n"
	@go test -short ./...
	@printf "$(GREEN)✅ All tests passed!$(RESET)\n"

test-cover:
	@printf "\n"
	@printf "$(BOLD)$(CYAN)🧪 Running tests with coverage...$(RESET)\n"
	@go test -cover ./...
	@printf "$(GREEN)✅ All tests passed!$(RESET)\n"

test-cover-html:
	@printf "\n"
	@printf "$(BOLD)$(CYAN)🧪 Generating coverage report...$(RESET)\n"
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@printf "$(GREEN)✅ Coverage report generated: coverage.html$(RESET)\n"

lint:
	@printf "\n"
	@printf "$(BOLD)$(CYAN)🔍 Running linter...$(RESET)\n"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		printf "$(YELLOW)⚠️  golangci-lint not found, installing...$(RESET)\n"; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
		$(MAKE) lint; \
	fi
	@printf "$(GREEN)✅ Linting passed!$(RESET)\n"

vulncheck:
	@printf "\n"
	@printf "$(BOLD)$(CYAN)🔒 Running vulnerability check...$(RESET)\n"
	@if command -v govulncheck >/dev/null 2>&1; then \
		govulncheck ./...; \
	else \
		printf "$(YELLOW)⚠️  govulncheck not found, installing...$(RESET)\n"; \
		go install golang.org/x/vuln/cmd/govulncheck@latest; \
		$(MAKE) vulncheck; \
	fi
	@printf "$(GREEN)✅ No vulnerabilities found!$(RESET)\n"

audit: fmt lint test vulncheck
	@printf "\n"
	@printf "$(BOLD)$(GREEN)✅ All audit checks passed!$(RESET)\n"

# ============================================================================
# Dependencies
# ============================================================================
upgrade-deps:
	@printf "\n"
	@printf "$(BOLD)$(CYAN)⬆️  Upgrading dependencies...$(RESET)\n"
	@go get -u ./...
	@go mod tidy
	@printf "$(GREEN)✅ Dependencies upgraded!$(RESET)\n"

tidy-deps:
	@printf "\n"
	@printf "$(BOLD)$(CYAN)🧹 Tidying dependencies...$(RESET)\n"
	@go mod tidy
	@printf "$(GREEN)✅ Dependencies tidied!$(RESET)\n"

# ============================================================================
# Clean
# ============================================================================
clean:
	@printf "\n"
	@printf "$(BOLD)$(CYAN)🧹 Cleaning build artifacts...$(RESET)\n"
	@rm -f coverage.out coverage.html
	@go clean -cache -testcache
	@printf "$(GREEN)✅ Cleaned!$(RESET)\n"

# ============================================================================
# Tools
# ============================================================================
sandbox:
	@printf "\n"
	@printf "$(BOLD)$(CYAN)🎮 Starting Sandbox...$(RESET)\n"
	@go run cmd/sandbox/main.go

# ============================================================================
# System info
# ============================================================================
system-info:
	@printf "\n"
	@printf "$(BOLD)$(CYAN)📊 System Information$(RESET)\n"
	@printf "$(CYAN)━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━$(RESET)\n"
	@printf "$(YELLOW)OS:$(RESET)              $(OS) ($(UNAME_S))\n"
	@printf "$(YELLOW)Go Version:$(RESET)      $(shell go version 2>/dev/null || echo 'not installed')\n"
	@printf "$(YELLOW)golangci-lint:$(RESET)   $(shell golangci-lint --version 2>/dev/null | head -1 || echo 'not installed')\n"
	@printf "$(YELLOW)govulncheck:$(RESET)     $(shell govulncheck -version 2>/dev/null | head -1 || echo 'not installed')\n"
	@printf "$(YELLOW)GOPATH:$(RESET)          $(shell go env GOPATH 2>/dev/null || echo 'not set')\n"
	@printf "$(YELLOW)GOBIN:$(RESET)           $(GOBIN)\n"
	@printf "$(CYAN)━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━$(RESET)\n"

# ============================================================================
# Help command
# ============================================================================
help:
	@printf "\n"
	@printf "$(BOLD)$(CYAN)━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━$(RESET)\n"
	@printf "$(BOLD)$(CYAN)           WPD Message Gateway - Makefile                   $(RESET)\n"
	@printf "$(BOLD)$(CYAN)━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━$(RESET)\n"
	@printf "\n"
	@printf "$(BOLD)$(GREEN)🚀 SETUP$(RESET)\n"
	@printf "   $(YELLOW)make setup$(RESET)            Complete setup (tools + go modules)\n"
	@printf "   $(YELLOW)make setup-tools$(RESET)      Install development tools\n"
	@printf "   $(YELLOW)make setup-go$(RESET)         Set up Go modules\n"
	@printf "   $(YELLOW)make system-info$(RESET)      Show system information\n"
	@printf "\n"
	@printf "$(BOLD)$(GREEN)🔨 BUILD$(RESET)\n"
	@printf "   $(YELLOW)make build$(RESET)            Build all packages\n"
	@printf "   $(YELLOW)make fmt$(RESET)              Format all code\n"
	@printf "   $(YELLOW)make clean$(RESET)            Clean build artifacts\n"
	@printf "\n"
	@printf "$(BOLD)$(GREEN)🧪 TESTING$(RESET)\n"
	@printf "   $(YELLOW)make test$(RESET)             Run all tests (verbose)\n"
	@printf "   $(YELLOW)make test-short$(RESET)       Run tests (short mode)\n"
	@printf "   $(YELLOW)make test-cover$(RESET)       Run tests with coverage\n"
	@printf "   $(YELLOW)make test-cover-html$(RESET)  Generate HTML coverage report\n"
	@printf "\n"
	@printf "$(BOLD)$(GREEN)🛠️ TOOLS$(RESET)\n"
	@printf "   $(YELLOW)make sandbox$(RESET)          Run interactive sandbox\n"
	@printf "\n"
	@printf "$(BOLD)$(GREEN)🔍 QUALITY$(RESET)\n"
	@printf "   $(YELLOW)make lint$(RESET)             Run linter\n"
	@printf "   $(YELLOW)make vulncheck$(RESET)        Check for vulnerabilities\n"
	@printf "   $(YELLOW)make audit$(RESET)            Run all checks (fmt, lint, test, vuln)\n"
	@printf "\n"
	@printf "$(BOLD)$(GREEN)📦 DEPENDENCIES$(RESET)\n"
	@printf "   $(YELLOW)make upgrade-deps$(RESET)     Upgrade all dependencies\n"
	@printf "   $(YELLOW)make tidy-deps$(RESET)        Tidy dependencies\n"
	@printf "\n"
	@printf "$(BOLD)$(CYAN)━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━$(RESET)\n"
	@printf "$(BOLD)$(MAGENTA)💡 Quick Start:$(RESET) make setup → make test → make audit\n"
	@printf "$(BOLD)$(CYAN)━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━$(RESET)\n"
	@printf "\n"
