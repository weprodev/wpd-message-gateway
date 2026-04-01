#!/usr/bin/env bash

set -euo pipefail

printf "Installing Cursor CLI using official installer...\n"
curl https://cursor.com/install -fsS | bash
