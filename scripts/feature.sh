#!/usr/bin/env bash

set -euo pipefail

REPO_ROOT="$(CDPATH="" cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SPECS_DIR="$REPO_ROOT/specs"

SETUP_CONFIG="$REPO_ROOT/.specify/memory/setup.json"

# Colors
BOLD="\033[1m"
MAGENTA="\033[35m"
CYAN="\033[36m"
GREEN="\033[32m"
YELLOW="\033[33m"
RESET="\033[0m"

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

usage() {
  cat <<'EOF'
Usage:
  ./scripts/feature.sh <command>

Commands:
  specify   Create branch + spec folder, then (if gh authed) create GitHub issue
  sync      Sync local specs (spec.md, plan.md, tasks.md) to the linked GitHub issue body
  pr        Run sync, then create GitHub PR linked to the issue

Notes:
  - macOS only assumptions are avoided where possible, but clipboard copy uses pbcopy when available.
  - GitHub operations require: gh auth login
EOF
}

is_git_repo() {
  git rev-parse --is-inside-work-tree >/dev/null 2>&1
}

is_gh_authed() {
  command -v gh >/dev/null 2>&1 && gh auth status >/dev/null 2>&1
}

copy_to_clipboard_if_possible() {
  local text="$1"
  if [[ "$(uname -s)" == "Darwin" ]] && command -v pbcopy >/dev/null 2>&1; then
    printf "%s" "$text" | pbcopy
  fi
}

trim() {
  echo "$1" | sed -e 's/^[[:space:]]*//' -e 's/[[:space:]]*$//'
}

prompt_nonempty() {
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

read_default_branch() {
  # Default to master for this repo unless setup.json overrides it.
  local default_branch="master"
  if [[ -f "$SETUP_CONFIG" ]] && command -v python3 >/dev/null 2>&1; then
    local parsed
    parsed="$(python3 - <<PY 2>/dev/null || true
import json, pathlib
p = pathlib.Path(r"$SETUP_CONFIG")
try:
  data = json.loads(p.read_text())
  print(data.get("default_branch", ""))
except Exception:
  pass
PY
)"
    parsed="$(trim "${parsed:-}")"
    if [[ -n "$parsed" ]]; then
      default_branch="$parsed"
    fi
  fi
  printf "%s" "$default_branch"
}

git_sync_default_branch() {
  local default_branch="$1"

  git fetch --all --prune

  local current_branch
  current_branch="$(git branch --show-current)"

  if [[ "$current_branch" != "$default_branch" ]]; then
    # Create from origin if it doesn't exist locally
    if ! git show-ref --quiet "refs/heads/$default_branch"; then
      if git show-ref --quiet "refs/remotes/origin/$default_branch"; then
        git checkout -b "$default_branch" "origin/$default_branch"
      else
        printf "\n${YELLOW}WARNING: Default branch '%s' does not exist locally or on origin.${RESET}\n" "$default_branch" >&2
        printf "Please create it or select a valid default branch.\n" >&2
        exit 1
      fi
    else
      git checkout "$default_branch"
    fi
  fi

  # Only attempt pull if tracking branch exists.
  if git rev-parse --abbrev-ref --symbolic-full-name '@{u}' >/dev/null 2>&1; then
    git pull --ff-only || true
  fi
}

create_feature_branch_and_spec() {
  local feature_description="$1"

  # Use Spec Kit's branch+spec initializer to keep numbering consistent.
  local out
  out="$("$REPO_ROOT/.specify/scripts/bash/create-new-feature.sh" --json "$feature_description")"

  if command -v python3 >/dev/null 2>&1; then
    python3 - <<PY
import json
data = json.loads(r'''$out''')
print(data["BRANCH_NAME"])
print(data["SPEC_FILE"])
PY
  else
    # Very small fallback (not robust, but avoids hard fail).
    echo "$out" | sed -n 's/.*"BRANCH_NAME":"\([^"]*\)".*/\1/p'
    echo "$out" | sed -n 's/.*"SPEC_FILE":"\([^"]*\)".*/\1/p'
  fi
}

