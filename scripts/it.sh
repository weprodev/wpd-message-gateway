#!/usr/bin/env bash

set -euo pipefail

REPO_ROOT="$(CDPATH="" cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
source "$REPO_ROOT/.specify/scripts/bash/common.sh"
source "$REPO_ROOT/scripts/lib/install-helpers.sh" 2>/dev/null || true

SETUP_CONFIG="$REPO_ROOT/.specify/memory/setup.json"
SETUP_OPTIONS="$REPO_ROOT/scripts/setup-options.json"

# Colors
BOLD="\033[1m"
CYAN="\033[36m"
GREEN="\033[32m"
YELLOW="\033[33m"
MAGENTA="\033[35m"
RESET="\033[0m"

is_macos() {
  [[ "$(uname -s)" == "Darwin" ]]
}

copy_to_clipboard_if_possible() {
  local text="$1"
  if is_macos && command -v pbcopy >/dev/null 2>&1; then
    printf "%s" "$text" | pbcopy
    printf "  ${MAGENTA}✓ Copied to clipboard!${RESET}\n"
  fi
}

prompt_yes_no() {
  local prompt="$1"
  local default="${2:-yes}"
  
  local hint="[Y/n]"
  if [[ "$default" == "no" ]]; then
    hint="[y/N]"
  fi

  local answer
  while true; do
    printf "%s %s " "$prompt" "$hint" >&2
    read -r answer < /dev/tty || true
    if [[ -z "$answer" ]]; then
      answer="$default"
    fi
    case "${answer}" in
      [yY]|[yY][eE][sS]) return 0 ;;
      [nN]|[nN][oO])     return 1 ;;
      *)                 printf "${YELLOW}Please answer yes or no.${RESET}\n" >&2 ;;
    esac
  done
}

ensure_jq() {
  if ! command -v jq >/dev/null 2>&1; then
    printf "${YELLOW}jq is required for orchestration but not found.${RESET}\n"
    if prompt_yes_no "Do you want to install jq now (via Homebrew)?" "yes"; then
      if command -v brew >/dev/null 2>&1; then
        brew install jq
      else
        printf "${RED}Homebrew not found. Please run 'make setup' first.${RESET}\n"
        exit 1
      fi
    else
      printf "Please run 'make setup' to configure your workspace properly.\n"
      exit 1
    fi
  fi
}

get_configured_tool() {
  local role="$1" # master, backend-engineer, reviewer, etc.
  ensure_jq

  if [[ ! -f "$SETUP_CONFIG" ]]; then
    printf "${YELLOW}setup.json not found. Run 'make setup' first.${RESET}\n" >&2
    return 1
  fi

  local tool_id
  # Master role is a top-level key; all others are under .roles
  if [[ "$role" == "master" ]]; then
    tool_id="$(jq -r '.master.tool // empty' "$SETUP_CONFIG")"
  else
    tool_id="$(jq -r ".roles[\"$role\"].tool // empty" "$SETUP_CONFIG")"
    # Fall back to master tool if role not configured
    if [[ -z "$tool_id" ]]; then
      tool_id="$(jq -r '.master.tool // empty' "$SETUP_CONFIG")"
    fi
  fi

  printf "%s" "$tool_id"
}

get_configured_model() {
  local role="$1" # master, backend-engineer, reviewer, etc.
  ensure_jq

  if [[ ! -f "$SETUP_CONFIG" ]]; then
    return 1
  fi

  local model_id
  if [[ "$role" == "master" ]]; then
    model_id="$(jq -r '.master.model // empty' "$SETUP_CONFIG")"
  else
    model_id="$(jq -r ".roles[\"$role\"].model // empty" "$SETUP_CONFIG")"
    if [[ -z "$model_id" ]]; then
      model_id="$(jq -r '.master.model // empty' "$SETUP_CONFIG")"
    fi
  fi
  # Final fallback to default_model
  if [[ -z "$model_id" ]]; then
    model_id="$(jq -r '.default_model // empty' "$SETUP_CONFIG")"
  fi

  printf "%s" "$model_id"
}

get_tool_command() {
  local tool_id="$1"
  ensure_jq

  if [[ ! -f "$SETUP_OPTIONS" ]]; then
    printf "${YELLOW}setup-options.json not found.${RESET}\n" >&2
    return 1
  fi

  # Find the command associated with the tool id
  local cmd_name
  cmd_name="$(jq -r ".installable_tools[] | select(.id == \"$tool_id\") | .command // empty" "$SETUP_OPTIONS")"
  printf "%s" "$cmd_name"
}

