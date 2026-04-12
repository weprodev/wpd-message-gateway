#!/usr/bin/env bash
# scripts/agents.sh — tmux multi-agent session manager
#
# Usage:
#   bash scripts/agents.sh start   — Create and attach to the agent session
#   bash scripts/agents.sh kill    — Destroy the active agent session
#   bash scripts/agents.sh status  — List panes and their roles

set -euo pipefail

REPO_ROOT="$(CDPATH="" cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
# Normalise to canonical filesystem case via os.listdir() path traversal.
# IMPORTANT: os.path.realpath() does NOT fix case on macOS HFS+ (case-insensitive
# filesystem). os.listdir() returns the real on-disk entry names, regardless of
# the case used when navigating (e.g. 'weprodev' → 'WeProDev').
if command -v python3 > /dev/null 2>&1; then
  _fix_case_script='import os, sys
p = os.path.abspath(sys.argv[1])
parts = []
while True:
    h, t = os.path.split(p)
    if h == p:
        parts.insert(0, p)
        break
    try:
        real = next((e for e in os.listdir(h) if e.lower() == t.lower()), t)
    except OSError:
        real = t
    parts.insert(0, real)
    p = h
print(os.path.join(*parts))'
  REPO_ROOT="$(python3 -c "$_fix_case_script" "$REPO_ROOT" 2>/dev/null || echo "$REPO_ROOT")"
  unset _fix_case_script
fi
SETUP_CONFIG="$REPO_ROOT/.specify/memory/setup.json"
SETUP_OPTIONS="$REPO_ROOT/scripts/setup-options.json"
ROLES_DIR="$REPO_ROOT/scripts/roles"
SIGNALS_DIR="$REPO_ROOT/.specify/signals"
SESSION="wpd-gateway"

# ── Colors ────────────────────────────────────────────────────────────────────
BOLD="\033[1m"
CYAN="\033[36m"
GREEN="\033[32m"
YELLOW="\033[33m"
RED="\033[31m"
RESET="\033[0m"

# ── Guards ────────────────────────────────────────────────────────────────────

require_tmux() {
  if ! command -v tmux > /dev/null 2>&1; then
    printf "${RED}Error: tmux is not installed.${RESET}\n"
    printf "Run ${YELLOW}make setup${RESET} to install it, or: brew install tmux\n"
    exit 1
  fi
}

require_config() {
  if [[ ! -f "$SETUP_CONFIG" ]]; then
    printf "${RED}Error: setup.json not found.${RESET}\n"
    printf "Run ${YELLOW}make setup${RESET} first to configure your roles and tools.\n"
    exit 1
  fi
  if ! command -v jq > /dev/null 2>&1; then
    printf "${RED}Error: jq is not installed.${RESET}\n"
    printf "Run ${YELLOW}make setup${RESET} first.\n"
    exit 1
  fi
}

# ── Config readers ────────────────────────────────────────────────────────────

get_tool_command() {
  local tool_id="$1"
  jq -r ".installable_tools[] | select(.id == \"$tool_id\") | .command // empty" "$SETUP_OPTIONS"
}

get_role_display() {
  local role_id="$1"
  jq -r ".roles[] | select(.id == \"$role_id\") | .display // \"$role_id\"" "$SETUP_OPTIONS"
}

session_layout()  { jq -r '.layout // "tabs"' "$SETUP_CONFIG"; }
master_tool_id()  { jq -r '.master.tool'  "$SETUP_CONFIG"; }
master_model()    { jq -r '.master.model' "$SETUP_CONFIG"; }
pipeline_role_id() { jq -r ".pipeline[\"$1\"] // \"none\"" "$SETUP_CONFIG"; }
role_tool_id()    { jq -r ".roles[\"$1\"].tool"  "$SETUP_CONFIG"; }
role_model()      { jq -r ".roles[\"$1\"].model" "$SETUP_CONFIG"; }
auto_approve_enabled() { jq -r ".auto_approve // false" "$SETUP_CONFIG"; }