write_feature_meta() {
  local feature_dir="$1"
  local issue_number="${2:-}"
  local issue_url="${3:-}"

  mkdir -p "$feature_dir"
  cat > "$feature_dir/meta.json" <<EOF
{
  "issue_number": "${issue_number}",
  "issue_url": "${issue_url}"
}
EOF
}

make_specify() {
  if ! is_git_repo; then
    printf "ERROR: Not a git repository.\n" >&2
    exit 1
  fi

  printf "\nSpec Kit — Feature kickoff\n\n"
  printf "Write a 1–3 sentence feature description.\n"
  printf "Good: \"Portal: add workspace icon picker\"\n"
  printf "Avoid: file paths / DB column names (that comes later).\n\n"

  local desc
  desc="$(prompt_nonempty "Feature description: ")"

  local default_branch
  default_branch="$(read_default_branch)"

  local current_branch
  current_branch="$(git branch --show-current)"
  
  local has_uncommitted=0
  if ! git diff-index --quiet HEAD --; then
    has_uncommitted=1
  fi

  local move_changes=0
  if [[ $has_uncommitted -eq 1 ]]; then
    printf "\n${YELLOW}WARNING: You have uncommitted changes on branch '${BOLD}%s${RESET}${YELLOW}'.${RESET}\n" "$current_branch" >&2
    if prompt_yes_no "Do you want to temporarily stash and move these changes to the new feature branch?" "yes"; then
      git stash --include-untracked
      move_changes=1
    else
      printf "\n${BOLD}${CYAN}Tip:${RESET} Please commit or stash your changes natively before kicking off a new feature.\n" >&2
      exit 1
    fi
  elif [[ "$current_branch" != "$default_branch" ]]; then
    printf "\n${YELLOW}WARNING: You are on branch '${BOLD}%s${RESET}${YELLOW}', not the default branch ('${BOLD}%s${RESET}${YELLOW}').${RESET}\n" "$current_branch" "$default_branch" >&2
    if ! prompt_yes_no "Are you sure you want to sync and create the new feature branch from '${default_branch}'?" "yes"; then
      exit 1
    fi
  fi

  printf "\nSyncing default branch (%s)...\n" "$default_branch"
  git_sync_default_branch "$default_branch"

  printf "\nCreating feature branch + spec folder...\n"
  local parsed
  parsed="$(create_feature_branch_and_spec "$desc")"
  local branch spec_file
  branch="$(printf "%s" "$parsed" | sed -n '1p')"
  spec_file="$(printf "%s" "$parsed" | sed -n '2p')"

  local feature_dir="$SPECS_DIR/$branch"

  printf "✓ Branch: %s\n" "$branch"
  printf "✓ Spec:   %s\n" "$spec_file"

  if [[ $move_changes -eq 1 ]]; then
    printf "\n${CYAN}Applying stashed changes to '%s'...${RESET}\n" "$branch"
    git stash pop || true
  fi

  local issue_number="" issue_url=""
  if is_gh_authed; then
    printf "\nCreating GitHub issue...\n"
    # Use the description as title (trim to first line).
    local title
    title="$(printf "%s" "$desc" | head -n 1)"
    local body
    body="$(cat <<EOF
## Summary
$desc

## Artifacts
- Spec: \`$spec_file\`
- Branch: \`$branch\`
EOF
)"
    issue_url="$(gh issue create --title "$title" --body "$body" --json url,number --jq '.url')"
    issue_number="$(gh issue view "$issue_url" --json number --jq '.number')"
    printf "✓ Issue:  %s\n" "$issue_url"
  else
    printf "\nGitHub issue creation skipped (gh not authenticated).\n"
    printf "Run: gh auth login\n"
  fi

  write_feature_meta "$feature_dir" "$issue_number" "$issue_url"

  local next="/speckit.specify $desc"
  printf "\nNext (paste into your AI agent chat):\n\n%s\n\n" "$next"
  copy_to_clipboard_if_possible "$next"
}

make_sync() {
  if ! is_gh_authed; then
    printf "ERROR: gh is not authenticated. Run: gh auth login\n" >&2
    exit 1
  fi

  local branch
  branch="$(git branch --show-current)"
  if [[ -z "$branch" ]]; then
    printf "ERROR: Unable to determine current branch.\n" >&2
    exit 1
  fi

  local feature_dir="$SPECS_DIR/$branch"
  local meta="$feature_dir/meta.json"
  local issue_number=""
  if [[ -f "$meta" ]] && command -v python3 >/dev/null 2>&1; then
    issue_number="$(python3 - <<PY 2>/dev/null || true
import json, pathlib
data=json.loads(pathlib.Path(r"$meta").read_text())
print(data.get("issue_number",""))
PY
)"
    issue_number="$(trim "${issue_number:-}")"
  fi

  if [[ -z "$issue_number" ]]; then
    printf "No linked GitHub issue found in %s. Skipping sync.\n" "$meta"
    return 0
  fi

  printf "Syncing local specs to GitHub Issue #%s...\n" "$issue_number"

  local tmp_body
  tmp_body="$(mktemp)"

  printf "## Specification\n\n" >> "$tmp_body"
  if [[ -f "$feature_dir/spec.md" ]]; then
    cat "$feature_dir/spec.md" >> "$tmp_body"
  else
    printf "*(No spec.md found)*\n" >> "$tmp_body"
  fi

  printf "\n\n## Implementation Plan\n\n" >> "$tmp_body"
  if [[ -f "$feature_dir/plan.md" ]]; then
    cat "$feature_dir/plan.md" >> "$tmp_body"
  else
    printf "*(No plan.md found)*\n" >> "$tmp_body"
  fi

  printf "\n\n## Tasks\n\n" >> "$tmp_body"
  if [[ -f "$feature_dir/tasks.md" ]]; then
    cat "$feature_dir/tasks.md" >> "$tmp_body"
  else
    printf "*(No tasks.md found)*\n" >> "$tmp_body"
  fi

  gh issue edit "$issue_number" --body-file "$tmp_body"
  rm -f "$tmp_body"
  
  printf "✓ Successfully synced specs to Issue #%s\n" "$issue_number"
}

make_pr() {
  if ! is_gh_authed; then
    printf "ERROR: gh is not authenticated. Run: gh auth login\n" >&2
    exit 1
  fi

  local branch
  branch="$(git branch --show-current)"
  if [[ -z "$branch" ]]; then
    printf "ERROR: Unable to determine current branch.\n" >&2
    exit 1
  fi

  local feature_dir="$SPECS_DIR/$branch"
  local meta="$feature_dir/meta.json"
  local issue_number=""
  if [[ -f "$meta" ]] && command -v python3 >/dev/null 2>&1; then
    issue_number="$(python3 - <<PY 2>/dev/null || true
import json, pathlib
data=json.loads(pathlib.Path(r"$meta").read_text())
print(data.get("issue_number",""))
PY
)"
    issue_number="$(trim "${issue_number:-}")"
  fi

  if [[ -n "$issue_number" ]]; then
    make_sync || true
  fi

  if ! git rev-parse --abbrev-ref --symbolic-full-name '@{u}' >/dev/null 2>&1; then
    git push -u origin HEAD
  else
    git push
  fi

  local title="$branch"
  local body="## Summary\n\nImplements work for \`$branch\`.\n"
  if [[ -n "$issue_number" ]]; then
    body+="\nCloses #$issue_number\n"
    title="[$issue_number] $branch"
  fi

  gh pr create --title "$title" --body "$body"
}

main() {
  local cmd="${1:-}"
  case "$cmd" in
    specify) make_specify ;;
    sync)    make_sync ;;
    pr)      make_pr ;;
    -h|--help|help|"") usage ;;
    *)
      printf "Unknown command: %s\n\n" "$cmd" >&2
      usage
      exit 1
      ;;
  esac
}

main "${1:-}"

