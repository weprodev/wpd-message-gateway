#!/usr/bin/env bash

set -euo pipefail

# Colors
GREEN="\033[32m"
YELLOW="\033[33m"
RESET="\033[0m"

ensure_homebrew() {
  if command -v brew >/dev/null 2>&1; then
    return 0
  fi

  printf "\nHomebrew is not installed.\n"
  printf "Install it from https://brew.sh then re-run: make setup\n\n"
  return 1
}

ensure_npm() {
  if command -v npm >/dev/null 2>&1; then
    return 0
  fi

  printf "npm is required for this installation.\n" >&2
  printf "Install Node.js (e.g. via Homebrew: brew install node) then re-run: make setup\n" >&2
  return 1
}

install_via_brew() {
  local display="$1"
  local cmd="$2"
  local formula="${3:-$cmd}"

  if command -v "$cmd" >/dev/null 2>&1; then
    printf "${GREEN}✓${RESET} %s already installed.\n" "$display"
    return 0
  fi

  ensure_homebrew

  printf "Installing %s...\n" "$display"
  if ! brew install "$formula"; then
    printf "Failed to install '%s' via brew.\n" "$formula" >&2
    printf "You can install it manually, then re-run: make setup\n" >&2
    return 1
  fi
  printf "${GREEN}✓${RESET} Installed %s.\n" "$display"
}

install_via_npm() {
  local display="$1"
  local cmd="$2"
  local package="${3:-$cmd}"

  if command -v "$cmd" >/dev/null 2>&1; then
    printf "${GREEN}✓${RESET} %s already installed.\n" "$display"
    return 0
  fi

  ensure_npm

  printf "Installing %s via npm...\n" "$display"
  if ! npm install -g "$package"; then
    printf "Failed to install '%s' via npm.\n" "$package" >&2
    return 1
  fi
  
  if command -v "$cmd" >/dev/null 2>&1; then
    printf "${GREEN}✓${RESET} Installed %s.\n" "$display"
  else
    printf "Installation completed, but '%s' is not on PATH.\n" "$cmd" >&2
    return 1
  fi
}

install_custom_script() {
  local display="$1"
  local cmd="$2"
  local script_name="$3"

  local script_path="$REPO_ROOT/scripts/installers/$script_name"

  if command -v "$cmd" >/dev/null 2>&1; then
    printf "${GREEN}✓${RESET} %s already installed.\n" "$display"
    return 0
  fi

  if [[ ! -f "$script_path" ]]; then
    printf "Custom installer '%s' not found.\n" "$script_path" >&2
    return 1
  fi

  printf "Running custom installer for %s...\n" "$display"
  bash "$script_path"

  if command -v "$cmd" >/dev/null 2>&1; then
    printf "✓ Installed %s.\n" "$display"
  else
    printf "Custom installation finished, but '%s' is not on PATH.\n" "$cmd" >&2
    return 1
  fi
}