# Attach or switch to session depending on whether we're already inside tmux
tmux_attach() {
  local session="$1"
  if [[ -n "${TMUX:-}" ]]; then
    # Already inside a tmux session — switch to target session
    tmux switch-client -t "$session"
  else
    tmux attach-session -t "$session"
  fi
}

build_agent_cmd() {
  local tool_id="$1"
  local model_id="$2"
  local role_id="$3"
  local phase_display="${4:-}"
  local fire_signal="${5:-}"

  local cmd_name
  cmd_name="$(get_tool_command "$tool_id")"

  if [[ -z "$cmd_name" ]]; then
    printf "echo 'ERROR: No command found for tool \"%s\". Check setup-options.json.'; sleep 9999" "$tool_id"
    return
  fi

  if ! command -v "$cmd_name" > /dev/null 2>&1; then
    printf "echo 'ERROR: %s not found in PATH. Run make setup to install it.'; sleep 9999" "$cmd_name"
    return
  fi

  local role_prompt_file="$ROLES_DIR/${role_id}.md"
  local model_flag=""

  local auto_approve
  auto_approve="$(auto_approve_enabled)"

  case "$tool_id" in
    gemini)      
      model_flag="-m $model_id"
      # Pipeline agents run with stdin redirected from /dev/null (< /dev/null).
      # Without -y, every tool-call confirmation prompt receives EOF and the tool
      # is silently CANCELLED — so write_file, run_shell_command, etc. never execute.
      # The agent then falls back to printing content to the terminal instead of disk.
      # ALWAYS force -y for pipeline agents; respect auto_approve only for master.
      if [[ "$auto_approve" == "true" ]] || [[ -n "$fire_signal" ]]; then
        model_flag+=" -y"
      fi
      ;;
    claude)      
      model_flag="--model $model_id"
      if [[ "$auto_approve" == "true" ]] || [[ -n "$fire_signal" ]]; then
        model_flag+=" --dangerously-skip-permissions"
      fi
      ;;
    cursor)      
      model_flag="--model $model_id" 
      ;;
  esac

  # Banner — for pipeline agents (fire_signal set) we skip printing the role prompt
  # text to the terminal; the AI already receives it via the injected prompt argument.
  local banner_cmd
  banner_cmd="printf '\\033[36m━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\\033[0m\\n'"
  banner_cmd+=" && printf ' \\033[1mRole: ${role_id}  |  Tool: ${tool_id}  |  Model: ${model_id}\\033[0m\\n'"
  banner_cmd+=" && printf '\\033[36m━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\\033[0m\\n'"

  local prompt_inject=""
  if [[ -f "$role_prompt_file" ]]; then
    if [[ -z "$fire_signal" ]]; then
      # Interactive: Keep terminal clean by telling the AI to read its file, rather
      # than injecting the full markdown text (which chat CLIs will echo).
      prompt_inject="Please silently read your system instructions from ${role_prompt_file} and execute YOUR IMMEDIATE STARTUP ACTION."
    else
      # Pipeline: Inject the full text inline to guarantee strict adherence
      prompt_inject="\$(cat '${role_prompt_file}')"
    fi
  fi

  # For pipeline agents: enrich the prompt with the spec context and exact output path
  # so the AI knows exactly where to write its artifact (fixes plan.md 404 errors).
  if [[ -n "$fire_signal" ]]; then
    # Determine the expected output file at launch-time.
    # IMPORTANT: use double quotes so $REPO_ROOT expands NOW (at script-generation time).
    # Single quotes inside $() prevent variable expansion entirely — the find command
    # would literally search for a path named '$REPO_ROOT/specs/' and always return empty.
    local baked_feature_dir
    baked_feature_dir="$(find "$REPO_ROOT/specs/" -mindepth 1 -maxdepth 1 -type d 2>/dev/null | sort -r | head -n 1)"

    # Determine which upstream artifact to inject as context
    local context_file_name=""
    case "$fire_signal" in
      plan)   context_file_name="spec" ;;
      tasks)  context_file_name="plan" ;;
      impl)   context_file_name="tasks" ;;
      review) context_file_name="tasks" ;;
    esac

    prompt_inject+="

