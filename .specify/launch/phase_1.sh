signal_file='/Users/michael/Sites/WeProDev/wpd-message-gateway/.specify/signals/spec.ready'; printf '\033[33m⏳ Waiting for Master Agent (spec step) to complete (spec.ready)...\033[0m\n'; while [[ ! -f "$signal_file" ]]; do sleep 3; done; tmux resize-pane -t "$TMUX_PANE" -y 999 2>/dev/null || true; tmux select-pane -t "$TMUX_PANE" 2>/dev/null || true; printf '\033[32m✅ Signal received: spec\033[0m\n\n' &&   while true; do
    # ── Clear Gemini session env vars (Root Cause 3: subagent tools disabled) ──
    # When 'make agents' is run from inside an active Gemini CLI session, tmux panes
    # inherit its env vars. Gemini's LocalAgentExecutor then detects nesting and
    # disables subagent tools (codebase_investigator, cli_help, etc.).
    # Unsetting these markers ensures every pipeline pane starts as a top-level agent.
    unset GEMINI_API_KEY_SOURCE GEMINI_CONVERSATION_ID GEMINI_PARENT_SESSION           GEMINI_AGENT_ID GEMINI_CODE_EXECUTION_SANDBOX           GOOGLE_CLOUD_PROJECT GOOGLE_GENAI_USE_VERTEXAI           GOOGLE_VERTEX_REGION CLOUD_ML_JOB_ID 2>/dev/null || true
    # ──────────────────────────────────────────────────────────────────────────
    printf '\033[36m━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\033[0m\n' && printf ' \033[1mRole: team-lead  |  Tool: gemini  |  Model: gemini-3-pro-preview\033[0m\n' && printf '\033[36m━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\033[0m\n'
    gemini -m gemini-3-pro-preview -y "$(cat '/Users/michael/Sites/WeProDev/wpd-message-gateway/scripts/roles/team-lead.md')

IMPORTANT PIPELINE DIRECTIVE:
You are assigned to Phase: Plan.
Your WORKING DIRECTORY is: /Users/michael/Sites/WeProDev/wpd-message-gateway

OUTPUT FILE — you MUST write your artifact to this exact path:
  /Users/michael/Sites/WeProDev/wpd-message-gateway/specs/017-refactor-readmefile-technical-writer/plan.md

HOW TO WRITE THE FILE (choose one — DO NOT use subagent tools):
  Option 1 — use the built-in write_file tool with the exact path above.
  Option 2 — use run_shell_command with a heredoc:
    mkdir -p '/Users/michael/Sites/WeProDev/wpd-message-gateway/specs/017-refactor-readmefile-technical-writer' && cat > '/Users/michael/Sites/WeProDev/wpd-message-gateway/specs/017-refactor-readmefile-technical-writer/plan.md' << 'PLAN_EOF'
    [your full content here]
    PLAN_EOF

DO NOT print the content to the terminal only.
DO NOT delegate to a generalist or any other subagent — write the file directly.

CONTEXT — read the following upstream artifact before starting:
  /Users/michael/Sites/WeProDev/wpd-message-gateway/specs/017-refactor-readmefile-technical-writer/spec.md

Complete your phase tasks as thoroughly as possible. When you finish, exit immediately." < /dev/null
    
    printf "\n\n\033[36m━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\033[0m\n"
    printf "\033[1m✅  %s Phase Output Generated\033[0m\n" "Plan"
    printf "\033[36m━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\033[0m\n\n"

    # Determine target file dynamically (searching inside specs/)
    feature_dir=$(find "/Users/michael/Sites/WeProDev/wpd-message-gateway/specs/" -mindepth 1 -maxdepth 1 -type d | sort -r | head -n 1)
    target_file="$feature_dir/plan.md"

    # ── GATE: required output file must exist before user can Approve ──────────
    while [[ ! -f "$target_file" ]]; do
      printf "\033[31m✗ ERROR: Required output file not found:\033[0m\n"
      printf "    %s\n" "$target_file"
      printf "\033[33m  The agent did not create the expected file (or you are using a manual IDE). Possible causes:\033[0m\n"
      printf "    • The AI exited before writing the file (check its output above)\n"
      printf "    • The agent wrote to a different path — check output above\n"
      printf "    • You are using Cursor/IDE and need to generate it manually\n\n"
      
      error_action=$(gum choose "Retry Agent" "Check Again (I saved it)" "Create empty placeholder manually")
      if [[ "$error_action" == "Retry Agent" ]]; then
        printf "\033[33m  Re-running the agent now...\033[0m\n\n"
        break
      elif [[ "$error_action" == "Check Again (I saved it)" ]]; then
        if [[ -f "$target_file" ]]; then
          break
        else
          printf "\033[31mStill not found. Please save the file.\033[0m\n\n"
          sleep 2
        fi
      else
        # Manual fallback
        mkdir -p "$(dirname "$target_file")"
        touch "$target_file"
        printf "\033[32mEmpty placeholder created.\033[0m\n"
        break
      fi
    done

    # If the user chose "Retry Agent", the file still doesn't exist.
    # We should continue the outer loop to run the agent again!
    if [[ ! -f "$target_file" ]]; then
      continue
    fi
    # ──────────────────────────────────────────────────────────────────────────

    if [[ "plan" == "spec" || "plan" == "plan" || "plan" == "tasks" || "plan" == "impl" || "plan" == "review" ]]; then
      read_mode=$(gum choose "Read in Terminal" "Read in Editor" "Skip Reading")
      if [[ "$read_mode" == "Read in Terminal" ]]; then
        cat "$target_file" | gum format
      elif [[ "$read_mode" == "Read in Editor" ]]; then
        if command -v cursor >/dev/null 2>&1; then
          cursor "$target_file" 2>/dev/null || open "$target_file" 2>/dev/null
        else
          open "$target_file" 2>/dev/null
        fi
      fi
    fi

    decision=$(gum choose "Approve" "Needs changes")
    if [[ "$decision" == "Approve" ]]; then
      touch "/Users/michael/Sites/WeProDev/wpd-message-gateway/.specify/signals/plan.ready"
      printf "\033[32m✅ Proceeding to next phase.\033[0m\n"
      break
    else
      printf "\n\033[33m⏳ Waiting for Master agent to tell me the next step...\033[0m\n"
      tmux select-pane -t "wpd-gateway:0.0" 2>/dev/null || true
      tmux send-keys -t "wpd-gateway:0.0" "I have rejected the Plan phase. Please ask me what needs changing, update the file, and then execute precisely this command in your shell tool to resume the pipeline: touch /Users/michael/Sites/WeProDev/wpd-message-gateway/.specify/signals/plan.revision.ready" Enter
      
      while [[ ! -f "/Users/michael/Sites/WeProDev/wpd-message-gateway/.specify/signals/plan.revision.ready" ]]; do
        sleep 3
      done
      rm -f "/Users/michael/Sites/WeProDev/wpd-message-gateway/.specify/signals/plan.revision.ready"
    fi
  done