launch_agent() {
  local role="$1"
  local prompt_str="$2"
  
  local tool_id
  tool_id="$(get_configured_tool "$role")"
  
  if [[ -z "$tool_id" || "$tool_id" == "null" ]]; then
    printf "${YELLOW}No tool configured for role '${role}'. Please run 'make setup'.${RESET}\n"
    return 1
  fi

  local model_id
  model_id="$(get_configured_model "$role")"

  local cmd_name
  cmd_name="$(get_tool_command "$tool_id")"

  if [[ -z "$cmd_name" ]]; then
    printf "${YELLOW}Could not resolve executable command for tool '${tool_id}'.${RESET}\n"
    return 1
  fi

  if ! command -v "$cmd_name" >/dev/null 2>&1; then
    printf "${YELLOW}Executable '${cmd_name}' not found in PATH.${RESET}\n"
    printf "Please ensure it is installed correctly.\n"
    return 1
  fi

  printf "\n"
  printf "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${RESET}\n"
  printf " ${BOLD}Launching AI Agent (${cmd_name})${RESET}\n"
  printf "   ${BOLD}Role:${RESET}  $role\n"
  printf "   ${BOLD}Tool:${RESET}  $tool_id\n"
  printf "   ${BOLD}Model:${RESET} $model_id\n"
  printf "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${RESET}\n"
  
  copy_to_clipboard_if_possible "$prompt_str"
  printf "${YELLOW}Paste the command into the agent once it starts.${RESET}\n\n"

  printf "\n${GREEN}Starting interaction... Type /exit or use native CLI exit commands to return to the menu.${RESET}\n\n"
  
  # Resolve model flag per tool (each CLI has a different flag convention)
  local model_flag=""
  case "$tool_id" in
    gemini) model_flag="-m" ;;
    claude) model_flag="--model" ;;
    cursor) model_flag="--model" ;;
  esac

  if [[ -n "$model_flag" && -n "$model_id" ]]; then
    "$cmd_name" "$model_flag" "$model_id" < /dev/tty
  else
    "$cmd_name" < /dev/tty
  fi
  
  printf "\n${CYAN}Agent session ended.${RESET}\n"
}

print_menu() {
  clear
  printf "\n"
  printf "${BOLD}${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${RESET}\n"
  printf " ${BOLD}Spec Kit — Open Terminal Orchestrator${RESET}\n"
  printf "${BOLD}${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${RESET}\n"
  printf " Use this menu to step through your feature pipeline.\n\n"
  
  printf "  ${BOLD}1)${RESET} Specify Feature \t ${YELLOW}(Kickoff)${RESET}\n"
  printf "  ${BOLD}2)${RESET} Plan            \t ${YELLOW}(Architecture)${RESET}\n"
  printf "  ${BOLD}3)${RESET} Tasks           \t ${YELLOW}(Breakdown)${RESET}\n"
  printf "  ${BOLD}4)${RESET} Implement       \t ${YELLOW}(Code Writing)${RESET}\n"
  printf "  ${BOLD}5)${RESET} Analyze / QA    \t ${YELLOW}(Review)${RESET}\n"
  printf "  ${BOLD}6)${RESET} Finish & PR     \t ${YELLOW}(Sync to GitHub)${RESET}\n"
  printf "\n"
  printf "  ${BOLD}Q)${RESET} Quit\n"
  printf "\n"
}

main() {
  ensure_jq

  while true; do
    print_menu
    local choice
    read -r -p "Select next step: " choice < /dev/tty || exit 0
    
    printf "\n"
    
    case "$choice" in
      1)
        printf "${BOLD}${GREEN}==> Run Kickoff${RESET}\n"
        make specify
        launch_agent "master" "/speckit.specify"
        ;;
      2)
        printf "${BOLD}${GREEN}==> Run Plan${RESET}\n"
        launch_agent "master" "/speckit.plan"
        ;;
      3)
        printf "${BOLD}${GREEN}==> Break down Tasks${RESET}\n"
        launch_agent "master" "/speckit.tasks"
        ;;
      4)
        printf "${BOLD}${GREEN}==> Implement${RESET}\n"
        launch_agent "master" "/speckit.implement"
        ;;
      5)
        printf "${BOLD}${GREEN}==> Analyze / Quality Assurance${RESET}\n"
        launch_agent "master" "/speckit.analyze"
        ;;
      6)
        printf "${BOLD}${GREEN}==> Finish & Pull Request${RESET}\n"
        make pr
        printf "\nPress Enter to return to menu..."
        read -r < /dev/tty || true
        ;;
      q|Q|quit|exit)
        printf "Goodbye!\n"
        exit 0
        ;;
      *)
        printf "${YELLOW}Invalid option.${RESET}\n"
        sleep 1
        ;;
    esac
  done
}

main "$@"
