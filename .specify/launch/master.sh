  while true; do
    printf '\033[36m━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\033[0m\n' && printf ' \033[1mRole: master  |  Tool: gemini  |  Model: gemini-3-pro-preview\033[0m\n' && printf '\033[36m━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\033[0m\n'
    gemini -m gemini-3-pro-preview "Please silently read your system instructions from /Users/michael/Sites/WeProDev/wpd-message-gateway/scripts/roles/master.md and execute YOUR IMMEDIATE STARTUP ACTION."
    
    printf "\n\n\033[33m━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\033[0m\n"
    printf "\033[33mMaster Agent session ended.\033[0m\n"
    printf "Press \033[1mEnter\033[0m to restart, or \033[1mCtrl+C\033[0m to exit.\n"
    printf "\033[33m━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\033[0m\n"
    read -r _master_restart_input
  done
