#!/usr/bin/env sh
# Print shell lines to put moonbasic/moonrun on PATH (optional — IDE does not need this).
set -e
HERE="$(cd "$(dirname "$0")" && pwd)"
echo ""
echo "moonBASIC — optional PATH setup"
echo "  Folder: $HERE"
echo ""
echo "The IDE finds moonbasic/moonrun beside itself. PATH is only for terminal use."
echo ""
echo "Add to ~/.bashrc or ~/.zshrc:"
echo ""
echo "  export PATH=\"$HERE:\$PATH\""
echo ""
case "$(uname -s)" in
  Darwin)
    PROFILE="$HOME/.zprofile"
    [ -f "$HOME/.zshrc" ] && PROFILE="$HOME/.zshrc"
    ;;
  *)
    PROFILE="$HOME/.bashrc"
    [ -f "$HOME/.zshrc" ] && PROFILE="$HOME/.zshrc"
    ;;
esac
printf "Append that line to %s now? [y/N] " "$PROFILE"
read -r ans || true
case "$ans" in
  y|Y|yes|YES)
    echo "" >> "$PROFILE"
    echo "# moonBASIC toolchain ($(date +%Y-%m-%d))" >> "$PROFILE"
    echo "export PATH=\"$HERE:\$PATH\"" >> "$PROFILE"
    echo "Appended. Run:  source $PROFILE"
    ;;
  *)
    echo "Skipped. Copy the export line above when you want."
    ;;
esac
