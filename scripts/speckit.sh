#!/usr/bin/env bash

set -euo pipefail

usage() {
  cat <<'EOF'
Usage:
  ./scripts/speckit.sh <command> [args...]

Commands:
  specify        /speckit.specify <feature description>
  clarify        /speckit.clarify
  plan           /speckit.plan
  tasks          /speckit.tasks
  implement      /speckit.implement
  analyze        /speckit.analyze
  checklist      /speckit.checklist <domain>
  constitution   /speckit.constitution <principles...>
  taskstoissues  /speckit.taskstoissues

Behavior:
  - Prints the slash command to run in your AI agent chat.
  - On macOS, copies it to clipboard via pbcopy (if available).

Examples:
  make specify FEATURE="Add SendGrid provider"
  make checklist DOMAIN="security"
EOF
}

is_macos() {
  [[ "$(uname -s)" == "Darwin" ]]
}

trim() {
  # Trim leading/trailing whitespace.
  # shellcheck disable=SC2001
  echo "$1" | sed -e 's/^[[:space:]]*//' -e 's/[[:space:]]*$//'
}

copy_to_clipboard_if_possible() {
  local text="$1"
  if is_macos && command -v pbcopy >/dev/null 2>&1; then
    printf "%s" "$text" | pbcopy
    printf "✓ Copied to clipboard.\n"
  fi
}

prompt_multiline_hint() {
  local title="$1"
  local hint="$2"
  printf "\n%s\n" "$title"
  printf "%s\n" "$hint"
  printf "\n"
}

prompt_line() {
  local prompt="$1"
  local value=""
  while true; do
    read -r -p "$prompt" value
    value="$(trim "$value")"
    if [[ -n "$value" ]]; then
      printf "%s" "$value"
      return 0
    fi
    printf "Please enter a non-empty value.\n"
  done
}

main() {
  local cmd="${1:-}"
  shift || true

  if [[ -z "$cmd" ]]; then
    usage
    exit 1
  fi

  local slash=""
  case "$cmd" in
    specify)        slash="/speckit.specify" ;;
    clarify)        slash="/speckit.clarify" ;;
    plan)           slash="/speckit.plan" ;;
    tasks)          slash="/speckit.tasks" ;;
    implement)      slash="/speckit.implement" ;;
    analyze)        slash="/speckit.analyze" ;;
    checklist)      slash="/speckit.checklist" ;;
    constitution)   slash="/speckit.constitution" ;;
    taskstoissues)  slash="/speckit.taskstoissues" ;;
    -h|--help|help) usage; exit 0 ;;
    *)
      printf "Unknown command: %s\n\n" "$cmd" >&2
      usage
      exit 1
      ;;
  esac

  # Preserve spacing; avoid shell escaping tricks—user is pasting into chat.
  local raw_args=""
  raw_args="$(trim "${*:-}")"

  # Interactive prompts for commands that require user input.
  case "$cmd" in
    specify)
      if [[ -z "$raw_args" ]]; then
        prompt_multiline_hint \
          "Spec Kit — Specify" \
          "Write a 1–3 sentence feature description.\n\nGood examples:\n- \"Add a SendGrid email provider integration\"\n- \"Add rate limiting to /v1/email per workspace\"\n- \"Portal: add a workspace icon picker\"\n\nAvoid implementation details (no file paths, no DB column names) — that comes in plan/tasks."
        raw_args="$(prompt_line "Feature description: ")"
      fi
      ;;
    checklist)
      if [[ -z "$raw_args" ]]; then
        prompt_multiline_hint \
          "Spec Kit — Checklist" \
          "Choose the checklist focus area (requirements quality, not implementation testing).\n\nCommon domains:\n- security\n- api\n- ux\n- performance"
        raw_args="$(prompt_line "Checklist domain: ")"
      fi
      ;;
  esac

  local args=""
  if [[ -n "$raw_args" ]]; then
    args=" $raw_args"
  fi

  local full="${slash}${args}"

  printf "\nRun this in your AI agent chat:\n\n"
  printf "%s\n\n" "$full"
  copy_to_clipboard_if_possible "$full"
}

main "$@"