IMPORTANT PIPELINE DIRECTIVE:
You are assigned to Phase: ${phase_display}.
Your WORKING DIRECTORY is: ${REPO_ROOT}

OUTPUT FILE — you MUST write your artifact to this exact path:
  ${baked_feature_dir}/${fire_signal}.md

HOW TO WRITE THE FILE (choose one — DO NOT use subagent tools):
  Option 1 — use the built-in write_file tool with the exact path above.
  Option 2 — use run_shell_command with a heredoc:
    mkdir -p '${baked_feature_dir}' && cat > '${baked_feature_dir}/${fire_signal}.md' << 'PLAN_EOF'
    [your full content here]
    PLAN_EOF

DO NOT print the content to the terminal only.
DO NOT delegate to a generalist or any other subagent — write the file directly."

    if [[ -n "$context_file_name" ]]; then
      prompt_inject+="

CONTEXT — read the following upstream artifact before starting:
  ${baked_feature_dir}/${context_file_name}.md"
    fi

    prompt_inject+="

Complete your phase tasks as thoroughly as possible. When you finish, exit immediately."
  fi

  local cli_invoke="$cmd_name"
  [[ -n "$model_flag" ]] && cli_invoke+=" $model_flag"
  if [[ -n "$prompt_inject" ]]; then
    case "$tool_id" in
      claude) cli_invoke+=" -p \"${prompt_inject}\"" ;;
      cursor) 
        # Cursor is a GUI IDE, not a headless CLI. Pass the workspace dir instead of invalid flags.
        cli_invoke="printf \"\n\033[33m⚠️  You have selected Cursor IDE as the agent tool.\033[0m\n\""
        cli_invoke+=" && printf \"Please use the Cursor Composer/Chat to complete this phase manually.\n\""
        cli_invoke+=" && echo \"${prompt_inject}\" | pbcopy 2>/dev/null || true"
        cli_invoke+=" && printf \"\033[32mThe prompt has been copied to your clipboard! Paste it into Cursor.\033[0m\n\""
        cli_invoke+=" && cursor \"${REPO_ROOT}\" 2>/dev/null || true"
        ;;
      *)      cli_invoke+=" \"${prompt_inject}\"" ;;
    esac
  fi

  local script_block
  if [[ -n "$fire_signal" ]]; then
    # Force standard CLIs into non-interactive batch mode so they exit after outputting
    if [[ "$tool_id" != "cursor" ]]; then
      cli_invoke+=" < /dev/null"
    fi
    
    # Escape inner variables carefully so they only evaluate at runtime inside the launcher scripts
    script_block=$(cat <<EOF
  while true; do
    # ── Clear Gemini session env vars (Root Cause 3: subagent tools disabled) ──
    # When 'make agents' is run from inside an active Gemini CLI session, tmux panes
    # inherit its env vars. Gemini's LocalAgentExecutor then detects nesting and
    # disables subagent tools (codebase_investigator, cli_help, etc.).
    # Unsetting these markers ensures every pipeline pane starts as a top-level agent.
    unset GEMINI_API_KEY_SOURCE GEMINI_CONVERSATION_ID GEMINI_PARENT_SESSION \
          GEMINI_AGENT_ID GEMINI_CODE_EXECUTION_SANDBOX \
          GOOGLE_CLOUD_PROJECT GOOGLE_GENAI_USE_VERTEXAI \
          GOOGLE_VERTEX_REGION CLOUD_ML_JOB_ID 2>/dev/null || true
    # ──────────────────────────────────────────────────────────────────────────
    $banner_cmd
    $cli_invoke
    
    printf "\n\n\033[36m━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\033[0m\n"
    printf "\033[1m✅  %s Phase Output Generated\033[0m\n" "$phase_display"
    printf "\033[36m━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\033[0m\n\n"

    # Determine target file dynamically (searching inside specs/)
    feature_dir=\$(find "$REPO_ROOT/specs/" -mindepth 1 -maxdepth 1 -type d | sort -r | head -n 1)
    target_file="\$feature_dir/${fire_signal}.md"

    # ── GATE: required output file must exist before user can Approve ──────────
    while [[ ! -f "\$target_file" ]]; do
      printf "\033[31m✗ ERROR: Required output file not found:\033[0m\n"
      printf "    %s\n" "\$target_file"
      printf "\033[33m  The agent did not create the expected file (or you are using a manual IDE). Possible causes:\033[0m\n"
      printf "    • The AI exited before writing the file (check its output above)\n"
      printf "    • The agent wrote to a different path — check output above\n"
      printf "    • You are using Cursor/IDE and need to generate it manually\n\n"
      
      error_action=\$(gum choose "Retry Agent" "Check Again (I saved it)" "Create empty placeholder manually")
      if [[ "\$error_action" == "Retry Agent" ]]; then
        printf "\033[33m  Re-running the agent now...\033[0m\n\n"
        break
      elif [[ "\$error_action" == "Check Again (I saved it)" ]]; then
        if [[ -f "\$target_file" ]]; then
          break
        else
          printf "\033[31mStill not found. Please save the file.\033[0m\n\n"
          sleep 2
        fi
      else
        # Manual fallback
        mkdir -p "\$(dirname "\$target_file")"
        touch "\$target_file"
        printf "\033[32mEmpty placeholder created.\033[0m\n"
        break
      fi
    done

    # If the user chose "Retry Agent", the file still doesn't exist.
    # We should continue the outer loop to run the agent again!
    if [[ ! -f "\$target_file" ]]; then
      continue
    fi
    # ──────────────────────────────────────────────────────────────────────────

    if [[ "$fire_signal" == "spec" || "$fire_signal" == "plan" || "$fire_signal" == "tasks" || "$fire_signal" == "impl" || "$fire_signal" == "review" ]]; then
      read_mode=\$(gum choose "Read in Terminal" "Read in Editor" "Skip Reading")
      if [[ "\$read_mode" == "Read in Terminal" ]]; then
        cat "\$target_file" | gum format
      elif [[ "\$read_mode" == "Read in Editor" ]]; then
        if command -v cursor >/dev/null 2>&1; then
          cursor "\$target_file" 2>/dev/null || open "\$target_file" 2>/dev/null
        else
          open "\$target_file" 2>/dev/null
        fi
      fi
    fi

    decision=\$(gum choose "Approve" "Needs changes")
    if [[ "\$decision" == "Approve" ]]; then
      touch "$SIGNALS_DIR/${fire_signal}.ready"
      printf "\033[32m✅ Proceeding to next phase.\033[0m\n"
      break
    else
      printf "\n\033[33m⏳ Waiting for Master agent to tell me the next step...\033[0m\n"
      tmux select-pane -t "${SESSION}:0.0" 2>/dev/null || true
      tmux send-keys -t "${SESSION}:0.0" "I have rejected the ${phase_display} phase. Please ask me what needs changing, update the file, and then execute precisely this command in your shell tool to resume the pipeline: touch $SIGNALS_DIR/${fire_signal}.revision.ready" Enter
      
      while [[ ! -f "$SIGNALS_DIR/${fire_signal}.revision.ready" ]]; do
        sleep 3
      done
      rm -f "$SIGNALS_DIR/${fire_signal}.revision.ready"
    fi
  done
EOF
)
  else
    # Master agent — interactive mode.
    # Do NOT auto-restart: wait for the user to press Enter so rapid restart
    # loops do not eat context for small tasks.
    script_block=$(cat <<EOF
  while true; do
    $banner_cmd
    $cli_invoke
    
    printf "\n\n\033[33m━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\033[0m\n"
    printf "\033[33mMaster Agent session ended.\033[0m\n"
    printf "Press \033[1mEnter\033[0m to restart, or \033[1mCtrl+C\033[0m to exit.\n"
    printf "\033[33m━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\033[0m\n"
    read -r _master_restart_input
  done
EOF
)
  fi

  printf "%s" "$script_block"
}

