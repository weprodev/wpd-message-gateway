#!/usr/bin/env bash

set -euo pipefail

REPO_ROOT="$(CDPATH="" cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
CONFIG_DIR="$REPO_ROOT/.specify/memory"
CONFIG_FILE="$CONFIG_DIR/setup.json"

# Colors
BOLD="\033[1m"
CYAN="\033[36m"
GREEN="\033[32m"
YELLOW="\033[33m"
RESET="\033[0m"

is_macos() {
  [[ "$(uname -s)" == "Darwin" ]]
}

print_header() {
  printf "\n"
  printf "${BOLD}${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${RESET}\n"
  printf " ${BOLD}Message Gateway — Interactive Setup (macOS)${RESET}\n"
  printf "${BOLD}${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${RESET}\n"
  printf "\n"
}

prompt_yes_no() {
  local prompt="$1"
  local default="${2:-yes}" # yes|no
  local answer=""

  while true; do
    if [[ "$default" == "yes" ]]; then
      read -r -p "$prompt [Y/n] " answer < /dev/tty
      answer="${answer:-Y}"
    else
      read -r -p "$prompt [y/N] " answer < /dev/tty
      answer="${answer:-N}"
    fi

    case "$answer" in
      [yY]|[yY][eE][sS]) return 0 ;;
      [nN]|[nN][oO])     return 1 ;;
      *)                 printf "Please answer yes or no.\n" >&2 ;;
    esac
  done
}

