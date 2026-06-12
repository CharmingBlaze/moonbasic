#!/usr/bin/env sh
set -e
cd "$(dirname "$0")"
if [ -x ./moonbasic ]; then
  ./moonbasic install-vscode
elif command -v moonbasic >/dev/null 2>&1; then
  moonbasic install-vscode
else
  echo "Could not find moonbasic in this folder or on PATH."
  echo "Extract the moonBASIC release archive first, then run: ./INSTALL-VSCODE.sh"
  exit 1
fi
echo "Open VS Code or Cursor and open any .mb file."