# Block current pane until a signal file exists (polling every 3s)
# Usage: embed this as the shell command prefix for a blocking pane
wait_cmd() {
  local signal="$1"
  local display="$2"
  # Use exact bash inline string so tmux send-keys doesn't misinterpret actual newlines
  printf "signal_file='%s/%s.ready'; printf '\\\\033[33m⏳ Waiting for %s to complete (%s.ready)...\\\\033[0m\\\\n'; while [[ ! -f \"\$signal_file\" ]]; do sleep 3; done; tmux resize-pane -t \"\$TMUX_PANE\" -y 999 2>/dev/null || true; tmux select-pane -t \"\$TMUX_PANE\" 2>/dev/null || true; printf '\\\\033[32m✅ Signal received: %s\\\\033[0m\\\\n\\\\n'" "$SIGNALS_DIR" "$signal" "$display" "$signal" "$signal"
}

# ── Session management ────────────────────────────────────────────────────────

cmd_kill() {
  if tmux has-session -t "$SESSION" 2>/dev/null; then
    tmux kill-session -t "$SESSION"
    printf "${GREEN}✅ Session '%s' killed.${RESET}\n" "$SESSION"
  else
    printf "${YELLOW}No active session named '%s'.${RESET}\n" "$SESSION"
  fi
}