prompt_select() {
  local prompt="$1"
  shift
  local options=("$@")

  # Modern interactive terminal UI if gum is available
  if command -v gum >/dev/null 2>&1; then
    local choice
    # gum renders to stderr/tty natively and prints selection to stdout
    if choice="$(gum choose --header="$prompt" "${options[@]}" < /dev/tty)"; then
      if [[ -n "$choice" ]]; then
        printf "%s" "$choice"
        printf "  ${GREEN}✓ %s${RESET}\n" "$choice" >&2
        return 0
      fi
    # If gum is interrupted by Ctrl+C inside subshell, exit whole script
    else
      exit 1
    fi
  else
    # Fallback to standard CLI numeric selection
    printf "\n%s\n" "$prompt" >&2
    local i=1
    for opt in "${options[@]}"; do
      printf "  %d) %s\n" "$i" "$opt" >&2
      i=$((i + 1))
    done

    local choice=""
    while true; do
      read -r -p "Select an option [1-${#options[@]}]: " choice < /dev/tty
      if [[ "$choice" =~ ^[0-9]+$ ]] && (( choice >= 1 && choice <= ${#options[@]} )); then
        printf "%s" "${options[$((choice - 1))]}"
        return 0
      fi
      printf "Invalid selection.\n" >&2
    done
  fi
}

require_cmd() {
  local cmd="$1"
  if ! command -v "$cmd" >/dev/null 2>&1; then
    printf "Missing required command: %s\n" "$cmd" >&2
    return 1
  fi
  return 0
}

source "$REPO_ROOT/scripts/lib/install-helpers.sh"

write_config() {
  local default_branch="$1"
  local default_model="$2"
  local default_tools_json="$3"
  local default_tools_model="$4"
  local reviewer_tools_json="$5"
  local reviewer_tools_model="$6"
  local testdocs_tools_json="$7"
  local testdocs_tools_model="$8"

  mkdir -p "$CONFIG_DIR"

  cat > "$CONFIG_FILE" <<EOF
{
  "platform": "macos",
  "default_branch": "${default_branch}",
  "tools": {
    "default": { "tools": ${default_tools_json}, "model": "${default_tools_model}" },
    "reviewer": { "tools": ${reviewer_tools_json}, "model": "${reviewer_tools_model}" },
    "test_docs": { "tools": ${testdocs_tools_json}, "model": "${testdocs_tools_model}" }
  },
  "default_model": "${default_model}"
}
EOF

  printf "\n${GREEN}✓ Wrote configuration: %s${RESET}\n" "$CONFIG_FILE"
}



main() {
  if ! is_macos; then
    printf "This setup script currently supports macOS only.\n" >&2
    exit 1
  fi

  print_header

  printf "Verifying required setup engine frameworks (Homebrew will be used if missing)...\n" >&2
  install_via_brew "jq (JSON Processor)" "jq"
  install_via_brew "gum (UI Components)" "gum"
  printf "\n" >&2

  local options_file="$REPO_ROOT/scripts/setup-options.json"
  if [[ ! -f "$options_file" ]]; then
    printf "Error: options config not found at %s\n" "$options_file" >&2
    exit 1
  fi

  printf "${BOLD}${YELLOW}A) Tools configuration${RESET}\n"
  printf "\n"

  while IFS= read -r tool_payload; do
    local tid tdisplay ttype tpkg tscript tcmd
    tid="$(echo "$tool_payload" | jq -r '.id')"
    tdisplay="$(echo "$tool_payload" | jq -r '.display')"
    ttype="$(echo "$tool_payload" | jq -r '.type')"
    tpkg="$(echo "$tool_payload" | jq -r '.package // empty')"
    tscript="$(echo "$tool_payload" | jq -r '.script // empty')"
    tcmd="$(echo "$tool_payload" | jq -r '.command // empty')"
    
    if prompt_yes_no "Enable $tdisplay?" "yes"; then
      case "$ttype" in
        "brew")
          install_via_brew "$tdisplay" "$tcmd" "$tpkg" || true
          ;;
        "npm")
          install_via_npm "$tdisplay" "$tcmd" "$tpkg" || true
          ;;
        "custom")
          install_custom_script "$tdisplay" "$tcmd" "$tscript" || true
          ;;
        *)
          printf "Unknown tool installation type '%s' for '%s'.\n" "$ttype" "$tid" >&2
          ;;
      esac
    fi
    printf "\n" >&2
  done < <(jq -c '.installable_tools[]' "$options_file")

  printf "\n" >&2
  printf "${BOLD}${YELLOW}B) Project Configuration${RESET}\n"
  printf "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n"

  local -a branches=()
  while IFS= read -r line; do
    if [[ -n "$line" ]]; then branches+=("$line"); fi
  done < <(jq -r '.branches[]' "$options_file")
  branches+=("custom")

  local -a models=()
  while IFS= read -r line; do
    if [[ -n "$line" ]]; then models+=("$line"); fi
  done < <(jq -r '.models[]' "$options_file")
  models+=("custom")

  local -a tools=()
  while IFS= read -r line; do
    if [[ -n "$line" ]]; then tools+=("$line"); fi
  done < <(jq -r '.tools[]' "$options_file")

  local -a tools_models=("inherit(default)")
  for m in "${models[@]}"; do
    if [[ "$m" != "custom" ]]; then
      tools_models+=("$m")
    fi
  done
  tools_models+=("custom")

  local default_branch
  default_branch="$(prompt_select "Default git branch name for this repo:" "${branches[@]}")"
  if [[ "$default_branch" == "custom" ]]; then
    read -r -p "Enter default branch name: " default_branch < /dev/tty
    default_branch="${default_branch:-master}"
  fi

  printf "\n" >&2
  printf "${BOLD}${YELLOW}C) Base Model Configuration${RESET}\n"
  printf "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n"

  local default_model
  default_model="$(prompt_select "Choose base infrastructure model:" "${models[@]}")"

  if [[ "$default_model" == "custom" ]]; then
    read -r -p "Enter custom default model string: " default_model < /dev/tty
    default_model="${default_model:-gpt-5}"
  fi

  printf "\n" >&2
  printf "${BOLD}${YELLOW}D) Role Agents Setup${RESET}\n"
  printf "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n" >&2
  
  printf "${CYAN}--- [1/3] Default Role ---${RESET}\n" >&2
  local default_tools_model reviewer_tools_model testdocs_tools_model
  
  local default_tool
  default_tool="$(prompt_select "Choose tool for default role:" "${tools[@]}")"
  local default_tools_json="[\"${default_tool}\"]"
  
  default_tools_model="$(prompt_select "Choose model for default tools:" "${tools_models[@]}")"
  if [[ "$default_tools_model" == "inherit(default)" ]]; then default_tools_model="$default_model"; fi
  if [[ "$default_tools_model" == "custom" ]]; then
    read -r -p "Enter custom model for default tools: " default_tools_model < /dev/tty
    default_tools_model="${default_tools_model:-$default_model}"
  fi

  printf "\n${CYAN}--- [2/3] Reviewer Role ---${RESET}\n" >&2
  local reviewer_tool
  reviewer_tool="$(prompt_select "Choose tool for reviewer role:" "${tools[@]}")"
  local reviewer_tools_json="[\"${reviewer_tool}\"]"

  reviewer_tools_model="$(prompt_select "Choose model for reviewer tools:" "${tools_models[@]}")"
  if [[ "$reviewer_tools_model" == "inherit(default)" ]]; then reviewer_tools_model="$default_model"; fi
  if [[ "$reviewer_tools_model" == "custom" ]]; then
    read -r -p "Enter custom model for reviewer tools: " reviewer_tools_model < /dev/tty
    reviewer_tools_model="${reviewer_tools_model:-$default_model}"
  fi

  printf "\n${CYAN}--- [3/3] Test & Docs Role ---${RESET}\n" >&2
  local testdocs_tool
  testdocs_tool="$(prompt_select "Choose tool for test/docs role:" "${tools[@]}")"
  local testdocs_tools_json="[\"${testdocs_tool}\"]"

  testdocs_tools_model="$(prompt_select "Choose model for test/docs tools:" "${tools_models[@]}")"
  if [[ "$testdocs_tools_model" == "inherit(default)" ]]; then testdocs_tools_model="$default_model"; fi
  if [[ "$testdocs_tools_model" == "custom" ]]; then
    read -r -p "Enter custom model for test/docs tools: " testdocs_tools_model < /dev/tty
    testdocs_tools_model="${testdocs_tools_model:-$default_model}"
  fi

  write_config \
    "$default_branch" \
    "$default_model" \
    "$default_tools_json" "$default_tools_model" \
    "$reviewer_tools_json" "$reviewer_tools_model" \
    "$testdocs_tools_json" "$testdocs_tools_model"

  printf "\n${BOLD}${GREEN}Done.${RESET}\n"
}

main "$@"