cmd_status() {
  if ! tmux has-session -t "$SESSION" 2>/dev/null; then
    printf "${YELLOW}Session '%s' is not running.${RESET}\n" "$SESSION"
    return
  fi
  printf "${BOLD}${CYAN}Active tmux windows in session '%s':${RESET}\n" "$SESSION"
  tmux list-windows -t "$SESSION"
}

cmd_start() {
  require_tmux
  require_config

  # Check if we should branch first before running anything
  local current_branch
  current_branch="$(git branch --show-current 2>/dev/null || true)"
  local default_branch
  default_branch="$(jq -r '.default_branch // "develop"' "$SETUP_CONFIG")"
  
  if [[ -n "$current_branch" ]] && [[ "$current_branch" == "$default_branch" || "$current_branch" == "master" || "$current_branch" == "main" ]]; then
    printf "\n${YELLOW}WARNING: You are on the default branch (${BOLD}%s${RESET}${YELLOW}).${RESET}\n" "$current_branch"
    if gum confirm "Create a new feature branch before spinning up the multi-agent session?"; then
      if ! make specify; then
        printf "\n${RED}Feature branch creation cancelled or failed. Exiting.${RESET}\n"
        exit 1
      fi
    fi
  fi

  # Ensure signals dir exists
  mkdir -p "$SIGNALS_DIR"

  # Clear stale signals from previous run
  rm -f "$SIGNALS_DIR"/*.ready 2>/dev/null || true

  # Setup launcher scripts directory (avoids tmux send-keys character limit and quote issues)
  local launch_dir="$REPO_ROOT/.specify/launch"
  mkdir -p "$launch_dir"
  rm -f "$launch_dir"/*.sh 2>/dev/null || true

  if tmux has-session -t "$SESSION" 2>/dev/null; then
    printf "${YELLOW}Session '%s' is already running.${RESET}\n" "$SESSION"
    printf "Switching to it... (use ${BOLD}Ctrl+B D${RESET} to detach, ${BOLD}Ctrl+B w${RESET} to list windows)\n\n"
    tmux_attach "$SESSION"
    return
  fi

  printf "\n${BOLD}${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${RESET}\n"
  printf " ${BOLD}Spinning up Multi-Agent tmux session: %s${RESET}\n" "$SESSION"
  printf "${BOLD}${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${RESET}\n\n"

  local layout
  layout="$(session_layout)"

  # ── Window 0: Master Agent ────────────────────────────────────────────────
  local m_tool m_model m_cmd
  m_tool="$(master_tool_id)"
  m_model="$(master_model)"
  m_cmd="$(build_agent_cmd "$m_tool" "$m_model" "master")"
  echo "$m_cmd" > "$launch_dir/master.sh"

  printf "  ${GREEN}[0]${RESET} ${BOLD}Master Agent${RESET}  (%s / %s)\n" "$m_tool" "$m_model"

  if [[ "$layout" == "dashboard" ]]; then
    tmux new-session -d -s "$SESSION" -n "dashboard" -x 260 -y 60
    # Enable sticky labels at the top of each pane
    tmux set-window-option -t "${SESSION}:dashboard" pane-border-status top
    tmux set-window-option -t "${SESSION}:dashboard" pane-border-format " #[bold]#{pane_title}#[default] "
    
    tmux send-keys -t "${SESSION}:dashboard" "cd '$REPO_ROOT' && source '$launch_dir/master.sh'" Enter
    tmux select-pane -t "${SESSION}:dashboard" -T "👑 Master Agent  |  Tool: ${m_tool}  |  Model: ${m_model}"
  else
    tmux new-session -d -s "$SESSION" -n "master" -x 220 -y 50
    tmux set-window-option -t "${SESSION}:master" pane-border-status top
    tmux set-window-option -t "${SESSION}:master" pane-border-format " #[bold]#{pane_title}#[default] "
    
    tmux send-keys -t "${SESSION}:master" "cd '$REPO_ROOT' && source '$launch_dir/master.sh'" Enter
    tmux select-pane -t "${SESSION}:master" -T "👑 Master Agent  |  Tool: ${m_tool}  |  Model: ${m_model}"
  fi

  # ── Windows 1..N: Role agents ─────────────────────────────────────────────
  local win_index=1
  local -a phases=("plan" "tasks" "impl" "review")

  for phase in "${phases[@]}"; do
    local role_id
    role_id="$(pipeline_role_id "$phase")"
    
    [[ "$role_id" == "none" || -z "$role_id" || "$role_id" == "null" ]] && continue

    local r_tool r_model r_display r_cmd
    r_tool="$(role_tool_id "$role_id")"
    r_model="$(role_model "$role_id")"
    r_display="$(get_role_display "$role_id")"

    local block_signal=""
    local phase_display=""
    case "$phase" in
      plan)   block_signal="spec"; phase_display="Plan" ;;
      tasks)  block_signal="plan"; phase_display="Tasks" ;;
      impl)   block_signal="tasks"; phase_display="Implement" ;;
      review) block_signal="impl"; phase_display="Analyze/QA" ;;
    esac

    r_cmd="$(build_agent_cmd "$r_tool" "$r_model" "$role_id" "$phase_display" "$phase")"
    
    printf "  ${GREEN}[%d]${RESET} ${BOLD}%-15s${RESET} -> %s  (%s / %s)" "$win_index" "Phase ${win_index}: ${phase_display}" "$r_display" "$r_tool" "$r_model"

    if [[ -n "$block_signal" ]]; then
      printf "  ${YELLOW}⏳ blocks on: %s.ready${RESET}" "$block_signal"
      local wait_prefix
      wait_prefix="$(wait_cmd "$block_signal" "Master Agent ($block_signal step)")"
      r_cmd="${wait_prefix} && ${r_cmd}"
    fi
    printf "\n"

    echo "$r_cmd" > "$launch_dir/phase_${win_index}.sh"

    if [[ "$layout" == "dashboard" ]]; then
      tmux split-window -t "${SESSION}:0" -v
      tmux send-keys -t "${SESSION}:0.${win_index}" "cd '$REPO_ROOT' && source '$launch_dir/phase_${win_index}.sh'" Enter
      tmux select-pane -t "${SESSION}:0.${win_index}" -T "🤖 Phase ${win_index}: ${phase_display} (${r_display})  |  Tool: ${r_tool}  |  Model: ${r_model}"
    else
      tmux new-window -t "${SESSION}:${win_index}" -n "$phase"
      tmux set-window-option -t "${SESSION}:${win_index}" pane-border-status top
      tmux set-window-option -t "${SESSION}:${win_index}" pane-border-format " #[bold]#{pane_title}#[default] "
      tmux send-keys -t "${SESSION}:${win_index}" "cd '$REPO_ROOT' && source '$launch_dir/phase_${win_index}.sh'" Enter
      tmux select-pane -t "${SESSION}:${win_index}" -T "🤖 Phase ${win_index}: ${phase_display} (${r_display})  |  Tool: ${r_tool}  |  Model: ${r_model}"
    fi

    win_index=$((win_index + 1))
  done

  if [[ "$layout" == "dashboard" ]]; then
    # Dashboard mode: Master agent gets ~100 characters on the left, others tile on the right
    tmux select-layout -t "${SESSION}:0" main-vertical
    tmux set-window-option -t "${SESSION}:0" main-pane-width 100
    tmux select-layout -t "${SESSION}:0" main-vertical
    
    # ENSURE FOCUS ON MASTER PANE
    tmux select-pane -t "${SESSION}:0.0"
    
    printf "\n${BOLD}${GREEN}✅ Session started in Dashboard mode (Master left, team right).${RESET}\n"
    printf "\n"
    printf "  ${BOLD}tmux cheatsheet (Panes):${RESET}\n"
    printf "  • ${CYAN}Ctrl+B arrow-keys${RESET} — Move between panes\n"
    printf "  • ${CYAN}Ctrl+B z${RESET}          — Zoom current pane full screen (press again to unzoom)\n"
    printf "  • ${CYAN}Ctrl+B D${RESET}          — Detach (session keeps running)\n"
  else
    tmux select-window -t "${SESSION}:0"
    printf "\n${BOLD}${GREEN}✅ Session started with %d window(s).${RESET}\n" "$win_index"
    printf "\n"
    printf "  ${BOLD}tmux cheatsheet (Windows):${RESET}\n"
    printf "  • ${CYAN}Ctrl+B [0-9]${RESET}      — Switch to window by number\n"
    printf "  • ${CYAN}Ctrl+B n / p${RESET}      — Next / previous window\n"
    printf "  • ${CYAN}Ctrl+B D${RESET}          — Detach (session keeps running)\n"
  fi
  printf "  • ${CYAN}make agents${RESET}       — Re-attach any time\n"
  printf "  • ${CYAN}make agents-kill${RESET}  — Destroy all agents\n"
  printf "\n"
  printf "  ${YELLOW}Signal files location: .specify/signals/${RESET}\n"
  printf "  ${YELLOW}Write a .ready file to unblock downstream agents:${RESET}\n"
  printf "    touch .specify/signals/spec.ready\n"
  printf "    touch .specify/signals/plan.ready\n"
  printf "    touch .specify/signals/tasks.ready\n"
  printf "\n"

  tmux_attach "$SESSION"
}

# ── Entry point ───────────────────────────────────────────────────────────────

CMD="${1:-start}"
case "$CMD" in
  start)  cmd_start ;;
  kill)   cmd_kill ;;
  status) cmd_status ;;
  *)
    printf "Usage: %s [start|kill|status]\n" "$(basename "$0")" >&2
    exit 1
    ;;
esac
